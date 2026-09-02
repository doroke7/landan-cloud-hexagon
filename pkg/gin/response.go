package pkgGin

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (oSelf *Response) Set(oContext *gin.Context, iStatus int, iCode int, sMessage string, oResult any, iTotal int, sAuthorization string, oErr error) {
	oContext.Set("code", iCode)
	oContext.Set("message", sMessage)
	oContext.Set("status", iStatus)
	oContext.Set("result", oResult)
	oContext.Set("total", iTotal)
	oContext.Set("authorization", sAuthorization)

	// oErr 可選，不需要時傳 nil；帶進來時存到 context 給 logger / error middleware 取用。
	if oErr != nil {
		oContext.Set("error", oErr)
	}
}
