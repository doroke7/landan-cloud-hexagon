package middleware_game

import (
	"github.com/gin-gonic/gin"
)

type RequestMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewRequestMiddleware(oAbstractMiddleware *AbstractMiddleware) *RequestMiddleware {
	return &RequestMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
//
// debug 模式下允許開發者直接傳未加密的 ?search={...}&option={...}（跳過 AES/RSA），
// 這已經是 handler 直接想要的形狀（單一 query 參數塞一包 JSON 字串），
// 不需要再額外處理，這裡純粹是通過。
func (oSelf *RequestMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {
		oContext.Next()
	}
}
