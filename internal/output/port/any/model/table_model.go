package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableModel interface {
	AddOne(oTable *domain.TableValue) error
	ShowOneById(iId uint) (*domain.Table, error)
	EditOneById(oTable *domain.TableValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
