package usecasePortAnyAdminAuthentication

type AuthenticatorUsecase interface {
	SignIn(sName string, sPassword string, sSecret string) (sAuthorization string, oError error)
	Refresh(sJwt string, sSecret string) (sAuthorization string, oError error)
}
