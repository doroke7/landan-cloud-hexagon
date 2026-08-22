package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableRecordModel interface {
	AddOne(oTableRecord *domain.TableRecordValue) (bool, error)
	ShowOneById(iId uint) (*domain.TableRecord, error)
	EditOneById(oTableRecord *domain.TableRecordValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecord, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
