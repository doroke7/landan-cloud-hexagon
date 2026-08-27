package usecasePortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, int64, error)
}
