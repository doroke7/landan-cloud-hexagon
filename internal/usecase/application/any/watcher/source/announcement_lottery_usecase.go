package usecaseApplicationAnyWatcherSource

import (
	"go.uber.org/zap"

	pkgUtility "example/pkg/utility"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyWatcher "example/internal/usecase/application/any/watcher"
	usecasePortAnyWatcherSource "example/internal/usecase/port/any/watcher/source"
)

type AnnouncementLotteryUsecase struct {
	*usecaseApplicationAnyWatcher.AbstractUsecase
	outputPortAnyModel.LotteryModel
}

func NewAnnouncementLotteryUsecase(oAbstractUsecase *usecaseApplicationAnyWatcher.AbstractUsecase, oLotteryRepository outputPortAnyModel.LotteryModel) usecasePortAnyWatcherSource.AnnouncementLotteryUsecase {
	return &AnnouncementLotteryUsecase{
		AbstractUsecase: oAbstractUsecase,
		LotteryModel:    oLotteryRepository,
	}
}

// Watch 收到一筆開獎資料，用 Round 當 key 落地存起來。
func (oSelf *AnnouncementLotteryUsecase) Watch(oValue *domain.LotteryValue) error {

	var oLottery domain.Lottery

	if oValue.Round != nil {
		oLottery.Round = *oValue.Round
	}

	if oValue.Time != nil {
		oLottery.Time = *oValue.Time
	}

	if oValue.Numbers != nil {
		oLottery.Numbers = *oValue.Numbers
	}

	if err := oSelf.LotteryModel.EditOneByKey(oValue, oLottery.Round); err != nil {
		pkgUtility.Logger(pkgUtility.Client).Error("儲存開獎資料失敗",
			zap.Uint("id", oLottery.Id),
			zap.String("round", oLottery.Round),
			zap.Error(err),
		)
		return err
	}

	pkgUtility.Logger(pkgUtility.Client).Info("收到開獎資料",
		zap.Uint("id", oLottery.Id),
		zap.String("round", oLottery.Round),
		zap.String("numbers", oLottery.Numbers),
	)

	return nil
}
