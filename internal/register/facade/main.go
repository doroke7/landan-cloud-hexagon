package registerFacade

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pkgGrpc "example/pkg/grpc"

	container "example/container"

	pbFacadeAdminAuthentication "example/pb/facade/admin/authentication"
	pbFacadeAdminResource "example/pb/facade/admin/resource"
	pbFacadeRegister "example/pb/facade/register"
	pbFacadeTable "example/pb/facade/table"
)

func facadeInterceptors(oContainer *container.FacadeContainer) grpc.UnaryServerInterceptor {

	aGameInterceptors := []grpc.UnaryServerInterceptor{
		oContainer.FacadeGameErrorInterceptor.Handle(),
		oContainer.FacadeGameStatusInterceptor.Handle(),
		oContainer.FacadeGameLoggerInterceptor.Handle(),
	}

	aAdminInterceptors := []grpc.UnaryServerInterceptor{
		oContainer.FacadeAdminErrorInterceptor.Handle(),
		oContainer.FacadeAdminStatusInterceptor.Handle(),
		oContainer.FacadeAdminLoggerInterceptor.Handle(),
		oContainer.FacadeAdminSignatureInterceptor.Handle(),
		oContainer.FacadeAdminDecryptionInterceptor.Handle(),
		// Encryption 自己管理 holder，不依賴其他攔截器的順序，放在請求階段的
		// Signature/Decryption 之後只是邏輯上比較好懂（回應階段的事放最後）。
		oContainer.FacadeAdminEncryptionInterceptor.Handle(),
	}

	oRouter := pkgGrpc.NewGrpcRouter()

	oRouter.Group("/pb.facade.game.", aGameInterceptors...)
	oRouter.Group("/pb.facade.game.authentication.")
	oRouter.Group("/pb.facade.game.resource.", oContainer.FacadeGameAuthenticationInterceptor.Handle())

	oRouter.Group("/pb.facade.admin.", aAdminInterceptors...)

	return oRouter.Build()
}

func Init(oContainer *container.FacadeContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(facadeInterceptors(oContainer)),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    1 * time.Second,
				Timeout: 5 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	pbFacadeTable.RegisterScannerControllerServer(oGrpcServer, oContainer.FacadeTableScanner)
	pbFacadeRegister.RegisterAuthenticatorControllerServer(oGrpcServer, oContainer.FacadeTableAuthenticator)
	pbFacadeAdminAuthentication.RegisterAuthenticatorControllerServer(oGrpcServer, oContainer.FacadeAdminAuthenticationAuthenticator)
	pbFacadeAdminResource.RegisterGameControllerServer(oGrpcServer, oContainer.FacadeAdminResourceGame)

	return oGrpcServer
}
