package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableRecordModel interface {
	AddOne(oTableRecord *domain.TableRecordVariable) error
	ShowOneById(iId uint64) (*domain.TableRecord, error)
	EditOneById(oTableRecord *domain.TableRecordVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecord, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
