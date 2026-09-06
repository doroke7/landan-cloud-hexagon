package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeModel interface {
	AddOne(oGameType *domain.GameTypeValue) error
	EditOneById(oGameType *domain.GameTypeValue, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnesByParentId(iParentId uint64) ([]*domain.GameType, error)
	TotalByParentId(iParentId uint64) (uint64, error)

	ShowOnes() ([]*domain.GameType, error)

	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
