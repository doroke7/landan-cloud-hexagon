package authentication

import (
	"github.com/gin-gonic/gin"

	pkg "example/pkg"

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

	oRequest := &pkg.Request{Context: oContext}

	oValue := &struct {
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}{}

	if oErr := oRequest.Bind("value", oValue); oErr != nil {
		oSelf.Response.Set(oContext, 200, -1, "pagination 格式錯誤", struct{}{}, 0, "")
		return
	}

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.LogIn(
		*oValue.Name,
		*oValue.Password,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {
			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Message, struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, 0, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "成功登入", struct{}{}, 0, sAuthorization)

}

func (oSelf *AuthenticatorHandler) Refresh(oContext *gin.Context) {

	sJwt := oContext.GetHeader("Authorization")

	sAuthorization, oErr := oSelf.AuthenticatorUsecase.Refresh(
		sJwt,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {

		if oDefaultError, bOk := oErr.(*pkg.DefaultError); bOk {
			oSelf.Response.Set(oContext, 200, int(oDefaultError.Code), oDefaultError.Message, struct{}{}, 0, "")
			return
		}

		oSelf.Response.Set(oContext, 200, -3, oErr.Error(), struct{}{}, 0, "")
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "成功刷新", struct{}{}, 0, sAuthorization)

}
