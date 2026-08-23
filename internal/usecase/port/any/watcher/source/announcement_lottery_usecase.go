package usecasePortAnyWatcherSource

import (
	domain "example/internal/domain"
)

type AnnouncementLotteryUsecase interface {
	Watch(oLottery *domain.LotteryValue) error
}
