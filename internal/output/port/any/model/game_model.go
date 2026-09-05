package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameModel interface {
	AddOne(oGame *domain.GameValue) error
	ShowOneById(iId uint) (*domain.Game, error)
	ShowOneByKey(sKey string) (*domain.Game, error)
	EditOneById(oGame *domain.GameValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByGameTypeId(iGameTypeId uint) ([]*domain.Game, error)
	TotalByGameTypeId(iGameTypeId uint) (uint, error)

	ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
