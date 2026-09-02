package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeModel interface {
	AddOne(oGameType *domain.GameTypeValue) (bool, error)
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oGameType *domain.GameTypeValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByParentId(iParentId uint) ([]*domain.GameType, error)
	TotalByParentId(iParentId uint) (uint64, error)

	ShowOnes() ([]*domain.GameType, error)

	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
