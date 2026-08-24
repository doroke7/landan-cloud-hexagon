package register

import (
	container "example/container"
	types "example/types"
)

func websocketAdminMiddlewares(oContainer *container.WebsocketContainer) []types.WebsocketMiddlewareFunc {
	return []types.WebsocketMiddlewareFunc{
		oContainer.WebsocketAdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.WebsocketAdminErrorMiddleware.Handle(),
		oContainer.WebsocketAdminSignatureMiddleware.Handle(),
		oContainer.WebsocketAdminDecryptionMiddleware.Handle(),
		oContainer.WebsocketAdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.WebsocketAdminResponseMiddleware.Handle(),
		oContainer.WebsocketAdminEncryptionMiddleware.Handle(),
	}
}

// websocketChain 手動把 middleware 由外而內包住 fnHandler，最後一個 middleware
// 包最裡層、離 handler 最近，跟 pkg.WebsocketRouter 內部 chain() 的順序一致。
func websocketChain(fnHandler types.WebsocketNextFunc, aMiddlewares []types.WebsocketMiddlewareFunc) types.WebsocketNextFunc {
	for i := len(aMiddlewares) - 1; i >= 0; i-- {
		fnMiddleware := aMiddlewares[i]
		fnNext := fnHandler
		fnHandler = func(oConn types.WebsocketConn, oReq types.WebsocketRequest) types.WebsocketResponse {
			return fnMiddleware(oConn, oReq, fnNext)
		}
	}
	return fnHandler
}

// WebsocketInit 只回傳 method -> handler 對照表，不碰 upgrade/ping/dispatch 這些
// websocket 協定細節，那些是 cmd/websocket.go 用 pkg.WebsocketRouter 組裝的事。
func WebsocketInit(oContainer *container.WebsocketContainer) map[string]types.WebsocketNextFunc {
	aMiddlewares := websocketAdminMiddlewares(oContainer)

	return map[string]types.WebsocketNextFunc{
		"Admin/Authentication/Authenticator.SignIn": websocketChain(
			oContainer.WebsocketAdminAuthenticationAuthenticator.SignIn,
			aMiddlewares,
		),
	}
}
