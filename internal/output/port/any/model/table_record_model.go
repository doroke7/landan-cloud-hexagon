package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableRecordModel interface {
	AddOne(oTableRecord *domain.TableRecordValue) error
	ShowOneById(iId uint) (*domain.TableRecord, error)
	EditOneById(oTableRecord *domain.TableRecordValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecord, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
