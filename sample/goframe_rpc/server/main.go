package main

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"goframe_rpc/internal/service"
	"goframe_rpc/protobuf"
)

func main() {
	// 1. 初始化 GoFrame gRPC Server
	oServer := grpcx.Server.New(&grpcx.GrpcServerConfig{
		Address: ":8888", // 改這裡
	})

	// 2. 註冊服務（將我們定義的 MonitorService 掛載上去）
	protobuf.RegisterMonitorServer(oServer.Server, service.NewMonitor())

	// 3. 啟動伺服器 (預設端口為 8000)
	oServer.Run()
}
