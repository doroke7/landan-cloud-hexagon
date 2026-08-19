package register

import (
	container "example/container"
	pkg "example/pkg"
	"example/types"
)

func websocketAdminMiddlewares(oContainer *container.WebsocketContainer) []types.WebsocketMiddlewareFunc {
	return []types.WebsocketMiddlewareFunc{
		oContainer.WebsocketAdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.WebsocketAdminErrorMiddleware.Handle(),
		oContainer.WebsocketAdminSignatureMiddleware.Handle(),
		oContainer.WebsocketAdminDecryptionMiddleware.Handle(),
		oContainer.WebsocketAdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.WebsocketAdminResponseMiddleware.Handle(),
		oContainer.WebsocketAdminEncryptionMiddleware.Handle(),
	}
}

func WebsocketInit(oContainer *container.WebsocketContainer) *pkg.WebsocketRouter {
	oRouter := pkg.NewWebsocketRouter("/ws")

	oAdminGroup := oRouter.Group("Admin.")
	oAdminGroup.Use(websocketAdminMiddlewares(oContainer))
	oAdminGroup.HandleFunc("Authentication.Authenticator.SignIn", oContainer.WebsocketAdminAuthenticationAuthenticator.SignIn)

	return oRouter
}
