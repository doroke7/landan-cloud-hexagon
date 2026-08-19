package service

import (
	"github.com/gogf/gf/v2/frame/g"
	"goframe_rpc/protobuf"
	"io"
)

// 定義 Service 實體
type Monitor struct {
	protobuf.UnimplementedMonitorServer
}

// 導出一個變數方便外部調用

func NewMonitor() *Monitor {
	return &Monitor{}
}

// 實作推送邏輯
func (s *Monitor) PushData(stream protobuf.Monitor_PushDataServer) error {
	var count int32
	ctx := stream.Context()

	g.Log().Info(ctx, "🚀 [Service] gRPC 串流訂閱啟動...")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 結束串流並回傳總數
			g.Log().Info(ctx, "✅ [Service] 客戶端推送完成")
			return stream.SendAndClose(&protobuf.DataResponse{Count: count})
		}
		if err != nil {
			g.Log().Errorf(ctx, "❌ [Service] 接收出錯: %v", err)
			return err
		}

		count++
		// 這裡可以寫入資料庫或進行業務計算
		g.Log().Infof(ctx, "📥 [Service] 處理數據 - 設備: %s, 內容: %s", req.DeviceId, req.Payload)
	}
}
