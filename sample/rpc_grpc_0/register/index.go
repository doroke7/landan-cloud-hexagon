package register

import (
	"grpc_client_stream/protobuf"
	"grpc_client_stream/service"

	"google.golang.org/grpc"
)

func Register(oServer *grpc.Server) {
	protobuf.RegisterMonitorServer(oServer, &service.MonitorService{})
}
