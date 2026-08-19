package interceptor_facade_game

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type ErrorInterceptor struct {
	*AbstractInterceptor
}

func NewErrorInterceptor(oInterceptor *AbstractInterceptor) *ErrorInterceptor {
	return &ErrorInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *ErrorInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(
		oContext context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// before

		fmt.Println("Before ErrorInterceptor...")

		oCtx, oReq := handler(oContext, req)
		fmt.Println("After ErrorInterceptor...")

		// after

		return oCtx, oReq
	}

}
