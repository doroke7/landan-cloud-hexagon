package port

import (
	domain "example/internal/domain"
)

type LotteryModel interface {
	EditOneByKey(oLottery *domain.LotteryValue, sKey string) (bool, error)
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
