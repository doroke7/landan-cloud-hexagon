package handler

import (
	pkg "example/pkg"

	helper "example/internal/helper"
)

// AbstractHandler 放 http 這個 input adapter 自己專用的共用依賴，
// 跟其他 input adapter（client / grpc / consumer）的抽象類各自獨立，互不共用。
type AbstractHandler struct {
	*pkg.Response
	*helper.AesHelper
	*helper.JwtHelper
}

func NewAbstractHandler(oResponse *pkg.Response, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractHandler {
	return &AbstractHandler{
		AesHelper: oAesHelper,
		JwtHelper: oJwtHelper,
		Response:  oResponse,
	}
}
