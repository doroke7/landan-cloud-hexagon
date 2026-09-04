package inputApplicationHttpAdminAuthentication

import (
	"github.com/gin-gonic/gin"

	bootstrap "example/bootstrap"
	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pkgGin "example/pkg/gin"
	pkgUtility "example/pkg/utility"
)

type AuthenticatorHandler struct {
	*inputApplicationHttp.AbstractHandler
	AdminAuthenticationAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

// NewUserHandler 構造函數 (Go 的慣用法)，
// 相当 PHP 的 __construct()

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:                         oAbstractHandler,
		AdminAuthenticationAuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext *gin.Context) {

	// 模擬 遠程長時間 ，方便前端做效果
	// time.Sleep(2 * time.Second)

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &struct {
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}{}

	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("pagination 格式錯誤", -1, 200))
		return
	}

	sAuthorization, oErr := oSelf.AdminAuthenticationAuthenticatorUsecase.SignIn(
		*oValue.Name,
		*oValue.Password,
		bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "AdminUser 成功登入", struct{}{}, 0, sAuthorization)

}

func (oSelf *AuthenticatorHandler) Refresh(oContext *gin.Context) {

	// 模擬 遠程長時間 ，方便前端做效果
	// time.Sleep(2 * time.Second)

	sJwt := oContext.GetHeader("Authorization")

	// Authorization
	sAuthorization, oErr := oSelf.AdminAuthenticationAuthenticatorUsecase.Refresh(
		sJwt,
		bootstrap.CONFIG.SERVICES.HTTP.ADMIN.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "成功刷新", struct{}{}, 0, sAuthorization)

}
