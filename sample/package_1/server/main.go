package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// TCP 的「拆包／黏包」 這句話有語病，其實應該是 TCP 之上的應用應該要定義自己的拆包黏包


func handle(conn net.Conn) {

	defer conn.Close()

	for {

		// 讀 4 bytes
		header := make([]byte, 4)

		_, err := io.ReadFull(
			conn,
			header,
		)

		if err != nil {
			return
		}


		length := binary.BigEndian.Uint32(header)


		// 按長度讀 body
		body := make([]byte, length)

		_, err = io.ReadFull(
			conn,
			body,
		)

		if err != nil {
			return
		}


		fmt.Println(
			"收到:",
			string(body),
		)
	}
}


func main(){

	listener,_ :=
		net.Listen(
			"tcp",
			":9000",
		)


	for {

		conn,_ :=
			listener.Accept()

		go handle(conn)
	}
}