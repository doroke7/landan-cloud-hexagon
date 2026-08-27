package registerTcp

import (
	pkgTcp "example/pkg/tcp"

	container "example/container"
)

func Init(oContainer *container.TcpContainer) *pkgTcp.TcpRouter {
	oTcpRouter := pkgTcp.NewTcpRouter()
	oTcpRouter.HandleFunc("Admin.Authentication.Authenticator.SignIn", oContainer.TcpAdminAuthenticationSignIn.SignIn)

	return oTcpRouter
}
