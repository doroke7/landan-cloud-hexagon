package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hello/pb"
)

func main() {
	// 建立連線
	oConnection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("無法連線: %v", err)
	}
	defer oConnection.Close()

	oClient := pb.NewGreeterClient(oConnection)

	// 呼叫服務
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	oResponse, err := oClient.SayHello(ctx, &pb.HelloRequest{Name: "Gemini"})
	if err != nil {
		log.Fatalf("呼叫失敗: %v", err)
	}
	log.Printf("伺服器回傳: %s", oResponse.GetMessage())
}
