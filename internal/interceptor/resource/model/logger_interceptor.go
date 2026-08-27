package interceptorResourceModel

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
)

type LoggerInterceptor struct {
	*AbstractInterceptor
}

func NewLoggerInterceptor(oInterceptor *AbstractInterceptor) *LoggerInterceptor {
	return &LoggerInterceptor{
		AbstractInterceptor: oInterceptor,
	}
}

func (oSelf *LoggerInterceptor) Handle() grpc.UnaryServerInterceptor {

	return func(oContext context.Context, oRequest any, oServerInfo *grpc.UnaryServerInfo, fnHandler grpc.UnaryHandler) (any, error) {

		iTime1 := pkgUtility.Time[int](true)

		if bootstrap.CONFIG.LOGGERS.INTERCEPTOR.STATUS {
			pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Info(
				"進入 grpc",
				zap.String("method", oServerInfo.FullMethod),
				zap.Any("request", oRequest),
			)
		}

		oResponse, oErr := fnHandler(oContext, oRequest)
		iTime2 := pkgUtility.Time[int](true)

		if bootstrap.CONFIG.LOGGERS.INTERCEPTOR.STATUS {
			pkgUtility.Logger(pkgUtility.ResourceModelInterceptor).Info(
				"結束 grpc",
				zap.String("method", oServerInfo.FullMethod),
				zap.Any("response", oResponse),
				zap.Error(oErr),
				zap.Int("time1(進入時間ms)", iTime1),
				zap.Int("time2(結束時間ms)", iTime2),
				zap.Int("time2-time1(經過時間ms)", iTime2-iTime1),
			)
		}

		return oResponse, oErr
	}

}
