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

	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
