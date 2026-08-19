package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// writeMessage/readMessage 跟 server 那份一模一樣——這是故意的，兩邊本來就要照
// 同一套規則組包/拆包才讀得懂彼此，不是誰依賴誰。

func main() {
	oConn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer oConn.Close()

	if err := writeMessage(oConn, "hello"); err != nil {
		panic(err)
	}

	oReader := bufio.NewReader(oConn)
	sResponse, err := readMessage(oReader)
	if err != nil {
		panic(err)
	}

	fmt.Println("client received:", sResponse)
}

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
