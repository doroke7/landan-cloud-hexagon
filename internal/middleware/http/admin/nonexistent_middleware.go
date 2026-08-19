package middleware_admin

import (
	"github.com/gin-gonic/gin"
)

type NonexistentMiddleware struct {
	*AbstractMiddleware
}

// 2. 在結構體上定義一個「構造函數」
func NewNonexistentMiddleware(oAbstractMiddleware *AbstractMiddleware) *NonexistentMiddleware {
	return &NonexistentMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *NonexistentMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

	}
}
