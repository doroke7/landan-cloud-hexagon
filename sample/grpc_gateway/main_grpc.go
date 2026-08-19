package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"project/pb"
)

type server struct {
	pb.UnimplementedTranslatorServer
}

func (s *server) Translate(ctx context.Context, in *pb.TransRequest) (*pb.TransReply, error) {
	return &pb.TransReply{TranslatedText: "核心處理中: " + in.Text}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterTranslatorServer(s, &server{})
	log.Println("gRPC 核心啟動於 :50051")
	s.Serve(lis)
}
