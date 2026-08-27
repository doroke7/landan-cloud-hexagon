package usecasePortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, int64, error)
}
