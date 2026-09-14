package inputApplicationFacadeAdminAuthentication

import (
	"context"

	bootstrap "example/bootstrap"
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
)

type AuthenticatorHandler struct {
	pbFacadeAdminAuthentication.UnimplementedAuthenticatorControllerServer
	*inputApplicationFacade.AbstractHandler
	AdminAuthenticationAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase
}

func NewAuthenticatorHandler(oAuthenticatorUsecase usecasePortAnyAdminAuthentication.AuthenticatorUsecase, oAbstractHandler *inputApplicationFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler:                         oAbstractHandler,
		AdminAuthenticationAuthenticatorUsecase: oAuthenticatorUsecase,
	}
}

func (oSelf *AuthenticatorHandler) SignIn(oContext context.Context, oRequest *pbFacadeAdminAuthentication.AuthenticatorSignInRequest) (*pbFacadeAdminAuthentication.AuthenticatorSignInResponse, error) {

	oVariable := oRequest.GetVariable()
	sName := oVariable.GetName()
	sPassword := oVariable.GetPassword()

	sAuthorization, oErr := oSelf.AdminAuthenticationAuthenticatorUsecase.SignIn(sName, sPassword, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.SECRET)

	if oErr != nil {
		return nil, oErr
	}

	if oAuthorization, bOk := oContext.Value("a").(*string); bOk {
		*oAuthorization = sAuthorization
	}

	return &pbFacadeAdminAuthentication.AuthenticatorSignInResponse{}, nil
}
