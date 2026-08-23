package middlewareWebsocketAdmin

import (
	"encoding/json"
	"strconv"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
	types "example/types"
)

type DecryptionMiddleware struct {
	*AbstractMiddleware
}

func NewDecryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *DecryptionMiddleware {
	return &DecryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 職責跟 http 版本的 DecryptionMiddleware 一樣：從 Header.K 解出這次請求的 AES
// key/iv，把 client 送來的密文 C/V/F/P 解密回明文，分別換成 oReq.Code/oReq.Value/
// oReq.Filters/oReq.Pagination 再往下傳給 handler——handler（例如 SignIn）完全不用知道
// Value 曾經是密文，跟 EncryptionMiddleware 用同一把 Header.K 各自獨立解一次是同一個
// 道理（詳見該檔案註解）。C/V/F/P 分別對應 Code/Value/Filters/Pagination，
// 跟 http 版本的 c/v/f/p 是同一套命名慣例。
//
// 跟 http 版本不同的地方：
//  1. http 版本把 f（filters）、p（pagination）解密後還要攤平進 query／PostForm，給 gin 的
//     query binding 用；websocket 這邊解密完直接整份取代 oReq.Filters/oReq.Pagination 就好，
//     不用 Flatten。
//  2. http 版本沒有 K 時會直接 panic（AES key 長度不對），靠 ErrorMiddleware 接住變成
//     系統錯誤；這裡刻意做成 K 是空的就當作「client 還沒走加密交握」，原樣往下傳
//     （不解密）——這樣在真的有 client 開始帶 key 之前，現有的明文測試流程不會被打壞。
//     K 有值但解不開（金鑰本身壞掉）才是真正的協定錯誤，直接擋掉、不呼叫 fnNext。
//  3. Header 的 authorization（A）解密，等 AuthenticationMiddleware 真正要用的時候
//     再一併處理，這裡先不做。
func (oSelf *DecryptionMiddleware) Handle() types.WebsocketMiddlewareFunc {
	return func(oConn types.WebsocketConn, oReq types.WebsocketRequest, fnNext types.WebsocketNextFunc) types.WebsocketResponse {
		var oHeader signatureHeader
		if len(oReq.Header) > 0 {
			_ = json.Unmarshal(oReq.Header, &oHeader)
		}

		if oHeader.K == "" {
			return fnNext(oConn, oReq)
		}

		sKeys, err := oSelf.rsaHelper.Decrypt(oHeader.K, bootstrap.CONFIG.SERVICES.WEBSOCKET.ADMIN.PRIVATE_KEY)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "金鑰解密失敗"}
		}

		oKeys, err := utility.JsonDecode[struct {
			Key string `json:"key"`
			Iv  string `json:"iv"`
		}](sKeys)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "金鑰解密失敗"}
		}

		sCode, err := oSelf.aesHelper.Decrypt(oReq.C, oKeys.Key, oKeys.Iv)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "內容解密失敗"}
		}
		if sCode != "" {
			iCode, err := strconv.Atoi(sCode)
			if err != nil {
				return types.WebsocketResponse{Code: -1, Message: "內容解密失敗"}
			}
			oReq.Code = iCode
		}

		oByteValue, err := oSelf.aesHelper.Decrypt(oReq.V, oKeys.Key, oKeys.Iv)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "內容解密失敗"}
		}
		oReq.Value = json.RawMessage(oByteValue)

		oByteSearch, err := oSelf.aesHelper.Decrypt(oReq.F, oKeys.Key, oKeys.Iv)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "內容解密失敗"}
		}
		oReq.Filters = json.RawMessage(oByteSearch)

		oBytePagination, err := oSelf.aesHelper.Decrypt(oReq.P, oKeys.Key, oKeys.Iv)
		if err != nil {
			return types.WebsocketResponse{Code: -1, Message: "內容解密失敗"}
		}
		oReq.Pagination = json.RawMessage(oBytePagination)

		return fnNext(oConn, oReq)
	}
}
