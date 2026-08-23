package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableLogic interface {
	ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, int64, error)
}
