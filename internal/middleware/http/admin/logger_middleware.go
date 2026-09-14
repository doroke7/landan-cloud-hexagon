package middlewareHttpAdmin

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
)

type LoggerMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewLoggerMiddleware(oAbstractMiddleware *AbstractMiddleware) *LoggerMiddleware {
	return &LoggerMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *LoggerMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		iTime1 := pkgUtility.Time[int](true)
		sPath := oContext.Request.URL.Path
		sRawQuery := oContext.Request.URL.RawQuery
		oMapHeaders := oContext.Request.Header

		if bootstrap.CONFIG.LOGGERS.MIDDLEWARE.STATUS {
			oLogger := pkgUtility.Logger(pkgUtility.HttpAdminMiddleware)
			oLogger.Info(
				"進入 http",
				zap.String("path", sPath),
				zap.String("query", sRawQuery),
				zap.Any("headers", oMapHeaders),
			)
		}

		oContext.Next()
		iTime2 := pkgUtility.Time[int](true)

		if bootstrap.CONFIG.LOGGERS.MIDDLEWARE.STATUS {
			oLogger := pkgUtility.Logger(pkgUtility.HttpAdminMiddleware)
			oLogger.Info(
				"結束 http",
				zap.String("path", sPath),
				zap.String("query", sRawQuery),
				zap.Any("headers", oMapHeaders),
				zap.Int("time1(進入時間ms)", iTime1),
				zap.Int("time2(結束時間ms)", iTime2),
				zap.Int("time2-time1(經過時間ms)", iTime2-iTime1),
			)
		}

	}
}
