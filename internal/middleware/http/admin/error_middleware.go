package middlewareHttpAdmin

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bootstrap "example/bootstrap"

	pkgUtility "example/pkg/utility"
)

type ErrorMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewErrorMiddleware(oAbstractMiddleware *AbstractMiddleware) *ErrorMiddleware {
	return &ErrorMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// 3. 定義一個方法，返回 gin.HandlerFunc
// 放在中间件链最前面，先 Next() 让后续中间件/handler 执行，
// 执行完毕后统一检查是否有错误，有则返回统一格式的 JSON 响应。
func (oSelf *ErrorMiddleware) Handle() gin.HandlerFunc {

	return func(oContext *gin.Context) {
		oContext.Set("result", struct{}{})         // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("code", 0)                    // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("message", "unknown message") // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("status", 200)                // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("authorization", "")          // 錯誤訊息給一個預設的，避免取 空 報錯
		oContext.Set("total", 0)                   // 錯誤訊息給一個預設的，避免取 空 報錯

		defer func() {
			bError := false

			// 捕获 panic
			if oError := recover(); oError != nil {
				aByteStack := make([]byte, 4096)
				iLen := runtime.Stack(aByteStack, false)

				switch oErrorType := oError.(type) {
				case *pkgUtility.DefaultError: // 需要用 *指標， 因為 controller 是用 指標
					pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Warn(
						"http 業務異常",
						zap.Any("error", oError),
						zap.Any("stack", aByteStack[:iLen]),
					)
					oSelf.Response.Set(oContext, int(oErrorType.Status), int(oErrorType.Code), oErrorType.Message, struct{}{}, 0, "")

				default:
					iLen := runtime.Stack(aByteStack, false)
					// Logger.Fatal 會再觸發 panic

					pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Error(
						"http 系統錯誤",
						zap.Any("error", oError),
						zap.Any("stack", aByteStack[:iLen]),
					)

					oSelf.Response.Set(oContext, 200, -4, "http system error", struct{}{}, 0, "")

				}

				bError = true
			}

			// oContext.Error 只是把 error 存入 context ，其實不會觸發 panic
			// oContext.Abort() : 會強制後面的 handler , middleware 都不執行
			// return 會提前結束 this handler, 如果 已經呼叫了 .Next(), 只後的會執行；如果 尚未呼叫 .Next(), 之後的不會執行

			if len(oContext.Errors) > 0 {
				oLastErr := oContext.Errors.Last()

				aByteStack := make([]byte, 4096)
				iLen := runtime.Stack(aByteStack, false)

				switch oErrorType := oLastErr.Err.(type) {
				case *pkgUtility.DefaultError:
					// 本進程內產生的業務錯誤（沒過 gRPC）
					pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Warn(
						"http 業務異常",
						zap.String("error", oLastErr.Error()),
						zap.Any("stack", aByteStack[:iLen]),
					)
					oSelf.Response.Set(oContext, int(oErrorType.Status), int(oErrorType.Code), oErrorType.Message, struct{}{}, 0, "")

				default:
					// 從 resource gRPC 來的：codes.Aborted = 業務錯誤，其餘 = 系統錯誤
					if oStatus, bOk := status.FromError(oLastErr.Err); bOk && oStatus.Code() == codes.Aborted {
						pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Warn(
							"resource 業務異常",
							zap.String("error", oLastErr.Error()),
						)
						oSelf.Response.Set(oContext, 200, -3, oStatus.Message(), struct{}{}, 0, "")
						break
					}

					if oStatus, bOk := status.FromError(oLastErr.Err); bOk {
						sLog := "resource **異常"
						iCode := -3
						sMessage := oStatus.Message()

						if oStatus.Code() == codes.Aborted {
							sLog = "resource 業務異常"
							iCode = -3
							sMessage = oStatus.Message()

						}

						if oStatus.Code() == codes.Unavailable {
							sLog = "resource 系統錯誤"
							iCode = -4
							sMessage = "resource system error"

						}

						if oStatus.Code() != codes.Aborted && oStatus.Code() != codes.Unavailable {
							sLog = "resource 其他錯誤"
							iCode = -4
							sMessage = "resource other error"

						}

						pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Warn(
							sLog,
							zap.String("error", oLastErr.Error()),
						)

						oSelf.Response.Set(oContext, 200, iCode, sMessage, struct{}{}, 0, "")
						break
					}

					pkgUtility.Logger(pkgUtility.HttpAdminMiddleware).Error(
						"http 系統錯誤",
						zap.String("error", oLastErr.Error()),
						zap.Any("stack", aByteStack[:iLen]),
					)
					oSelf.Response.Set(oContext, 200, -4, "http system error", struct{}{}, 0, "")

				}

				bError = true
			}

			if !bError {
				return
			}

			mStatus, _ := oContext.Get("status")
			mCode, _ := oContext.Get("code")
			mResult, _ := oContext.Get("result")
			mMessage, _ := oContext.Get("message")
			mKey, _ := oContext.Get("key")
			mIv, _ := oContext.Get("iv")

			sKey, _ := mKey.(string) // sKey 可能未定義
			sIv, _ := mIv.(string)

			if sKey == "" {
				sKey = pkgUtility.RandString(16)
			}

			if sIv == "" {
				sIv = pkgUtility.RandString(16)
			}

			sCode := fmt.Sprintf("%d", mCode)
			sMessage := mMessage.(string)
			iStatus := mStatus.(int)

			oKeys := map[string]interface{}{
				"key": sKey,
				"iv":  sIv,
			}
			sKeys, _ := pkgUtility.JsonEncode(oKeys)

			sTime := strconv.FormatInt(oSelf.clock.Now().Unix(), 10)
			sResultJson, _ := pkgUtility.JsonEncode(mResult)

			sR, _ := oSelf.aesHelper.Encrypt(sResultJson, sKey, sIv)
			sC, _ := oSelf.aesHelper.Encrypt(sCode, sKey, sIv)
			sM, _ := oSelf.aesHelper.Encrypt(sMessage, sKey, sIv)
			sT, _ := oSelf.aesHelper.Encrypt("0", sKey, sIv)

			aStrings := []string{sKeys, sTime, sC, sM, sR, bootstrap.CONFIG.SERVICES.HTTP.ADMIN.SALT}
			sHeaderSignature := pkgUtility.Md5(strings.Join(aStrings, ","))

			oJson := gin.H{
				"c": sC,
				"m": sM,
				"r": sR,
				"t": sT,
			}
			if bootstrap.CONFIG.DEFAULT.DEBUG {
				oJson["code"] = mCode
				oJson["message"] = sMessage
				oJson["result"] = mResult
				oJson["total"] = 0

			}
			oContext.Writer.Header().Set("Authorization", "")
			oContext.Writer.Header().Set("Time", sTime)
			oContext.Writer.Header().Set("Signature", sHeaderSignature)

			oContext.JSON(iStatus, oJson)

			// IMPORTANT 如果 ErrorMiddleware 捕獲到錯誤就不再往下
			// IMPORTANT 如果 ErrorMiddleware 捕獲到錯誤就不再往下
			// IMPORTANT 如果 ErrorMiddleware 捕獲到錯誤就不再往下

			oContext.Abort()
		}()

		oContext.Next()
	}
}
