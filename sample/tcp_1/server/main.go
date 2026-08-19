package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// 這個範例故意寫得很陽春：沒有 JSON、沒有 method 分發、沒有連線池、沒有 RequestId，
// 只示範最基本的「TCP 框架」概念——用一個 4 byte 長度前綴解決黏包/拆包問題，
// 一次只處理一個 request、處理完才收下一個。整個檔案看完就懂框架的核心，
// 之後再去看正式版的 pkg/tcp_router.go、internal/client/tcp_mux_client.go
// 在這個基礎上多做了什麼、為什麼要多做（連線池、並發、RequestId 配對）。

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

	for {
		sMessage, err := readMessage(oReader)
		if err != nil {
			return // 連線關掉或壞掉，結束這個連線的處理
		}

		fmt.Println("server received:", sMessage)

		sResponse := "echo: " + sMessage
		if err := writeMessage(oConn, sResponse); err != nil {
			return
		}
	}
}

// writeMessage 把 message 前面加上 4 byte 長度前綴再寫出去。
func writeMessage(oWriter io.Writer, sMessage string) error {
	aBody := []byte(sMessage)

	aLengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(aLengthBuf, uint32(len(aBody)))

	if _, err := oWriter.Write(aLengthBuf); err != nil {
		return err
	}
	if _, err := oWriter.Write(aBody); err != nil {
		return err
	}
	return nil
}

// readMessage 先讀 4 byte 拿到長度，再讀滿那個長度——這就是解決黏包/拆包的關鍵，
// 不管 TCP 底層一次送幾個 byte 過來，讀滿了才算一個完整的 message。
func readMessage(oReader io.Reader) (string, error) {
	aLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(oReader, aLengthBuf); err != nil {
		return "", err
	}

	iLength := binary.BigEndian.Uint32(aLengthBuf)

	aBody := make([]byte, iLength)
	if _, err := io.ReadFull(oReader, aBody); err != nil {
		return "", err
	}

	return string(aBody), nil
}
