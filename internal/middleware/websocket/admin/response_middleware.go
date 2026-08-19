package middleware_admin

import (
	"encoding/json"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
	types "example/types"
)

type ResponseMiddleware struct {
	*AbstractMiddleware
}

func NewResponseMiddleware(oAbstractMiddleware *AbstractMiddleware) *ResponseMiddleware {
	return &ResponseMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 ResponseMiddleware 一樣：production 只送密文（C/M/R），
// 明文（Code/Message/Result）跟 Header.Authorization 只在 DEBUG 模式下方便直接看，
// 不是 DEBUG 就清空，不會被送上線。ResponseMiddleware 註冊在 EncryptionMiddleware
// 外層（詳見 register/websocket.go 的順序），fnNext 回來的時候 Encryption 已經跑完、
// oResp.Header/C/M/R 都已經是加密後的內容。
func (oSelf *ResponseMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		oResp := fnNext(oConn, oReq)

		if bootstrap.CONFIG.DEFAULT.DEBUG {
			return oResp
		}

		var oHeader encryptionHeader
		if len(oResp.Header) > 0 {
			_ = json.Unmarshal(oResp.Header, &oHeader)
		}
		oHeader.Authorization = ""

		sHeader, _ := utility.JsonEncode(oHeader)
		oResp.Header = json.RawMessage(sHeader)

		oResp.Code = 0
		oResp.Message = ""
		oResp.Result = nil

		return oResp
	}
}
