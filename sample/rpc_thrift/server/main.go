package main

import (
	"fmt"
	"log"

	genGoApp "example.rpc.thrift/gen-go/app" // 引入自動生成的包
	service "example.rpc.thrift/service"
	"github.com/apache/thrift/lib/go/thrift"
)

func main() {

	oHelloService := service.NewHello()

	oHelloProcessor := genGoApp.NewHelloProcessor(oHelloService)

	// 4. 設定監聽的 Port (Transport 層)
	serverTransport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		log.Fatalln("監聽失敗:", err)
	}

	// 5. 設定傳輸快取與編碼格式 (Protocol 層)
	// 這裡使用經典的 Binary (二進位) 協議
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	// 6. 建立並啟動 Thrift Server
	server := thrift.NewTSimpleServer4(
		oHelloProcessor,
		serverTransport,
		transportFactory,
		protocolFactory,
	)

	fmt.Println("Thrift 服務端已啟動，監聽 port 9090...")
	if err := server.Serve(); err != nil {
		log.Fatalln("啟動失敗:", err)
	}
}
