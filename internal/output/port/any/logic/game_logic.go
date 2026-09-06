package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameLogic interface {
	ShowGameById(iId uint64) (*domain.Game, error)

	ShowGamesByGameTypeId(iGameTypeId uint64) ([]*domain.Game, error)

	ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error)
}
