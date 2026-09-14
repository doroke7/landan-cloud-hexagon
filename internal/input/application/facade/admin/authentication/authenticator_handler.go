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

	oVariable := oRequest.GetVariable()
	sName := oVariable.GetName()
	sPassword := oVariable.GetPassword()

	sAuthorization, oErr := oSelf.AdminAuthenticationAuthenticatorUsecase.SignIn(sName, sPassword, bootstrap.CONFIG.SERVICES.FACADE.ADMIN.JWT.SECRET)

	if oErr != nil {
		if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {
			sDefaultError := oDefaultError.Error()
			oStatusError := status.Error(codes.Aborted, sDefaultError)

			return nil, oStatusError
		}

		sError := oErr.Error()
		oStatusError := status.Error(codes.Internal, sError)

		return nil, oStatusError
	}

	if oAuthorization, bOk := oContext.Value("a").(*string); bOk {
		*oAuthorization = sAuthorization
	}

	return &pbFacadeAdminAuthentication.AuthenticatorSignInResponse{}, nil
}
