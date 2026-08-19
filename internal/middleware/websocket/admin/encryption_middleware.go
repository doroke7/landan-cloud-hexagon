package middleware_admin

import (
	"encoding/json"
	"fmt"
	"strings"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
	types "example/types"
)

// encryptionHeader 是 Encryption 回應時寫回 oResp.Header 的中繼資料，對應 http 版本
// 回應時寫進 Time/Signature/Authorization 幾個 HTTP Header 的做法；Authorization 是否
// 留在回應裡（還是被清空）由 ResponseMiddleware 依 DEBUG 模式決定。
type encryptionHeader struct {
	Time          string `json:"time"`
	Signature     string `json:"signature"`
	Authorization string `json:"authorization"`
}

type EncryptionMiddleware struct {
	*AbstractMiddleware
}

func NewEncryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *EncryptionMiddleware {
	return &EncryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 EncryptionMiddleware 一樣：等 fnNext 跑完拿到明文結果，
// 用這次請求帶來的 AES key/iv 把 Code/Message/Result 分別加密後放進 oResp.C/M/R，
// 再簽一次 md5 放進回應的 Header——跟 oReq.P/oReq.Q 是同一套慣例，C/M/R 才是真正送上線
// 的密文欄位，Code/Message/Result 保留明文只是方便 debug 直接看。
//
// 跟 http 版本不同的地方：
//  1. http 版本靠 gin.Context.Set/Get 從 DecryptionMiddleware 那邊拿 key/iv；我們的
//     middleware chain 是純函式組合、沒有共用的可變狀態，所以這裡直接從 oReq.Header.K
//     自己重新解一次（跟 DecryptionMiddleware 之後要做的事一樣，只是各自獨立算，不依賴
//     彼此的執行順序——請求還在，Header 不會變，重算一次不影響正確性）。
//  2. http 版本另外用固定的 JWT key/iv 加密 authorization；websocket 這邊還沒有「每次
//     回應都帶著」的 authorization 概念，先不處理。
//  3. Header.K 是空的或解不開（client 還沒走加密交握）就原樣回傳，不 fail——這樣在
//     Decryption 真正實作、client 真的開始帶 key 之前，這個 middleware 不會擋掉任何請求。
func (oSelf *EncryptionMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		oResp := fnNext(oConn, oReq)

		var oHeader signatureHeader
		if len(oReq.Header) > 0 {
			_ = json.Unmarshal(oReq.Header, &oHeader)
		}

		sKeys, err := oSelf.rsaHelper.Decrypt(oHeader.K, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.PRIVATE_KEY)
		if err != nil {
			return oResp
		}

		oKeys, err := utility.JsonDecode[struct {
			Key string `json:"key"`
			Iv  string `json:"iv"`
		}](sKeys)
		if err != nil {
			return oResp
		}

		sC, err := oSelf.aesHelper.Encrypt(fmt.Sprintf("%d", oResp.Code), oKeys.Key, oKeys.Iv)
		if err != nil {
			return oResp
		}
		sM, err := oSelf.aesHelper.Encrypt(oResp.Message, oKeys.Key, oKeys.Iv)
		if err != nil {
			return oResp
		}
		sR, err := oSelf.aesHelper.Encrypt(string(oResp.Result), oKeys.Key, oKeys.Iv)
		if err != nil {
			return oResp
		}

		sTime := utility.Time[string](false)
		aStrings := []string{sTime, sC, sM, sR, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.SALT}
		sSignature := utility.Md5(strings.Join(aStrings, ","))

		sHeader, _ := utility.JsonEncode(encryptionHeader{Time: sTime, Signature: sSignature})

		oResp.Header = json.RawMessage(sHeader)
		oResp.C = sC
		oResp.M = sM
		oResp.R = sR

		return oResp
	}
}
