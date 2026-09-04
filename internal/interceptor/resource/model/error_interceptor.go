package interceptorResourceModel

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

				pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Error(
					"resource 系統錯誤",
					zap.Any("panic", oPanic),
					zap.String("method", oServerInfo.FullMethod),
					zap.String("stack", string(oByteStack[:iLen])),
				)

				oResponse = nil
				oErr = status.Error(codes.Internal, "resource 系統錯誤")
			}
		}()

		oResponse, oErr = fnHandler(oContext, oRequest)

		if oErr == nil {
			return oResponse, nil
		}

		// 1. handler 直接回的 *DefaultError（業務錯誤）
		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {
			pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Warn(
				"resource 業務異常",
				zap.String("method", oServerInfo.FullMethod),
				zap.String("message", oDefaultError.Message),
				zap.Int16("code", oDefaultError.Code),
			)

			return nil, status.Error(codes.Aborted, oDefaultError.Message)
		}

		// 2. handler 用 status.Error 回的（也當業務錯誤）
		if oStatus, bOk := status.FromError(oErr); bOk {
			pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Warn(
				"resource 業務異常",
				zap.String("method", oServerInfo.FullMethod),
				zap.String("message", oStatus.Message()),
				zap.Int32("code", int32(oStatus.Code())),
			)

			return nil, status.Error(codes.Aborted, oStatus.Message())
		}

		// 3. 其他原生錯誤（MYSQL / Redis 等）
		pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Error(
			"resource 系統錯誤",
			zap.String("method", oServerInfo.FullMethod),
			zap.Error(oErr),
		)

		return nil, status.Error(codes.Unavailable, "resource 系統錯誤")
	}

}
