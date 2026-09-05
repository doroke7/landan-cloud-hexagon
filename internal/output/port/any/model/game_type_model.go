package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeModel interface {
	AddOne(oGameType *domain.GameTypeValue) error
	ShowOneById(iId uint) (*domain.GameType, error)
	EditOneById(oGameType *domain.GameTypeValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByParentId(iParentId uint) ([]*domain.GameType, error)
	TotalByParentId(iParentId uint) (uint, error)

	ShowOnes() ([]*domain.GameType, error)

	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
