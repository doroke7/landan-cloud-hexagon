package interceptorResourceLogic

import (
	"context"
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

	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (oResponse any, oErr error) {

		defer func() {

			if oPanic := recover(); oPanic != nil {

				oByteStack := make([]byte, 4096)
				iLen := runtime.Stack(oByteStack, false)

				// Logger.Fatal 會再觸發 panic

				pkgUtility.Logger(pkgUtility.ResourceLogicInterceptor).Error(
					"resource 系統錯誤",
					zap.Any("panic", oPanic),
					zap.String("method", oServerInfo.FullMethod),
					zap.String("stack", string(oByteStack[:iLen])),
				)

				oResponse = nil
				oErr = status.Error(codes.Internal, "resource system error")
			}
		}()

		oResponse, oErr = fnHandler(oContext, oRequest)

		if oErr != nil {

			oStatus, bOk := status.FromError(oErr)

			// 如果是 gRPC error 就打印 並且回傳 error 到前級 Http或Facade
			if bOk {
				pkgUtility.Logger(pkgUtility.ResourceLogicInterceptor).Error(
					"resource 業務異常",
					zap.String("method", oServerInfo.FullMethod),
					zap.String("message", oStatus.Message()),
					zap.Int32("code", int32(oStatus.Code())),
				)
				return nil, status.Error(codes.Aborted, oStatus.Message())

			}
			// 如果非 gRPC error 就打印 並且回傳 “內部錯誤” 到前級 Http或Facade
			// 例如 MYSQL 錯誤。Redis 錯誤
			if !bOk {
				pkgUtility.Logger(pkgUtility.ResourceLogicInterceptor).Error(
					"resource 系統錯誤",
					zap.String("method", oServerInfo.FullMethod),
					zap.Error(oErr),
				)

				return nil, status.Error(codes.Unavailable, "resource system error")

			}

		}

		return oResponse, oErr
	}

}
