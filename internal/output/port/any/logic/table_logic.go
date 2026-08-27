package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableLogic interface {
	ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, int64, error)
}
