package middlewareHttpAdmin

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	pkg "example/pkg"
)

type AuthenticationMiddleware struct {
	*AbstractMiddleware
}

// go的嵌入式繼承（組合繼承） 比較特殊， Abstract 類別 需要注入到子類別，這個其他語言不需要這個動作

// 2. 在結構體上定義一個「構造函數」
func NewAuthenticationMiddleware(oAbstractMiddleware *AbstractMiddleware) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		AbstractMiddleware: oAbstractMiddleware,
	}
}

// Handle 驗證 DecryptionMiddleware 解密後放進 context 的 JWT（key: "Authrization"）。
// 沒帶 token 就當作不需要驗證、原樣往下傳——這樣 SignIn 這種「用來換 token」的
// method 才不會卡在雞生蛋蛋生雞的問題；要不要對特定 method 強制要求 token，
// 之後再視需求擴充。（跟 websocket 版 AuthenticationMiddleware 是同一套邏輯）
func (oSelf *AuthenticationMiddleware) Handle() gin.HandlerFunc {
	return func(oContext *gin.Context) {

		if bootstrap.CONFIG.SERVICES.HTTP.ADMIN.AUTHENTICATION != true {
			oContext.Next()
			return
		}

		sAuthorization := oContext.GetHeader("Authorization")
		sJwt := sAuthorization

		if sJwt == "" {
			_ = oContext.Error(pkg.NewDefaultError("缺少 jwt", -2, 200))
			oContext.Abort()
			return
		}

		if _, oErr := oSelf.jwtHelper.Parse(sJwt); oErr != nil {
			_ = oContext.Error(pkg.NewDefaultError("身份驗證失敗", -2, 200))
			oContext.Abort()
			return
		}

		oContext.Next()

	}
}
