package usecasePortAnyGameAuthentication

type AuthenticatorUsecase interface {
	LogIn(sName string, sPassword string, sSecret string) (sAuthorization string, oError error)
	Refresh(sJwt string, sSecret string) (sAuthorization string, oError error)
}
