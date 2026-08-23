package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) (bool, error)
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oValue *domain.GameTypeValue, iId uint64) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error)
}
