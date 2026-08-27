package registerTcp

import (
	pkg "example/pkg"

	container "example/container"
)

func Init(oContainer *container.TcpContainer) *pkg.TcpRouter {
	oTcpRouter := pkg.NewTcpRouter()
	oTcpRouter.HandleFunc("Admin.Authentication.Authenticator.SignIn", oContainer.TcpAdminAuthenticationSignIn.SignIn)

	return oTcpRouter
}
