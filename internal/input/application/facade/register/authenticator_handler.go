package inputApplicationFacadeRegister

import (
	"context"

	inputApplicationFacade "example/internal/input/application/facade"
	pbFacadeRegister "example/pb/facade/register"
)

type AuthenticatorHandler struct {
	pbFacadeRegister.UnimplementedAuthenticatorControllerServer
	*inputApplicationFacade.AbstractHandler
}

func NewAuthenticatorHandler(oAbstractHandler *inputApplicationFacade.AbstractHandler) *AuthenticatorHandler {
	return &AuthenticatorHandler{
		AbstractHandler: oAbstractHandler,
	}
}

func (oSelf *AuthenticatorHandler) SingUp(oContext context.Context, oRequest *pbFacadeRegister.AuthenticatorSignInRequest) (*pbFacadeRegister.AuthenticatorSignInResponse, error) {

	return &pbFacadeRegister.AuthenticatorSignInResponse{
		Name: "AA",
	}, nil
}
