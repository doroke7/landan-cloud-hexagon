package middleware_admin

import (
	types "example/types"
)

type NonexistentMiddleware struct {
	*AbstractMiddleware
}

func NewNonexistentMiddleware(oAbstractMiddleware *AbstractMiddleware) *NonexistentMiddleware {
	return &NonexistentMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 跟 http 版本不同：http 版本的 NonexistentMiddleware 本來就沒有真正的邏輯，
// 也沒被接進 adminMiddlewares() 的鏈裡，是死代碼；這裡改接進 pkg.WebsocketRouter 的
// NoMethod，變成 method 真的找不到時會被呼叫的 hook——fnNext 是內建的
// ErrWebsocketMethodNotFound 回應，目前原樣回傳，之後想額外記錄未知 method、
// 或換一套對外訊息，都可以直接在這裡加，不用動 pkg.WebsocketRouter。
func (oSelf *NonexistentMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		return fnNext(oConn, oReq)
	}
}
