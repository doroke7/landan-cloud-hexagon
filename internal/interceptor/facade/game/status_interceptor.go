package interceptor_facade_game

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type StatusInterceptor struct {
	*AbstractInterceptor
}

func NewStatusInterceptor(oInterceptor *AbstractInterceptor) *StatusInterceptor {
	return &StatusInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *StatusInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(
		oContext context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// before
		fmt.Println("Before StatusInterceptor...")

		oCtx, oReq := handler(oContext, req)
		fmt.Println("After StatusInterceptor...")

		// after

		return oCtx, oReq
	}

}
