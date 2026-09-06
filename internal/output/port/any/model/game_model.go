package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameModel interface {
	AddOne(oGame *domain.GameValue) error
	ShowOneById(iId uint64) (*domain.Game, error)
	ShowOneByKey(sKey string) (*domain.Game, error)
	EditOneById(oGame *domain.GameValue, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByGameTypeId(iGameTypeId uint64) (uint, error)

	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
