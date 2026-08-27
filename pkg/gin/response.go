package pkgGin

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
}

func NewResponse() *Response {
	return &Response{}
}

func (oSelf *Response) Set(oContext *gin.Context, iStatus int, iCode int, sMessage string, oResult any, iTotal int, sAuthorization string) {
	oContext.Set("code", iCode)
	oContext.Set("message", sMessage)
	oContext.Set("status", iStatus)
	oContext.Set("result", oResult)
	oContext.Set("authorization", sAuthorization)
	oContext.Set("total", iTotal)

}
