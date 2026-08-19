package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_client_stream/protobuf"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("連線失敗: %v", err)
	}
	defer conn.Close()

	client := protobuf.NewMonitorClient(conn)

	// 1. 調用方法獲取 stream 對象
	stream, err := client.PushData(context.Background())
	if err != nil {
		log.Fatalf("開啟串流失敗: %v", err)
	}

	// 2. 模擬連續推送 5 筆數據
	devices := []string{"Node-A", "Node-B", "Node-C", "Node-D", "Node-E"}
	for _, name := range devices {
		req := &protobuf.DataRequest{
			DeviceId: name,
			Payload:  "Heartbeat Status: Normal",
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("發送數據失敗: %v", err)
		}
		log.Printf("📤 已推送: %s", name)
		time.Sleep(500 * time.Millisecond) // 模擬發送間隔
	}

	// 3. 告訴伺服器「我傳完了」，並接收回傳結果
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("接收伺服器響應失敗: %v", err)
	}

	log.Printf("🏁 伺服器回傳：共處理 %d 筆數據，狀態: %s", res.Count)
}
