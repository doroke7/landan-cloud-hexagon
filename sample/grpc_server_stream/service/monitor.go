package service

import (
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// 不再使用 pb 别名
	"grpc_server_stream/protobuf"
)

type MonitorService struct {
	// 直接使用生成的包名 monitor
	protobuf.UnimplementedMonitorServiceServer
}

func (s *MonitorService) WatchMetrics(req *protobuf.MonitorRequest, stream protobuf.MonitorService_WatchMetricsServer) error {
	resource := req.GetResourceName()
	if resource == "" {
		return status.Error(codes.InvalidArgument, "resource_name is required")
	}

	log.Printf("Starting monitor for: %s", resource)

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Monitor stopped for: %s", resource)
			return stream.Context().Err()
		default:
			res := &protobuf.MetricResponse{
				Value:     rand.Float32() * 100,
				Timestamp: time.Now().Unix(),
			}

			if err := stream.Send(res); err != nil {
				return err
			}

			time.Sleep(time.Second)
		}
	}
}
