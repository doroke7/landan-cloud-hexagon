package middlewareHttpGame

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
	types "example/types"
)

type DecryptionMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewDecryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *DecryptionMiddleware {
	return &DecryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *DecryptionMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		// sHeaderKeys := oContext.GetHeader("Keys")
		sHeaderK := oContext.GetHeader("K")
		sHeaderA := oContext.GetHeader("A")

		sQueryF := oContext.Request.URL.Query().Get("f")
		sQueryS := oContext.Request.URL.Query().Get("s")
		sQueryP := oContext.Request.URL.Query().Get("p")

		var oRequestBody types.HttpRequestBody

		_ = oContext.ShouldBindBodyWith(&oRequestBody, binding.JSON)

		sV := oRequestBody.V

		//
		sKeys, oErr := oSelf.rsaHelper.Decrypt(sHeaderK, bootstrap.CONFIG.SERVICES.HTTP.GAME.PRIVATE_KEY)
		if oErr != nil {
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
			_ = oContext.Error(pkgUtility.NewDefaultError("金鑰解密失敗", -1, 400))
			oContext.Abort()
			return
		}

		oKeys, _ := pkgUtility.JsonDecode[struct {
			Key string `json:"key"`
			Iv  string `json:"iv"`
		}](sKeys)

		oContext.Set("key", oKeys.Key)
		oContext.Set("iv", oKeys.Iv)

		sPagination, oErr := oSelf.aesHelper.Decrypt(sQueryP, oKeys.Key, oKeys.Iv)
		if oErr != nil {
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
			_ = oContext.Error(pkgUtility.NewDefaultError("pagination 解密失敗", -1, 400))
			oContext.Abort()

			return
		}

		sSorters, oErr := oSelf.aesHelper.Decrypt(sQueryS, oKeys.Key, oKeys.Iv)
		if oErr != nil {
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
			_ = oContext.Error(pkgUtility.NewDefaultError("sorter 解密失敗", -1, 400))
			oContext.Abort()

			return
		}

		sFilters, oErr := oSelf.aesHelper.Decrypt(sQueryF, oKeys.Key, oKeys.Iv)
		if oErr != nil {
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
			_ = oContext.Error(pkgUtility.NewDefaultError("filter 解密失敗", -1, 400))
			oContext.Abort()

			return
		}

		sValue, oErr := oSelf.aesHelper.Decrypt(sV, oKeys.Key, oKeys.Iv)
		if oErr != nil {
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
			_ = oContext.Error(pkgUtility.NewDefaultError("value 解密失敗", -1, 400))
			oContext.Abort()

			return
		}

		sAuthorizaion, oErr := oSelf.aesHelper.Decrypt(sHeaderA, bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.KEY, bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.IV)
		if oErr != nil {
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

			_ = oContext.Error(pkgUtility.NewDefaultError("Authorization 解密失敗", -1, 400))
			oContext.Abort()

			return
		}

		oUrlQuery := oContext.Request.URL.Query()

		oUrlQuery.Set("filters", sFilters)
		oUrlQuery.Set("sorters", sSorters)
		oUrlQuery.Set("pagination", sPagination)

		oContext.Request.URL.RawQuery = oUrlQuery.Encode()
		oContext.Request.PostForm = url.Values{"value": []string{sValue}}

		oContext.Request.Header.Add(
			"Authorization",
			sAuthorizaion,
		)

		oContext.Next()

	}
}
