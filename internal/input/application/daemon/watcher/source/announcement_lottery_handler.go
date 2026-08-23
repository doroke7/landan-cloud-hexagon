package applicationDaemonWatcherSource

import (
	"io"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"go.uber.org/zap"

	domain "example/internal/domain"
	inputApplicationDaemon "example/internal/input/application/daemon"
	usecasePortAnyWatcherSource "example/internal/usecase/port/any/watcher/source"
	pbSourceAnnouncement "example/pb/source/announcement"
	pkg "example/pkg"
)

type AnnouncementLotteryHandler struct {
	*inputApplicationDaemon.AbstractHandler
	lotteryUsecase usecasePortAnyWatcherSource.AnnouncementLotteryUsecase
}

func NewAnnouncementLotteryHandler(oLotteryUsecase usecasePortAnyWatcherSource.AnnouncementLotteryUsecase, oAbstractHandler *inputApplicationDaemon.AbstractHandler) *AnnouncementLotteryHandler {
	return &AnnouncementLotteryHandler{
		AbstractHandler: oAbstractHandler,
		lotteryUsecase:  oLotteryUsecase,
	}
}

// Watch 只負責讀 stream、轉呼叫 usecase；stream 要連誰、怎麼開，交給 register 決定，
// 這裡完全不知道 gRPC client 怎麼建立的——但開 stream 當下的錯誤，一樣由這裡統一判斷。
func (oSelf *AnnouncementLotteryHandler) Watch(oStream grpc.ServerStreamingClient[pbSourceAnnouncement.LotteryWatchResponse], err error) error {
	if err != nil {
		pkg.Logger(pkg.DeamonWatcher).Error("開啟 lottery stream 失敗", zap.Error(err))
		return err
	}

	for {
		oReply, err := oStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			pkg.Logger(pkg.Client).Error("讀取 lottery stream 失敗", zap.Error(err))
			return err
		}

		aNumbers := make([]string, 0, len(oReply.Numbers))
		for _, iNumber := range oReply.Numbers {
			aNumbers = append(aNumbers, strconv.Itoa(int(iNumber)))
		}

		sRound := oReply.Round
		iTime := oReply.Time
		sNumbers := strings.Join(aNumbers, ",")

		oLotteryValue := &domain.LotteryValue{
			Round:   &sRound,
			Time:    &iTime,
			Numbers: &sNumbers,
		}

		if err := oSelf.lotteryUsecase.Watch(oLotteryValue); err != nil {
			pkg.Logger(pkg.DeamonWatcher).Error("處理開獎資料失敗", zap.Error(err))
			continue
		}
	}
}
