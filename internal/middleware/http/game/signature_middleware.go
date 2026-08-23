package middleware_game

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	pkg "example/pkg"
	types "example/types"

	bootstrap "example/bootstrap"
	utility "example/internal/utility"
)

type SignatureMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewSignatureMiddleware(oAbstractMiddleware *AbstractMiddleware) *SignatureMiddleware {
	return &SignatureMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *SignatureMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		sVer := oContext.GetHeader("Ver")
		sVersion := oContext.GetHeader("Version")
		sK := oContext.GetHeader("K")
		sTime := oContext.GetHeader("Time")
		sHeaderSignature := oContext.GetHeader("Signature")

		// 故意不用 oContext.DefaultQuery/PostForm 讀這幾個「加密前」的原始值——那兩個
		// 方法會觸發 gin 的 queryCache/formCache 快取，之後 DecryptionMiddleware 改寫
		// RawQuery/PostForm 時就得反過來清快取。直接讀 c.Request.URL.Query()（每次都是
		// 現剖析 RawQuery，不會被 gin 快取）就完全不會有這個問題。
		sF := oContext.Request.URL.Query().Get("f")
		sS := oContext.Request.URL.Query().Get("s")
		sP := oContext.Request.URL.Query().Get("p")

		var oRequestPayload types.HttpRequestBody

		_ = oContext.ShouldBindBodyWith(&oRequestPayload, binding.JSON)

		sV := oRequestPayload.V
		// NOTE: 不要把 未加密的 search, option, param, 都加下去簽名，多次一舉
		aStrings := []string{sVer, sVersion, sK, sTime, sF, sS, sP, sV, bootstrap.CONFIG.SERVICES.HTTP.GAME.SALT}
		sStrings := strings.Join(aStrings, "|")
		sMd5Signature := utility.Md5(sStrings)

		if bootstrap.CONFIG.SERVICES.HTTP.GAME.SIGNATURE == true {
			if sMd5Signature != sHeaderSignature {

				/*
					1. panic 可以用嗎？ -> 可以用，但是不建議
					   1-1. 使用 panic 可以達到業務效果，並且流程簡單
					   1-2. 但是不建議用panic 實現，panic 性能較差

					2. 不用 panic 需要這樣搭配 -> 只能用這個了，但是很麻煩
					   2-1. 寫上 return 代表 這個 middleware 終止
					   2-2. 寫上 oContext.Abort() 代表 後面的 handler 跟 middleware 都不執行
					   2-3. 寫上 oContext.Error() 代表寫入 gin.Error, 在  ErrorMiddleware 的 defer 後收集 error 數據後寫入響應

					3. 也不能用 oContext.AbortWithError
					   3-1. 這個會 Abort + Error + 寫入 http-status
					   3-2. 但是 會使得 ErrorMiddleware 無法修改 http 響應了
				*/
				_ = oContext.Error(pkg.NewDefaultError("簽名失敗", -3, 406))
				oContext.Abort()
				return
			}
		}

		oContext.Next()

	}
}
