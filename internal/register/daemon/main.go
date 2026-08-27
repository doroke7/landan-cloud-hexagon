package registerDaemon

import (
	"context"
	pkgGrpc "example/pkg/grpc"

	container "example/container"
	pbSourceAnnouncement "example/pb/source/announcement"
)

func Init(oContainer *container.DaemonContainer) *pkgGrpc.ClientRouter {
	oRouter := pkgGrpc.NewClientRouter()

	oRouter.Handle(func(ctx context.Context) error {
		oStream, err := oContainer.SourceClient.Announcement.Lottery.Watch(ctx, &pbSourceAnnouncement.LotteryWatchRequest{Key: "default"})
		return oContainer.DaemonWatcherSourceAnnouncementLottery.Watch(oStream, err)
	})

	return oRouter
}
