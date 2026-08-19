package tcp

import (
	helper "example/internal/helper"
)

// AbstractHandler 放 tcp 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / http / consumer）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*helper.AesHelper
}

func NewAbstractHandler(oAesHelper *helper.AesHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
	}
}
