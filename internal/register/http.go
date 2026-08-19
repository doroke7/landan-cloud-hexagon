package register

import (
	"github.com/gin-gonic/gin"

	container "example/container"
)

func httpAdminAuthenticationMiddlewares(oContainer *container.HttpContainer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// ALL middleware
		oContainer.HttpAdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.HttpAdminErrorMiddleware.Handle(),
		oContainer.HttpAdminSignatureMiddleware.Handle(),
		oContainer.HttpAdminDecryptionMiddleware.Handle(),
		oContainer.HttpAdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.HttpAdminResponseMiddleware.Handle(),
		oContainer.HttpAdminEncryptionMiddleware.Handle(),
	}
}

func httpAdminResourceMiddlewares(oContainer *container.HttpContainer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// ALL middleware
		oContainer.HttpAdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.HttpAdminErrorMiddleware.Handle(),
		oContainer.HttpAdminSignatureMiddleware.Handle(),
		oContainer.HttpAdminDecryptionMiddleware.Handle(),
		oContainer.HttpAdminAuthenticationMiddleware.Handle(),
		oContainer.HttpAdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.HttpAdminResponseMiddleware.Handle(),
		oContainer.HttpAdminEncryptionMiddleware.Handle(),
	}
}

func httpAdminOptionMiddlewares(oContainer *container.HttpContainer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// ALL middleware
		oContainer.HttpAdminLoggerMiddleware.Handle(),

		// Before Middleware
		oContainer.HttpAdminErrorMiddleware.Handle(),
		oContainer.HttpAdminSignatureMiddleware.Handle(),
		oContainer.HttpAdminDecryptionMiddleware.Handle(),
		oContainer.HttpAdminAuthenticationMiddleware.Handle(),
		oContainer.HttpAdminRequestMiddleware.Handle(),

		// After Middleware
		oContainer.HttpAdminResponseMiddleware.Handle(),
		oContainer.HttpAdminEncryptionMiddleware.Handle(),
	}
}

func HttpInit(oGin *gin.Engine, oContainer *container.HttpContainer) *gin.Engine {

	oAdmin := oGin.Group("/Admin")
	{

		oAdminAuthentication := oAdmin.Group("/Authentication")
		oAdminAuthentication.Use(httpAdminAuthenticationMiddlewares(oContainer)...)
		{
			oAdminAuthentication.POST("/Authenticator/SignIn", oContainer.HttpAdminAuthenticationAuthenticator.SignIn)
			oAdminAuthentication.POST("/Authenticator/Refresh", oContainer.HttpAdminAuthenticationAuthenticator.Refresh)
		}

		oAdminResource := oAdmin.Group("/Resource")
		oAdminResource.Use(httpAdminResourceMiddlewares(oContainer)...)
		{
			oAdminResource.POST("/Game/AddOne", oContainer.HttpAdminResourceGame.AddOne)
			oAdminResource.PUT("/Game/EditOne", oContainer.HttpAdminResourceGame.EditOne)
			oAdminResource.DELETE("/Game/RemoveOne", oContainer.HttpAdminResourceGame.RemoveOne)
			oAdminResource.GET("/Game/ShowOne", oContainer.HttpAdminResourceGame.ShowOne)
			oAdminResource.GET("/Game/ShowOnes", oContainer.HttpAdminResourceGame.ShowOnes)

			oAdminResource.POST("/Table/AddOne", oContainer.HttpAdminResourceTable.AddOne)
			oAdminResource.GET("/Table/ShowOne", oContainer.HttpAdminResourceTable.ShowOne)
			oAdminResource.PUT("/Table/EditOne", oContainer.HttpAdminResourceTable.EditOne)
			oAdminResource.DELETE("/Table/RemoveOne", oContainer.HttpAdminResourceTable.RemoveOne)
			oAdminResource.GET("/Table/ShowOnes", oContainer.HttpAdminResourceTable.ShowOnes)

			oAdminResource.POST("/AdminUser/AddOne", oContainer.HttpAdminResourceAdminUser.AddOne)
			oAdminResource.GET("/AdminUser/ShowOnes", oContainer.HttpAdminResourceAdminUser.ShowOnes)

			oAdminResource.POST("/GameType/AddOne", oContainer.HttpAdminResourceGameType.AddOne)
			oAdminResource.PUT("/GameType/EditOne", oContainer.HttpAdminResourceGameType.EditOne)
			oAdminResource.DELETE("/GameType/RemoveOne", oContainer.HttpAdminResourceGameType.RemoveOne)
			oAdminResource.GET("/GameType/ShowOne", oContainer.HttpAdminResourceGameType.ShowOne)
			oAdminResource.GET("/GameType/ShowOnes", oContainer.HttpAdminResourceGameType.ShowOnes)
		}

		oAdminOption := oAdmin.Group("/Option")
		oAdminOption.Use(httpAdminOptionMiddlewares(oContainer)...)
		{
			oAdminOption.GET("/GameType/Select", oContainer.HttpAdminOptionGameType.Select)
		}

	}

	return oGin
}
