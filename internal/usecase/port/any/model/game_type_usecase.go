package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) (bool, error)
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oValue *domain.GameTypeValue, iId uint64) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error)
}
