package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameValue) (bool, error)
	ShowOneById(iId uint) (*domain.Game, error)
	EditOneById(oValue *domain.GameValue, iId uint64) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, error)
}
