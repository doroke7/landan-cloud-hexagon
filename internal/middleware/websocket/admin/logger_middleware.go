package middleware_admin

import (
	"time"

	"go.uber.org/zap"

	pkg "example/pkg"
	types "example/types"
)

type LoggerMiddleware struct {
	*AbstractMiddleware
}

func NewLoggerMiddleware(oAbstractMiddleware *AbstractMiddleware) *LoggerMiddleware {
	return &LoggerMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 LoggerMiddleware 一樣：記錄進入/結束的時間，
// 不同的是 http 記 path/query/headers，這裡改記這筆訊息的 method/param。
func (oSelf *LoggerMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		oStart := time.Now()

		pkg.Logger(pkg.WebsocketAdminMiddleware).Info(
			"進入 websocket",
			zap.String("method", oReq.Method),
			zap.ByteString("param", oReq.Param),
		)

		oResp := fnNext(oConn, oReq)

		pkg.Logger(pkg.WebsocketAdminMiddleware).Info(
			"結束 websocket",
			zap.String("method", oReq.Method),
			zap.Duration("經過時間", time.Since(oStart)),
		)

		return oResp
	}
}
