package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameModel interface {
	AddOne(oGame *domain.GameValue) (bool, error)
	ShowOneById(iId uint) (*domain.Game, error)
	EditOneById(oGame *domain.GameValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)

	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)

	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
