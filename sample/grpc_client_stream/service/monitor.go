package service

import (
	"context"
	"grpc_client_stream/protobuf" // 替換為你的實際 module 名稱
	"io"
	"log"
	"strings"
)

// MonitorService 獨立負責處理監控相關業務
type MonitorService struct {
	protobuf.UnimplementedMonitorServer
}

// PushData 實現 Client Streaming 邏輯
func (s *MonitorService) PushData(stream protobuf.Monitor_PushDataServer) error {
	var count int32
	log.Println("🛰️ [Service] 收到新的 Client Streaming 連線")

	for {
		// 接收客戶端傳來的數據
		req, err := stream.Recv()

		// 判斷是否傳輸結束
		if err == io.EOF {
			log.Printf("✅ [Service] 接收完成，總計處理 %d 筆數據", count)
			return stream.SendAndClose(&protobuf.DataResponse{
				Count: count,
			})
		}
		if err != nil {
			log.Printf("❌ [Service] 串流接收異常: %v", err)
			return err
		}

		// 實際業務處理邏輯 (例如寫入 DB)
		count++
		log.Printf("📥 [Service] 正在處理設備: %s, 負載: %s", req.DeviceId, req.Payload)
	}
}

func (s *MonitorService) Repeat(ctx context.Context, oReq *protobuf.DataRequest1) (*protobuf.DataResponse1, error) {
	sInput := oReq.GetInput()

	log.Printf("收到請求: Input=%s ", sInput)

	// 使用 strings.Repeat 處理邏輯
	sResult := strings.Repeat(sInput, 10)

	return &protobuf.DataResponse1{Result: sResult}, nil
}
