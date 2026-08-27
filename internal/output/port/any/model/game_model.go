package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameModel interface {
	AddOne(oGame *domain.GameValue) (bool, error)
	ShowOneById(iId uint) (*domain.Game, error)
	EditOneById(oGame *domain.GameValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)

	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
