package register

import (
	pkg "example/pkg"

	container "example/container"
)

func TcpInit(oContainer *container.TcpContainer) *pkg.TcpRouter {
	oTcpRouter := pkg.NewTcpRouter()
	oTcpRouter.HandleFunc("Admin.Authentication.Authenticator.SignIn", oContainer.TcpAdminAuthenticationSignIn.SignIn)

	return oTcpRouter
}
