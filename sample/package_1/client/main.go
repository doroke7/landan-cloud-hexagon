package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

func send(conn net.Conn, msg string) {

	var buf bytes.Buffer

	// 寫入長度
	binary.Write(
		&buf,
		binary.BigEndian,
		uint32(len(msg)),
	)

	// 寫入資料
	buf.WriteString(msg)

	conn.Write(buf.Bytes())
}


func main() {

	conn, _ := net.Dial(
		"tcp",
		"localhost:9000",
	)

	send(conn, "Hello")
	send(conn, "World")
	send(conn, "Go")
}