package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// TcpRequest/TcpResponse 是要送上線的內容，跟正式版 types.TcpRequest/TcpResponse
// 是同一套欄位，用 JSON 編碼——比 "|" 字串分隔更明確，也不用擔心欄位內容剛好
// 出現分隔符號。核心的 4-byte 長度前綴（解決黏包/拆包）沒有變。

type TcpRequest struct {
	Code   int    `json:"code"`
	Method string `json:"method"`
	Param  string `json:"param"`
}

type TcpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func main() {
	oListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer oListener.Close()

	fmt.Println("server listening on :9000")

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

	/* IMPORTANT 多路：
	一路：收到消息， 同步的寫入。一定是一個蘿蔔一個坑
	多路：收到消息後馬上 異步開協程 寫入，可能次序不同，此時靠 request-id
	*/

	for {
		var oReq TcpRequest
		if err := decodeFrame(oReader, &oReq); err != nil {
			return
		}

		fmt.Printf("server received: code=%d method=%s param=%s\n", oReq.Code, oReq.Method, oReq.Param)

		oResp := TcpResponse{
			Code:    1,
			Message: "成功處理 " + oReq.Method,
			Result:  "echo:" + oReq.Param,
		}

		if err := encodeFrame(oConn, oResp); err != nil {
			return
		}
	}
}

// encodeFrameAndWrite 把 oPayload 編成 JSON、加上 4-byte 長度前綴，寫出去。
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
