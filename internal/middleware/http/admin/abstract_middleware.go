package middleware_admin

import (
	"github.com/gin-gonic/gin"

	pkg "example/pkg"

	helper "example/internal/helper"
)

type AbstractMiddleware struct {
	*pkg.Response
	clock     *pkg.Clock
	rsaHelper *helper.RsaHelper
	aesHelper *helper.AesHelper
	jwtHelper *helper.JwtHelper
}

// 2. 在結構體上定義一個「構造函數」
func NewAbstractMiddleware(oResponse *pkg.Response, oClock *pkg.Clock, oRsaHelper *helper.RsaHelper, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractMiddleware {
	return &AbstractMiddleware{
		Response:  oResponse,
		clock:     oClock,
		rsaHelper: oRsaHelper,
		aesHelper: oAesHelper,
		jwtHelper: oJwtHelper,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *AbstractMiddleware) HandleAbstractMiddleware() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

	}
}
