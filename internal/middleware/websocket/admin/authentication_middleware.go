package middleware_admin

import (
	"encoding/json"

	bootstrap "example/bootstrap"
	types "example/types"
)

type AuthenticationMiddleware struct {
	*AbstractMiddleware
}

func NewAuthenticationMiddleware(oAbstractMiddleware *AbstractMiddleware) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 是 http 版本沒有真正邏輯可以參照的一個（http 的 AuthenticationMiddleware 本身
// 也只是 oContext.Next()，沒被接進 adminMiddlewares() 的鏈裡）；這裡是重新設計的行為：
// 把 Header.A 用固定的 JWT key/iv（跟 per-request 的 RSA key 不同一把）解密回 JWT 字串，
// 再交給 helper.JwtHelper.Parse 驗證簽章跟過期時間，驗證失敗就直接擋掉、不呼叫 fnNext。
//
// A 是空的（沒帶 token）就當作不需要驗證、原樣往下傳——這樣 SignIn 這種「用來換 token」
// 的 method 才不會卡在雞生蛋蛋生雞的問題；要不要對特定 method 強制要求 token，
// 之後再視需求擴充。
func (oSelf *AuthenticationMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {

		if bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.AUTHENTICATION != true {
			return fnNext(oConn, oReq)
		}

		var oHeader signatureHeader
		if len(oReq.Header) > 0 {
			_ = json.Unmarshal(oReq.Header, &oHeader)
		}

		if oHeader.A == "" {
			return fnNext(oConn, oReq)
		}

		sJwt, err := oSelf.aesHelper.Decrypt(oHeader.A, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.JWT.KEY, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.JWT.IV)
		if err != nil {
			return types.WebsocketResponse{Code: -2, Message: "身份驗證失敗"}
		}

		if _, err := oSelf.jwtHelper.Parse(sJwt); err != nil {
			return types.WebsocketResponse{Code: -2, Message: "身份驗證失敗"}
		}

		return fnNext(oConn, oReq)
	}
}
