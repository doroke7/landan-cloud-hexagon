package interceptorResourceModel

import (
	"context"

	"google.golang.org/grpc"
)

type AbstractInterceptor struct {
}

func NewAbstractInterceptor() *AbstractInterceptor {
	return &AbstractInterceptor{}
}

func (oSelf *AbstractInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(ctx, req)
	}
}
