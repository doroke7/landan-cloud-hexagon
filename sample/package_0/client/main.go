package main

import (
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:9000")

	conn.Write([]byte("Hello"))
	conn.Write([]byte("World"))
	conn.Write([]byte("Go"))

	time.Sleep(time.Second)
}