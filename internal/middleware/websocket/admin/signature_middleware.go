package middleware_admin

import (
	"encoding/json"
	"strings"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
	types "example/types"
)

// signatureHeader 對應 oReq.Header，跟 http 版本用的 Ver/Version/K/Time/Signature/A
// 這幾個 HTTP Header 是同一套慣例，只是 websocket 沒有 HTTP header，改用 Header 這個
// json.RawMessage 欄位承載；A 是給 AuthenticationMiddleware 用的加密 authorization。
type signatureHeader struct {
	Ver       string `json:"ver"`
	Version   string `json:"version"`
	K         string `json:"k"`
	A         string `json:"a"`
	Time      string `json:"time"`
	Signature string `json:"signature"`
}

type SignatureMiddleware struct {
	*AbstractMiddleware
}

func NewSignatureMiddleware(oAbstractMiddleware *AbstractMiddleware) *SignatureMiddleware {
	return &SignatureMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 SignatureMiddleware 一樣：把 Ver/Version/K/Time 跟這筆訊息
// 加密前的 s/o/p 三個欄位（對應 oReq.S/oReq.O/oReq.P）串起來算 md5，比對 Header 裡的
// Signature，簽名對不上就直接擋掉，不呼叫 fnNext；是否真的要擋用
// bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN 這組獨立設定，不再跟 http 共用。
func (oSelf *SignatureMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		var oHeader signatureHeader
		if len(oReq.Header) > 0 {
			if err := json.Unmarshal(oReq.Header, &oHeader); err != nil {
				return types.WebsocketResponse{Code: -1, Message: "invalid header, expect {\"ver\":..., \"version\":..., \"k\":..., \"time\":..., \"signature\":...}"}
			}
		}

		// NOTE: 簽的是 oReq.S/oReq.O/oReq.P 這三個「加密前」的密文欄位，不是
		// DecryptionMiddleware 解密後才會有值的 oReq.Search/oReq.Option/oReq.Param——
		// Signature 跑在 Decryption 之前，這時候 Search/Option/Param 都還是零值，簽那些
		// 等於什麼都沒簽到；跟 http 版本簽 s/o/p 是同一個道理：一律簽「加密前」的內容，
		// 不要簽解密後的明文，多此一舉。
		aStrings := []string{oHeader.Ver, oHeader.Version, oHeader.K, oHeader.Time, oReq.S, oReq.O, oReq.P, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.SALT}
		sMd5Signature := utility.Md5(strings.Join(aStrings, "|"))

		if bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.SIGNATURE && sMd5Signature != oHeader.Signature {
			return types.WebsocketResponse{Code: -3, Message: "簽名失敗"}
		}

		return fnNext(oConn, oReq)
	}
}
