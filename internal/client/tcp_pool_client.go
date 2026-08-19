package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"example/bootstrap"
	types "example/types"
)

const (
	tcpMaxBodyLength = 1 << 12 // 4KB，避免錯誤/惡意的長度前綴把記憶體打爆
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

var (
	ErrTcpBodyTooLarge   = errors.New("tcp: body too large")
	ErrTcpMethodNotFound = errors.New("tcp: method not found")
)

// tcpConn 是池子裡實際借還的單位：一條 net.Conn 綁一個專屬的 bufio.Reader。
// reader 要跟著連線的生命週期走、不能每次呼叫都重建（理由見 DecodeFrame 註解），
// 所以借出去、還回來都是這一整包一起借還，不會有連線和 reader 對不上的情況。
type TcpConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

// NewTcpConn 自己讀 bootstrap.CONFIG.CLIENTS.TCP 撥號，不接受外部注入的 net.Conn，
// 跟 TcpClient.NewTcpClient 的做法一致。
func NewTcpConn() (*TcpConn, error) {
	oConn, err := net.Dial("tcp", bootstrap.CONFIG.CLIENTS.TCP.HOSTS[0]+":"+bootstrap.CONFIG.CLIENTS.TCP.PORTS[0])
	if err != nil {
		return nil, err
	}

	return &TcpConn{
		conn:   oConn,
		reader: bufio.NewReader(oConn),
	}, nil
}

func (oSelf *TcpConn) Write(aBuf []byte) (int, error) {
	return oSelf.conn.Write(aBuf)
}

func (oSelf *TcpConn) Close() error {
	return oSelf.conn.Close()
}

// TcpPoolClient 自帶連線池：外部呼叫端不需要知道底層借了哪條連線，呼叫 method 時
// 自動借一條、用完自動還回去。channel 本身就是併發安全的借還機制，多個 goroutine
// 可以同時各自佔用一條連線並發送 request，不會像單一連線+mutex 那樣互相卡住。
type TcpPoolClient struct {
	pool chan *TcpConn
}

// NewTcpPoolClient 建一個空池子，池子大小（最多留幾條閒置連線）讀
// bootstrap.CONFIG.CLIENTS.TCP.POOL；連線是真正要用時才現撥（見 NewTcpConn），
// 不會一開始就撥滿。
func NewTcpPoolClient() *TcpPoolClient {
	return &TcpPoolClient{
		pool: make(chan *TcpConn, bootstrap.CONFIG.CLIENTS.TCP.POOL),
	}
}

// get 借一條連線：池子裡有現成的就直接拿，沒有的話（池子空或已達上限被借光）就現撥一條新的，
// 所以真正並發數不受池子大小限制，池子大小只限制「閒置不用時最多留幾條」。
func (oSelf *TcpPoolClient) get() (*TcpConn, error) {
	select {
	case oConn := <-oSelf.pool: // 如果channel 有連結變量，就從channel 拿
		return oConn, nil
	default:
		return NewTcpConn() // 如果channel 無連結變量，就重新建立
	}
}

// put 把用完的連線還回池子；池子已經滿了就直接把這條連線關掉，不強留超過上限的閒置連線。
func (oSelf *TcpPoolClient) put(oTcpConn *TcpConn) {
	select {
	case oSelf.pool <- oTcpConn: // 如果channel 還有空間，就丟回去
	default: //
		oTcpConn.Close() // 如果channel 無了空間，就棄用
	}
}

// do 借一條連線交給 fn 做 write/read，業務代碼只要專心處理 I/O 本身，
// 不用管連線要還回池子還是關掉：fn 成功就還回去，失敗就直接關掉
// （連線的讀寫狀態已經跟 frame 對不上，不安全再借給下一個呼叫端）。
func (oSelf *TcpPoolClient) do(cFn func(oTcpConn *TcpConn) error) error {
	oTcpConn, err := oSelf.get()
	if err != nil {
		log.Println("tcp client: get conn failed:", err)
		return err
	}

	if err := cFn(oTcpConn); err != nil {
		oTcpConn.Close()
		return err
	}

	oSelf.put(oTcpConn)
	return nil
}

func (oSelf *TcpPoolClient) AdminAuthenticationAuthenticatorSignIn(sName string, sPassword string) (*types.TcpResponse, error) {

	var oResp types.TcpResponse
	err := oSelf.do(func(oTcpConn *TcpConn) error {
		aFrame, err := oSelf.EncodeFrame(types.TcpRequest{
			Method: "Admin.Authentication.Authenticator.SignIn",
			Param:  sName + ":" + sPassword,
		})
		if err != nil {
			log.Println("tcp client: encode failed:", err)
			return err
		}

		if _, err := oTcpConn.Write(aFrame); err != nil {
			log.Println("tcp client: write failed:", err)
			return err
		}

		if err := oSelf.DecodeFrame(oTcpConn.reader, &oResp); err != nil {
			log.Println("tcp client: decode failed:", err)
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &oResp, nil
}

func (oSelf *TcpPoolClient) EncodeFrame(oPayload any) ([]byte, error) {
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

func (oSelf *TcpPoolClient) DecodeFrame(oReader *bufio.Reader, oPayload any) error {
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
