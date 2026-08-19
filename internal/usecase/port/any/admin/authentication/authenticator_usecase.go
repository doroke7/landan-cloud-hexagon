package any

type AuthenticatorUsecase interface {
	SignIn(name string, password string, secret string) (authorization string, err error)
	Refresh(jwt string, secret string) (authorization string, err error)
}
