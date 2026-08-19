package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 這份在 sample/tcp_2 的框架基礎上（4-byte 長度前綴解決黏包/拆包、struct+JSON）
// 加上 RequestId，讓 client 可以在同一條連線上「同時」掛好幾筆還沒回應的 request——
// 這就是多路復用（multiplexing）。跟 tcp_2 的差異都標了 // IMPORTANT 多路。

type TcpRequest struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Method    string `json:"method"`
	Param     string `json:"param"`
}

type TcpResponse struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Result    string `json:"result"`
}

func main() {
	oListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	defer oListener.Close()

	fmt.Println("server listening on :9001")

	for {
		oConn, err := oListener.Accept()
		if err != nil {
			continue
		}

		go handleConn(oConn)
	}
}

func handleConn(oConn net.Conn) {
	defer oConn.Close()

	oReader := bufio.NewReader(oConn)

	// IMPORTANT 多路：多個 goroutine 會同時想 Write 到同一個 net.Conn，寫入本身
	// 不是原子的，這裡只保護「寫一個完整 frame」這個動作；不然兩筆 response 的
	// bytes 可能交錯寫進去，變成誰都解不出來的垃圾。
	var oWriteMu sync.Mutex

	for {
		var oReq TcpRequest
		if err := decodeFrame(oReader, &oReq); err != nil {
			return
		}

		/* IMPORTANT 多路：
		           一路：收到消息， 同步的寫入。一定是一個蘿蔔一個坑
				   多路：收到消息後馬上 異步開協程 寫入，可能次序不同，此時靠 request-id
		*/
		go func(oReq TcpRequest) {
			iDelayMs := 100
			if oReq.Method == "Slow" {
				iDelayMs = 400
			}
			time.Sleep(time.Duration(iDelayMs) * time.Millisecond)

			fmt.Printf("server processed id=%s method=%s param=%s\n", oReq.RequestId, oReq.Method, oReq.Param)

			oResp := TcpResponse{
				RequestId: oReq.RequestId,
				Code:      1,
				Message:   "成功處理 " + oReq.Method,
				Result:    "echo:" + oReq.Param,
			}

			oWriteMu.Lock()
			defer oWriteMu.Unlock()

			if err := encodeFrame(oConn, oResp); err != nil {
				fmt.Println("server write failed:", err)
			}
		}(oReq)
	}
}

// encodeFrame 把 oPayload 編成 JSON、加上 4-byte 長度前綴，寫出去。
func encodeFrame(oWriter io.Writer, oPayload any) error {
	aBody, err := json.Marshal(oPayload)
	if err != nil {
		return err
	}

	aFrame := make([]byte, 4+len(aBody))
	binary.BigEndian.PutUint32(aFrame[0:4], uint32(len(aBody)))
	copy(aFrame[4:], aBody)

	_, err = oWriter.Write(aFrame)
	return err
}

// decodeFrame 先讀 4 byte 拿到長度、讀滿那個長度（解決黏包/拆包），
// 再把 body 的 JSON 解到 oPayload（傳 &TcpRequest{} 或 &TcpResponse{}）。
func decodeFrame(oReader io.Reader, oPayload any) error {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return err
	}

	iLength := binary.BigEndian.Uint32(aLengthBuf)

	aBody := make([]byte, iLength)
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return err
	}

	return json.Unmarshal(aBody, oPayload)
}
