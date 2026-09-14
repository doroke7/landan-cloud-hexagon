package outputPortAnyModel

import (
	domain "example/internal/domain"
)

type LotteryModel interface {
	EditOneByKey(oLottery *domain.LotteryVariable, sKey string) error
	WatchOneByKey(sKey string) (*domain.Lottery, error)
}
