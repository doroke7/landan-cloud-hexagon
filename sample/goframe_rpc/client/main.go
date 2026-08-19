package main

import (
	"fmt"
	"goframe_rpc/protobuf"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		ctx            = gctx.New()
		conn           = grpcx.Client.MustNewGrpcClientConn("127.0.0.1:8888")
		oMonitorClient = protobuf.NewMonitorClient(conn)
	)
	// 2. 開啟串流
	stream, err := oMonitorClient.PushData(ctx)
	if err != nil {
		// 處理錯誤
	}
	aStrings := []string{"data-1", "data-2", "data-3", "data-4", "data-5"}

	// 3. 持續推送資料
	for _, sString := range aStrings {
		err = stream.Send(&protobuf.DataRequest{
			DeviceId: "device-ABC",
			Payload:  sString,
		})
		if err != nil {
			// 處理錯誤
		}
	}

	// 4. 關閉串流，取得 server 回傳的統計結果
	resp, err := stream.CloseAndRecv()
	fmt.Println("共收到筆數:", resp.Count)
}
