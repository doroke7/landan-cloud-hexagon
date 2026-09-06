package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) error
	EditOneById(oValue *domain.GameTypeValue, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error)
}
