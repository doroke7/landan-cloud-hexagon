package middleware_admin

import (
	helper "example/internal/helper"
)

// AbstractMiddleware 放 websocket admin middleware 共用依賴，跟 http 版本的 AbstractMiddleware
// 是同一個概念：rsaHelper 解出 Header.K 帶的 AES key/iv，aesHelper 拿那把 key/iv 加解密
// Result，jwtHelper 驗證 A 解密後的 JWT。Signature/Decryption/Encryption/Authentication
// 都已經有真正的邏輯。
type AbstractMiddleware struct {
	rsaHelper *helper.RsaHelper
	aesHelper *helper.AesHelper
	jwtHelper *helper.JwtHelper
}

func NewAbstractMiddleware(oRsaHelper *helper.RsaHelper, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractMiddleware {
	return &AbstractMiddleware{
		rsaHelper: oRsaHelper,
		aesHelper: oAesHelper,
		jwtHelper: oJwtHelper,
	}
}
