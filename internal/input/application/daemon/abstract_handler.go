package daemon

import (
	helper "example/internal/helper"
)

// AbstractHandler 放 daemon 這個 input adapter 自己專用的共用依賴。
// gRPC client（SourceClient）不放這裡——開 stream 這件事交給 register 決定，
// handler 只管收到資料之後怎麼處理。
type AbstractHandler struct {
	*helper.AesHelper
}

func NewAbstractHandler(oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
	}
}
