package inputApplicationHttpGameAuthentication

import (
	"github.com/gin-gonic/gin"

	pkgGin "example/pkg/gin"
	pkgUtility "example/pkg/utility"

	bootstrap "example/bootstrap"
	inputApplicationHttp "example/internal/input/application/http"
	usecasePortAnyGameAuthentication "example/internal/usecase/port/any/game/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationHttp.AbstractHandler
	AuthenticatorUsecase usecasePortAnyGameAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAuthenticatorUsecase usecasePortAnyGameAuthentication.AuthenticatorUsecase) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) LogIn(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &struct {
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}{}

	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		_ = oContext.Error(pkgUtility.NewDefaultError("pagination 格式錯誤", -1, 200))
		return
	}

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.LogIn(
		*oValue.Name,
		*oValue.Password,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "AppUser 成功登入", struct{}{}, 0, sAuthorization, nil)

}

func (oSelf *AuthenticatorHandler) Refresh(oContext *gin.Context) {

	sJwt := oContext.GetHeader("Authorization")

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.Refresh(
		sJwt,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "成功刷新", struct{}{}, 0, sAuthorization, nil)

}
