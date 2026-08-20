package interceptor_resource_model

import (
	"context"

	"google.golang.org/grpc"

	pkg "example/pkg"
)

type AbstractInterceptor struct {
	Clock *pkg.Clock
}

func NewAbstractInterceptor(oClock *pkg.Clock) *AbstractInterceptor {
	return &AbstractInterceptor{
		Clock: oClock,
	}
}

func (oSelf *AbstractInterceptor) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return handler(ctx, req)
	}
}
