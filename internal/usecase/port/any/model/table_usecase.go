package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) (bool, error)
	ShowOneById(iId uint) (*domain.Table, error)
	EditOneById(oValue *domain.TableValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error)
}
