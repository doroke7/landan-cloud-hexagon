package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"

	bootstrap "example/bootstrap"
)

var oUdpCommand = &cobra.Command{
	Use:   "udp",
	Short: "啟動 UDP 服務",
	Run: func(cmd *cobra.Command, args []string) {

		addr, err := net.ResolveUDPAddr("udp", ":"+bootstrap.CONFIG.SERVICES.UDP.PORT)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		fmt.Println("UDP Server Start")

		buf := make([]byte, 1024)

		/*
		   這邊的 for 不會一直空轉
		   跟 go channel 很接近

		*/
		for {

			// 這裡會 很像 event 的機制， 一直讀取，直到沒有資料了就 sleep（阻塞）， 直到 有消息時候會喚醒程式碼

			n, remoteAddr, err := conn.ReadFromUDP(buf) // 會 （阻塞）
			if err != nil {
				continue
			}

			fmt.Println("Receive:", string(buf[:n]))

			conn.WriteToUDP([]byte("OK"), remoteAddr)
		}
	},
}

func init() {
	// 將 server 指令加入到 root 中
	oRootCommand.AddCommand(oUdpCommand)
}
