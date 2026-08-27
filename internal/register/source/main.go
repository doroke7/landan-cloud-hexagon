package registerSource

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pbSourceAnnouncement "example/pb/source/announcement"

	container "example/container"
	pkgGrpc "example/pkg/grpc"
)

func sourceInterceptors(_ *container.SourceContainer) grpc.UnaryServerInterceptor {

	// 全局攔截器：目前 resource 只需要 facade -> resource 的 Basic Auth 驗證。
	// Group("", ...) 註冊在 root 節點上，任何 method 都會先經過 root。
	aBase := []grpc.UnaryServerInterceptor{
		// oContainer.ResourceModelAuthenticationInterceptor.Handle(),
	}

	oRouter := pkgGrpc.NewGrpcRouter()

	oRouter.Group("", aBase...)
	// database 群組目前只需要全局驗證；未來個別 resource 服務需要額外攔截器
	// （例如特定表要多一層權限檢查）時，在這裡加一個更具體的 Group 前綴即可

	return oRouter.Build()
}

func Init(oContainer *container.SourceContainer) *grpc.Server {

	oGrpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(sourceInterceptors(oContainer)),
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
	pbSourceAnnouncement.RegisterLotteryControllerServer(oGrpcServer, oContainer.SourceAnnouncementLottery)

	return oGrpcServer
}
