package interceptor_facade_admin

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

	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {

		// before
		fmt.Println("Before StatusInterceptor...")

		oResponse, oReq := fnHandler(oContext, oRequest)
		fmt.Println("After StatusInterceptor...")

		// after

		return oResponse, oReq
	}

}
