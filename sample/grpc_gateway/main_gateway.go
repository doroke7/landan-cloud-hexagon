package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"project/pb"
)

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// 設定轉接到核心的連線資訊
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 將產生的適配代碼註冊到 mux
	_ = pb.RegisterTranslatorHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)

	println("Adapter 外殼啟動於 :8080 (HTTP)")
	http.ListenAndServe(":8080", mux)
}
