package client

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"test/protobuf" // 這裡替換成你生成的套件路徑
)

func main() {
	// 1. 建立連線
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	c := protobuf.NewChatServiceClient(conn)

	// 2. 呼叫雙向流方法
	stream, _ := c.Chat(context.Background())

	waitc := make(chan struct{})

	// 【協程一：接收流】持續監聽伺服器主動噴過來的訊息
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("接收失敗: %v", err)
			}
			log.Printf("📢 收到伺服器廣播: %s", in.Text)
		}
	}()

	// 【協程二：發送流】每 2 秒主動送一個指令給伺服器
	go func() {
		for i := 0; i < 5; i++ {
			msg := &protobuf.Message{
				User: "Architect_User",
				Text: "執行任務指令 #" + string(rune(i+48)),
			}
			if err := stream.Send(msg); err != nil {
				log.Printf("發送失敗: %v", err)
				return
			}
			time.Sleep(2 * time.Second)
		}
		// 傳完後關閉發送方向的流
		stream.CloseSend()
	}()

	<-waitc // 等待接收流結束
}
