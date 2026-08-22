package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameUsecase interface {
	ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, int64, error)
}
