package websocket

import (
	helper "example/internal/helper"
)

// AbstractHandler 放 websocket 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / http / consumer / cron）的抽象類各自獨立，互不共用。
//
// upgrade／心跳／優雅關機這些跟連線生命週期有關的技術細節都收在 pkg.WebsocketRouter 裡，
// 業務 handler（例如 authenticator_handler.go）不用碰，簽名直接對應 pkg.WebsocketHandlerFunc。
type AbstractHandler struct {
	*helper.AesHelper
}

func NewAbstractHandler(oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
	}
}
