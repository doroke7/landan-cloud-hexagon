package middlewareHttpAdmin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	pkgUtility "example/pkg/utility"
)

type EncryptionMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewEncryptionMiddleware(oAbstractMiddleware *AbstractMiddleware) *EncryptionMiddleware {
	return &EncryptionMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *EncryptionMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

		mKey, _ := oContext.Get("key")
		mIv, _ := oContext.Get("iv")

		mCode, _ := oContext.Get("code")
		mResult, _ := oContext.Get("result")
		mMessage, _ := oContext.Get("message")
		mAuthorization, _ := oContext.Get("authorization")
		mTotal, _ := oContext.Get("total")

		sCode := fmt.Sprintf("%d", mCode)
		sKey := fmt.Sprintf("%s", mKey)
		sIv := fmt.Sprintf("%s", mIv)
		sAuthorization := fmt.Sprintf("%s", mAuthorization)
		sTotal := fmt.Sprintf("%d", mTotal)

		sMessage := mMessage.(string)

		sTime := strconv.FormatInt(oSelf.clock.Now().Unix(), 10)
		sResult, _ := pkgUtility.JsonEncode(mResult)

		sR, oErr := oSelf.aesHelper.Encrypt(sResult, sKey, sIv)
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
			_ = oContext.Error(pkgUtility.NewDefaultError("加密失敗", -4, 500))
			oContext.Abort()

			return
		}
		sC, oErr := oSelf.aesHelper.Encrypt(sCode, sKey, sIv)
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
			_ = oContext.Error(pkgUtility.NewDefaultError("加密失敗", -4, 500))
			oContext.Abort()

			return
		}
		sM, oErr := oSelf.aesHelper.Encrypt(sMessage, sKey, sIv)
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
			_ = oContext.Error(pkgUtility.NewDefaultError("加密失敗", -4, 500))
			oContext.Abort()

			return
		}
		sT, oErr := oSelf.aesHelper.Encrypt(sTotal, sKey, sIv)
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
			_ = oContext.Error(pkgUtility.NewDefaultError("加密失敗", -4, 500))
			oContext.Abort()

			return
		}

		sA := ""
		if sAuthorization != "" {
			sA, oErr = oSelf.aesHelper.Encrypt(sAuthorization, bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.KEY, bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.IV)
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
				_ = oContext.Error(pkgUtility.NewDefaultError("加密失敗", -4, 500))
				oContext.Abort()

				return
			}
		}

		// NOTE: 不要把 未加密的 code, message, result 都加下去簽名，多次一舉
		aStrings := []string{sTime, sC, sM, sR, bootstrap.CONFIG.SERVICES.HTTP.ADMIN.SALT}

		sString := strings.Join(aStrings, ",")

		sHeaderSignature := pkgUtility.Md5(sString)

		oContext.Writer.Header().Set("Time", sTime)
		oContext.Writer.Header().Set("Signature", sHeaderSignature)

		oContext.Set("a", sA)
		oContext.Set("c", sC)
		oContext.Set("m", sM)
		oContext.Set("r", sR)
		oContext.Set("t", sT)

	}
}
