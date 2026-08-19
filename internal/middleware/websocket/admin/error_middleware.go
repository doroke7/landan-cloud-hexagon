package middleware_admin

import (
	pkg "example/pkg"
	types "example/types"
)

type ErrorMiddleware struct {
	*AbstractMiddleware
}

func NewErrorMiddleware(oAbstractMiddleware *AbstractMiddleware) *ErrorMiddleware {
	return &ErrorMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 ErrorMiddleware 一樣：放在鏈最前面，統一攔截下游 panic
// （包含 *pkg.DefaultError 這種業務錯誤），轉成統一格式的 types.WebsocketResponse，
// 個別 handler 就不用每個自己寫 recover；不同的是 http 版本要寫 gin.Context 的
// header/JSON，這裡直接把結果當函式回傳值交出去就好，不用額外的 Writer。
func (oSelf *ErrorMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) (oResp types.WebsocketResponse) {
		defer func() {
			oError := recover()
			if oError == nil {
				return
			}

			switch oErrorType := oError.(type) {
			case *pkg.DefaultError:
				oResp = types.WebsocketResponse{Code: int(oErrorType.Code), Message: oErrorType.Message}
			default:
				pkg.Logger(pkg.WebsocketMiddleware).Sugar().Errorf("panic: %v", oError)
				oResp = types.WebsocketResponse{Code: -4, Message: "系統錯誤"}
			}
		}()

		return fnNext(oConn, oReq)
	}
}
