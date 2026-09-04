package interceptorFacadeAdmin

import (
	"context"
	"fmt"
	"runtime"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pkgUtility "example/pkg/utility"
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

	return func(oContext context.Context, oReqeust any, oServerIno *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		oErr = nil

		fmt.Println("Before ErrorInterceptor...")

		defer func() {

			if oPanic := recover(); oPanic != nil {

				oByteStack := make([]byte, 4096)
				iLen := runtime.Stack(oByteStack, false)

				// Logger.Fatal 會再觸發 panic
				pkgUtility.Logger(pkgUtility.FacadeAdminInterceptor).Error(
					"upstream system error",
					zap.Any("panic", oPanic),
					zap.String("method", oServerIno.FullMethod),
					zap.String("stack", string(oByteStack[:iLen])),
				)

				oResponse = nil
				oErr = status.Error(codes.Unavailable, "upstream system error")
			}
		}()

		oResponse, oErr = fnHandler(oContext, oReqeust)
		fmt.Println("After ErrorInterceptor...")

		// after

		return oResponse, oErr
	}

}
