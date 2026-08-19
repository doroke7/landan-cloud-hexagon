package middleware_admin

import (
	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewAdminMiddleware(oAbstractMiddleware *AbstractMiddleware) *AdminMiddleware {

	return &AdminMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}

}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *AdminMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

	}
}
