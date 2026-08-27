package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableUsecase interface {
	AddOne(oValue *domain.TableValue) (bool, error)
	ShowOneById(iId uint) (*domain.Table, error)
	EditOneById(oValue *domain.TableValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error)
}
