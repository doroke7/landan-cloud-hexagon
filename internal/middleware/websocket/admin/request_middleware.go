package middleware_admin

import (
	types "example/types"
)

type RequestMiddleware struct {
	*AbstractMiddleware
}

func NewRequestMiddleware(oAbstractMiddleware *AbstractMiddleware) *RequestMiddleware {
	return &RequestMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 對應 http 版本的 RequestMiddleware：http 版本存在的理由是 gin 的 query binding
// 沒辦法直接綁巢狀 JSON，所以 debug 模式下要把 search/option 的巢狀 JSON 攤平成點記法
// 的 query key 才能被 controller 讀到。websocket 沒有 query string 這個限制，oReq.Param
// 本來就是一份完整的 JSON，各 handler 自己 json.Unmarshal 就能拿到巢狀結構，不需要攤平——
// 這裡的空實作是刻意維持的最終狀態，不是還沒補上的骨架。
func (oSelf *RequestMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		return fnNext(oConn, oReq)
	}
}
