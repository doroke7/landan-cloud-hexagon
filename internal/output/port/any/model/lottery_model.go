package outputPortAnyModel

import (
	domain "example/internal/domain"
)

type LotteryModel interface {
	EditOneByKey(oLottery *domain.LotteryValue, sKey string) error
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
