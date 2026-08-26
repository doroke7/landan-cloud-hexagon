package middlewareWebsocketAdmin

import (
	helper "example/internal/helper"
	pkg "example/pkg"
)

// AbstractMiddleware 放 websocket admin middleware 共用依賴，跟 http 版本的 AbstractMiddleware
// 是同一個概念。原本 Signature/Decryption/Encryption/Authentication 幾個 middleware 都掛
// 在這裡共用 rsaHelper/aesHelper/jwtHelper，但那套設計是在沒考慮 event／method 路由的情況下
// 寫的統一鏈式 middleware，跟現在 WebsocketEventer 依 event 分派的架構對不上，已經整批移除；
// 這幾個 helper 欄位先留著，之後要重新設計「依 event/method 路由」的 middleware 時再用。
type AbstractMiddleware struct {
	clock     *pkg.Clock
	rsaHelper *helper.RsaHelper
	aesHelper *helper.AesHelper
	jwtHelper *helper.JwtHelper
}

func NewAbstractMiddleware(oClock *pkg.Clock, oRsaHelper *helper.RsaHelper, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractMiddleware {
	return &AbstractMiddleware{
		clock:     oClock,
		rsaHelper: oRsaHelper,
		aesHelper: oAesHelper,
		jwtHelper: oJwtHelper,
	}
}
