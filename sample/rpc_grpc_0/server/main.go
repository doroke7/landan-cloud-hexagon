package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// 直接导入包
	"grpc_server_stream/protobuf"
	"grpc_server_stream/service"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// 注册服务，使用 monitor 包前缀
	protobuf.RegisterMonitorServiceServer(grpcServer, &service.MonitorService{})

	reflection.Register(grpcServer)

	log.Println("Server processing on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
