package middlewareHttpGame

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgUtility "example/pkg/utility"

	helper "example/internal/helper"
)

type AbstractMiddleware struct {
	*pkgGin.Response
	clock     *pkgUtility.Clock
	rsaHelper *helper.RsaHelper
	aesHelper *helper.AesHelper
	jwtHelper *helper.JwtHelper
}

// 2. 在結構體上定義一個「構造函數」
func NewAbstractMiddleware(oResponse *pkgGin.Response, oClock *pkgUtility.Clock, oRsaHelper *helper.RsaHelper, oAesHelper *helper.AesHelper, oJwtHelper *helper.JwtHelper) *AbstractMiddleware {
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
