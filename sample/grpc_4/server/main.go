package server // golang 的一個目錄 裡面的 go文件 全部package要同名

import (
	"io"
	"log"
	"net"
	"time"

	"test/protobuf" // 這裡替換成你生成後的套件路徑

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedChatServiceServer
}

// Chat 實作了雙向流邏輯
func (s *server) Chat(stream protobuf.ChatService_ChatServer) error {
	// 用於協調退出
	errChan := make(chan error)

	// 1. 【非同步寫入】：伺服器主動、持續地發送訊息
	go func() {
		for {
			msg := &protobuf.Message{
				User: "Server",
				Text: "心跳檢查: " + time.Now().Format("15:04:05"),
			}
			if err := stream.Send(msg); err != nil {
				errChan <- err
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// 2. 【同步讀取】：持續監聽客戶端傳來的指標數據
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				errChan <- nil
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			log.Printf("收到來自 [%s] 的消息: %s", in.User, in.Text)
		}
	}()

	return <-errChan
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	protobuf.RegisterChatServiceServer(s, &server{})
	log.Println("gRPC 全雙工伺服器啟動於 :50051")
	s.Serve(lis)
}
