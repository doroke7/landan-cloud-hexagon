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
	GameAuthenticationAuthenticatorUsecase usecasePortAnyGameAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationHttp.AbstractHandler, oAuthenticatorUsecase usecasePortAnyGameAuthentication.AuthenticatorUsecase) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:                        oAbstractHandler,
		GameAuthenticationAuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) LogIn(oContext *gin.Context) {

	oRequest := &pkgGin.Request{Context: oContext}

	oValue := &struct {
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}{}

	if oErr := oRequest.Bind("variable", oValue); oErr != nil {
		oDefaultError := pkgUtility.NewDefaultError("pagination format error", -1, 200)
		_ = oContext.Error(oDefaultError)
		return
	}

	sAuthorization, oErr := oSelf.GameAuthenticationAuthenticatorUsecase.LogIn(
		*oValue.Name,
		*oValue.Password,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "AppUser signed in successfully", struct{}{}, 0, sAuthorization)

}

func (oSelf *AuthenticatorHandler) Refresh(oContext *gin.Context) {

	sJwt := oContext.GetHeader("Authorization")

	sAuthorization, oErr := oSelf.GameAuthenticationAuthenticatorUsecase.Refresh(
		sJwt,
		bootstrap.CONFIG.SERVICES.HTTP.GAME.JWT.SECRET,
	)

	if oErr != nil {
		_ = oContext.Error(oErr)
		return
	}

	oSelf.Response.Set(oContext, 200, 1, "Refreshed successfully", struct{}{}, 0, sAuthorization)

}
