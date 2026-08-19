package service

import (
	"context"
	"time"

	pb "heartbeater.com.server/pb"
)

// HeartbeaterService 實作 pb.HeartbeaterServer 定義的 Ping 方法
type HeartbeaterService struct {
	pb.UnimplementedHeartbeaterServer
}

func (s *HeartbeaterService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{
		Status:     "OK",
		ServerTime: time.Now().Unix(),
	}, nil
}
