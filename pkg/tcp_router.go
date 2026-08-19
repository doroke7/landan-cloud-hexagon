package pkg

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	types "example/types"
)

/*
TCP 是位元組流（stream），沒有天然的訊息邊界，直接 conn.Read(buf) 會遇到：
  - 黏包：對方連續送兩筆小訊息，一次 Read 可能把兩筆都讀進來
  - 拆包：一筆訊息太大，一次 Read 讀不完，要跨好幾次 Read 才拼得齊

解法是自訂一個「帶長度前綴」的 frame 格式，讀取端固定「先讀長度、再讀滿長度」；
中間跨了幾次底層 Read 都由 bufio.Reader + io.ReadFull 自動處理。

Frame 格式（Big Endian）：

	┌──────────────┬─────────────────────────┐
	│  BodyLength  │   Body (JSON)           │
	│   4 bytes    │   N bytes               │
	└──────────────┴─────────────────────────┘

Body 是 JSON，內容依方向不同：
  - Request  ：{ code, method, param }
  - Response ：{ code, message, result }
*/

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
)

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

// tcpBodyPool 裝 DecodeFrame 每次讀 body 用的緩衝區，重複利用、減少每次 Read
// 都重新配置 []byte 造成的 GC 壓力。放 *[]byte 而不是 []byte，避免 slice header
// 被裝進 any 時多一次 heap allocation（sync.Pool 的慣用寫法）。
var tcpBodyPool = sync.Pool{
	New: func() any {
		aBuf := make([]byte, tcpMaxBodyLength)
		return &aBuf
	},
}

// tcpEncodeBufferPool 裝 EncodeFrame 組 JSON 用的暫存緩衝區。只在函式內部使用、
// 用完就 Put 回去，不會把池化的記憶體交出去給呼叫端——最終回傳的 aFrame 一定是
// 另外配置的一份，這樣才不會有「呼叫端還在用、池子卻把它重新分配出去」的風險。
var tcpEncodeBufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// TcpHandlerFunc 是一個 method 對應的處理方法，簽名統一，方便用 method name 當 key 做路由。
type TcpHandlerFunc func(oReq types.TcpRequest) types.TcpResponse

// TcpRouter 職責跟 ConsumerRouter 一樣：只負責把 method name 對應到一個處理方法，
// 不管 unmarshal／business 邏輯；Serve 時自己負責 accept 連線、讀 frame、分發、回包。
type TcpRouter struct {
	routes map[string]TcpHandlerFunc
}

func NewTcpRouter() *TcpRouter {
	return &TcpRouter{routes: make(map[string]TcpHandlerFunc)}
}

// HandleFunc 註冊一個 method 對應的處理方法，用法跟 ConsumerRouter.HandleFunc 一樣。
func (oSelf *TcpRouter) HandleFunc(sMethod string, fnHandler TcpHandlerFunc) *TcpRouter {
	oSelf.routes[sMethod] = fnHandler
	return oSelf
}

// Serve 監聽 sAddr，每個連線各自在自己的 goroutine 讀 frame、分發、回包，
// ctx 取消時關掉 listener，讓 Accept 中斷、Serve 返回。
func (oSelf *TcpRouter) Serve(ctx context.Context, sAddr string) error {
	oListener, err := net.Listen("tcp", sAddr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		oListener.Close()
	}()

	for {
		oConn, err := oListener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
		}

		go oSelf.serveConn(oConn)
	}
}

func (oSelf *TcpRouter) serveConn(oConn net.Conn) {
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	// oWriteMu 保護「寫一個完整 frame」這個動作。read 迴圈本身是循序的（一次只解一個
	// request），但每個 request 都丟到自己的 goroutine 去處理、回應順序可能跟收到的
	// 順序不一樣——這就是 RequestId 存在的意義，client 端靠它把亂序回來的 response
	// 配對回正確的呼叫方。多個 goroutine 同時 Write 到同一個 net.Conn 必須序列化，
	// 不然兩個 frame 的 bytes 會交錯寫進去，變成一包誰都解不出來的垃圾。
	var oWriteMu sync.Mutex

	for {
		var oReq types.TcpRequest
		if err := oSelf.DecodeFrame(oReader, &oReq); err != nil {
			return
		}

		go func(oReq types.TcpRequest) {
			oResp := oSelf.dispatch(oReq)
			oResp.RequestId = oReq.RequestId

			aFrame, err := oSelf.EncodeFrame(oResp)
			if err != nil {
				log.Printf("tcp: encode failed: method=%s err=%v", oReq.Method, err)
				return
			}

			oWriteMu.Lock()
			defer oWriteMu.Unlock()

			if _, err = oConn.Write(aFrame); err != nil {
				log.Printf("tcp: write failed: method=%s err=%v", oReq.Method, err)
			}
		}(oReq)
	}
}

func (oSelf *TcpRouter) dispatch(oReq types.TcpRequest) types.TcpResponse {
	fnHandler, ok := oSelf.routes[oReq.Method]
	if !ok {
		return types.TcpResponse{Code: -1, Message: ErrTcpMethodNotFound.Error()}
	}
	return fnHandler(oReq)
}

// EncodeFrame 把 oPayload（types.TcpRequest 或 types.TcpResponse）編碼成一個完整的 frame。
func (oSelf *TcpRouter) EncodeFrame(oPayload any) ([]byte, error) {
	oBuf := tcpEncodeBufferPool.Get().(*bytes.Buffer)
	oBuf.Reset()
	defer tcpEncodeBufferPool.Put(oBuf)

	if err := json.NewEncoder(oBuf).Encode(oPayload); err != nil {
		return nil, err
	}

	// json.Encoder.Encode 會多寫一個結尾的 \n，這裡要去掉，行為才會跟 json.Marshal 一致
	aBody := bytes.TrimRight(oBuf.Bytes(), "\n")

	if len(aBody) > tcpMaxBodyLength {
		return nil, ErrTcpBodyTooLarge
	}

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	return aFrame, nil
}

// DecodeFrame 從 reader 讀出一個完整 frame，解到 oPayload（傳 &types.TcpRequest{} 或 &types.TcpResponse{}）。
func (oSelf *TcpRouter) DecodeFrame(oReader *bufio.Reader, oPayload any) error {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return err
	}

	iBodyLength := binary.BigEndian.Uint32(aLengthBuf)
	if iBodyLength == 0 || iBodyLength > tcpMaxBodyLength {
		return ErrTcpBodyTooLarge
	}

	pBuf := tcpBodyPool.Get().(*[]byte)
	defer tcpBodyPool.Put(pBuf)

	aBody := (*pBuf)[:iBodyLength]
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return err
	}

	return json.Unmarshal(aBody, oPayload)
}
