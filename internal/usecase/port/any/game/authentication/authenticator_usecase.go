package any

type AuthenticatorUsecase interface {
	LogIn(name string, password string, secret string) (authorization string, err error)
	Refresh(jwt string, secret string) (authorization string, err error)
}
