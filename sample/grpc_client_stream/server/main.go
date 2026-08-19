package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc_client_stream/register" // 引入剛才拆分的 service 包
)

func main() {
	// 1. 開啟 TCP 監聽
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("❌ 無法監聽端口: %v", err)
	}

	// 2. 初始化原生 gRPC Server
	oGrpcServer := grpc.NewServer()

	// 3. 註冊我們的獨立服務
	register.Register(oGrpcServer)

	// 4. (選配) 開啟反射，方便使用 gRPC UI 或 Evans 調試
	reflection.Register(oGrpcServer)

	log.Println("🚀 原生 gRPC Server 啟動於 :8080")
	if err := oGrpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ 伺服器啟動失敗: %v", err)
	}
}
