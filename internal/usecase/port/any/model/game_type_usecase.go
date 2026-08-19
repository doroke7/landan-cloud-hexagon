package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) (bool, error)
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oValue *domain.GameTypeValue, iId uint64) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.GameType, error)
}
