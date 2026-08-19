package middleware_admin

import (
	types "example/types"
)

type AdminMiddleware struct {
	*AbstractMiddleware
}

func NewAdminMiddleware(oAbstractMiddleware *AbstractMiddleware) *AdminMiddleware {
	return &AdminMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 跟 http 版本的 AdminMiddleware 一樣是骨架，先讓請求往下傳。
func (oSelf *AdminMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		return fnNext(oConn, oReq)
	}
}
