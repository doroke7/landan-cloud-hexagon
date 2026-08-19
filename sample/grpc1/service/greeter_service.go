package service

import (
	"context"
	"hello/pb"
	"log"
)

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func (oSelf *GreeterService) SayHello(ctx context.Context, oHellRequest *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("收到客戶端請求: %v", oHellRequest.GetName())
	return &pb.HelloReply{Message: "你好 " + oHellRequest.GetName()}, nil
}
