package registerResource

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pbResourceEvent "example/pb/resource/event"
	pbResourceLogic "example/pb/resource/logic"
	pbResourceModel "example/pb/resource/model"

	container "example/container"
	pkgGrpc "example/pkg/grpc"
)

func resourceInterceptors(oContainer *container.ResourceContainer) grpc.UnaryServerInterceptor {

	oRouter := pkgGrpc.NewGrpcRouter()

	oRouter.Group("pb.resource.model",
		oContainer.ResourceModelErrorInterceptor.Handle(),
		oContainer.ResourceModelLoggerInterceptor.Handle(),
		oContainer.ResourceModelAuthenticationInterceptor.Handle(),
	)
	oRouter.Group("pb.resource.logic",
		oContainer.ResourceLogicErrorInterceptor.Handle(),
		oContainer.ResourceLogicLoggerInterceptor.Handle(),
		oContainer.ResourceLogicAuthenticationInterceptor.Handle(),
	)

	oRouter.Group("pb.resource.event",
		oContainer.ResourceModelErrorInterceptor.Handle(),
		oContainer.ResourceModelLoggerInterceptor.Handle(),
		oContainer.ResourceModelAuthenticationInterceptor.Handle(),
	)

	return oRouter.Build()
}

func Init(oContainer *container.ResourceContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(resourceInterceptors(oContainer)),
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

	pbResourceModel.RegisterAdminUserModelServer(oGrpcServer, oContainer.ResourceModelAdminUser)
	pbResourceModel.RegisterAdminRoleModelServer(oGrpcServer, oContainer.ResourceModelAdminRole)
	pbResourceModel.RegisterAdminPermissionModelServer(oGrpcServer, oContainer.ResourceModelAdminPermission)
	pbResourceModel.RegisterAppUserModelServer(oGrpcServer, oContainer.ResourceModelAppUser)
	pbResourceModel.RegisterGameModelServer(oGrpcServer, oContainer.ResourceModelGame)
	pbResourceModel.RegisterTableModelServer(oGrpcServer, oContainer.ResourceModelTable)
	pbResourceModel.RegisterGameTypeModelServer(oGrpcServer, oContainer.ResourceModelGameType)
	pbResourceModel.RegisterAdminPermissionGroupModelServer(oGrpcServer, oContainer.ResourceModelAdminPermissionGroup)
	pbResourceLogic.RegisterAdminUserLogicServer(oGrpcServer, oContainer.ResourceLogicAdminUser)
	pbResourceLogic.RegisterAdminRoleLogicServer(oGrpcServer, oContainer.ResourceLogicAdminRole)
	pbResourceLogic.RegisterGameLogicServer(oGrpcServer, oContainer.ResourceLogicGame)
	pbResourceLogic.RegisterTableLogicServer(oGrpcServer, oContainer.ResourceLogicTable)
	pbResourceLogic.RegisterGameTypeLogicServer(oGrpcServer, oContainer.ResourceLogicGameType)
	pbResourceLogic.RegisterAdminPermissionGroupLogicServer(oGrpcServer, oContainer.ResourceLogicAdminPermissionGroup)
	pbResourceEvent.RegisterAdminUserEventServer(oGrpcServer, oContainer.ResourceEventAdminUser)

	return oGrpcServer
}
