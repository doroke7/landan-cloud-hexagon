package main

import (
	"hello/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"hello/pb" // 換成你的專案路徑
)

func main() {
	oListen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("無法監聽: %v", err)
	}

	oServer := grpc.NewServer()
	pb.RegisterGreeterServer(oServer, service.NewGreeterService())

	log.Println("gRPC 伺服器啟動於 :50051...")
	if err := oServer.Serve(oListen); err != nil {
		log.Fatalf("無法啟動服務: %v", err)
	}
}
