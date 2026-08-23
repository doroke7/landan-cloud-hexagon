package middlewareHttpTable

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	bootstrap "example/bootstrap"
	pkg "example/pkg"
)

type ResponseMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewResponseMiddleware(oAbstractMiddleware *AbstractMiddleware) *ResponseMiddleware {
	return &ResponseMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
func (oSelf *ResponseMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		oContext.Next()

		mStatus, _ := oContext.Get("status")
		mCode, _ := oContext.Get("code")
		mMessage, _ := oContext.Get("message")
		mResult, _ := oContext.Get("result")
		mAuthorization, _ := oContext.Get("authorization")
		mTotal, _ := oContext.Get("total")

		mC, _ := oContext.Get("c")
		mM, _ := oContext.Get("m")
		mR, _ := oContext.Get("r")
		mA, _ := oContext.Get("a")
		mT, _ := oContext.Get("t")

		iCode := mCode.(int)
		sMessage := mMessage.(string)
		sAuthorization := mAuthorization.(string)
		iTotal := mTotal.(int)

		sA := mA.(string)
		sC := mC.(string)
		sM := mM.(string)
		sR := mR.(string)
		sT := mT.(string)

		iStatus := mStatus.(int)

		if iCode < 0 {
			var iLevel zapcore.Level = 1
			if iCode == -1 {
				iLevel = 1
			}

			if iCode <= -2 {
				iLevel = 2
			}
			aByteStack := make([]byte, 4096)
			iLen := runtime.Stack(aByteStack, false)
			pkg.Logger(pkg.HttpTableMiddleware).Log(
				iLevel,
				"前級系統錯誤4",
				zap.Any("message", sMessage),
				zap.Any("stack", aByteStack[:iLen]),
			)
		}

		if !oContext.Writer.Written() {
			oJson := gin.H{
				"c": sC,
				"m": sM,
				"r": sR,
				"t": sT,
			}

			if bootstrap.CONFIG.DEFAULT.DEBUG {
				oJson["code"] = iCode
				oJson["message"] = sMessage
				oJson["result"] = mResult
				oJson["total"] = iTotal

				oContext.Writer.Header().Set("Authorization", sAuthorization)

			}
			oContext.Writer.Header().Set("A", sA)
			oContext.JSON(iStatus, oJson)
		}

	}
}
