package main

import (
	"context"
	"fmt"
	"log"

	genGoApp "example.rpc.thrift/gen-go/app" // 引入自動生成的包
	"github.com/apache/thrift/lib/go/thrift"
)

func main() {

	// 1. 建立 Transport（TCP 連線），要跟 server 端監聽的位址一致
	oSocket := thrift.NewTSocketConf("localhost:9090", nil)

	// 2. 設定傳輸快取與編碼格式，必須跟 server 端一致
	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	oTransport, err := transportFactory.GetTransport(oSocket)
	if err != nil {
		log.Fatalln("建立 Transport 失敗:", err)
	}
	defer oTransport.Close()

	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

	if err := oTransport.Open(); err != nil {
		log.Fatalln("連線開啟失敗:", err)
	}

	iprot := protocolFactory.GetProtocol(oTransport)
	oprot := protocolFactory.GetProtocol(oTransport)

	// 3. 建立 Client Stub
	oHelloClient := genGoApp.NewHelloClientProtocol(oTransport, iprot, oprot)

	// 4. 呼叫遠端方法
	sRsp, err := oHelloClient.TestHello(context.Background(), "Tom")
	if err != nil {
		log.Fatalln("呼叫失敗:", err)
	}

	fmt.Println("收到回應:", sRsp)
}
