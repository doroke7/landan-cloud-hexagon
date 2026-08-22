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
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, error)
}
