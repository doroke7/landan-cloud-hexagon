package usecasePortAnyAnnouncement

import (
	domain "example/internal/domain"
)

type LotteryUsecase interface {
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
