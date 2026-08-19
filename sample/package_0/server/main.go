package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":9000")
	fmt.Println("[INFO]: 啟動 Server")

	for {
		conn, _ := listener.Accept()

		// 錯誤示範：直接 Read 當一個訊息
		// 錯誤示範：直接 Read 當一個訊息
		// 錯誤示範：直接 Read 當一個訊息
		go func() {
			buf := make([]byte, 1024)

			for {
				n, err := conn.Read(buf)
				if err != nil {
					return
				}

				fmt.Println("收到:", string(buf[:n]))
			}
		}()
	}
}