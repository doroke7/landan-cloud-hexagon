package inputApplicationFacadeAdminAuthentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bootstrap "example/bootstrap"
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminAuthentication "example/internal/usecase/port/any/admin/authentication"
	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
	pkgUtility "example/pkg/utility"
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

	sAuthorization, oErr := oSelf.AdminAuthenticationAuthenticatorUsecase.SignIn(oRequest.Name, oRequest.Password, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.SECRET)

	if oErr != nil {
		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {
			return nil, status.Error(codes.Aborted, oDefaultError.Error())
		}

		return nil, status.Error(codes.Internal, oErr.Error())
	}

	if oAuthorization, bOk := oContext.Value("a").(*string); bOk {
		*oAuthorization = sAuthorization
	}

	return &pbFacadeAdminAuthentication.AuthenticatorSignInResponse{}, nil
}
