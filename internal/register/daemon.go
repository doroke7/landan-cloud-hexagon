package register

import (
	"context"
	pkg "example/pkg"

	container "example/container"
	pbSourceAnnouncement "example/pb/source/announcement"
)

func DaemonInit(oContainer *container.DaemonContainer) *pkg.ClientRouter {
	oRouter := pkg.NewClientRouter()

	oRouter.Handle(func(ctx context.Context) error {
		oStream, err := oContainer.SourceClient.Announcement.Lottery.Watch(ctx, &pbSourceAnnouncement.LotteryWatchRequest{Key: "default"})
		return oContainer.DaemonWatcherSourceAnnouncementLottery.Watch(oStream, err)
	})

	return oRouter
}
