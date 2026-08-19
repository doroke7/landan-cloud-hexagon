package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

// TcpRequest/TcpResponse、encodeFrameAndWrite/decodeFrame 跟 server 那份一模一樣——
// 這是故意的，兩邊本來就要照同一套規則組包/拆包才讀得懂彼此。

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
	oConn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	oReq := TcpRequest{
		Code:   0,
		Method: "SignIn",
		Param:  "admin:520999",
	}

	if err := encodeFrame(oConn, oReq); err != nil {
		panic(err)
	}

	oReader := bufio.NewReader(oConn)

	var oResp TcpResponse
	if err := decodeFrame(oReader, &oResp); err != nil {
		panic(err)
	}

	fmt.Printf("client received: code=%d message=%s result=%s\n", oResp.Code, oResp.Message, oResp.Result)
}

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
