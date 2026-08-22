package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableModel interface {
	AddOne(oTable *domain.TableValue) (bool, error)
	ShowOneById(iId uint) (*domain.Table, error)
	EditOneById(oTable *domain.TableValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
