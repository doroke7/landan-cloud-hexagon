package authentication

import (
	inputApplicationCommand "example/internal/input/application/command"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
)

type AuthenticatorHandler struct {
	*inputApplicationCommand.AbstractHandler
	AuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationCommand.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:      oAbstractHandler,
		AuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

// SignIn 只負責轉呼叫 usecase，不知道自己是被 CLI 呼叫的，也不知道要用哪個 secret——
// cobra.Command 的組裝、container 的建立時機、JWT secret 要用哪個 carrier 的設定，
// 都交給 cmd/register 那層決定。
func (oSelf *AuthenticatorHandler) SignIn(sName string, sPassword string, sSecret string) (string, error) {
	return oSelf.AuthenticatorUsecase.SignIn(sName, sPassword, sSecret)
}
