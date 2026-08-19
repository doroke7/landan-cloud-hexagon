package main

import (
	"fmt"
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":12345")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	fmt.Println("安全 UDP 服務端已啟動...")

	// 接收緩衝區稍微開大一點（例如 65535 最大 UDP 限制），防止底層截斷
	buf := make([]byte, 65535)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		// 解包
		packet, err := Unpack(buf[:n])
		if err != nil {
			fmt.Printf("來自 %s 的無效封包: %v\n", remoteAddr, err)
			continue
		}

		fmt.Printf("[收到封包 #%d] 長度: %d, 內容: %s\n", packet.SeqID, packet.Len, string(packet.Data))
	}
}
