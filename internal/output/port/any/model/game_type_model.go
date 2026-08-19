package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameTypeModel interface {
	AddOne(oGameType *domain.GameTypeValue) (bool, error)
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oGameType *domain.GameTypeValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)

	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.GameType, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
}
