package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

// TableRecordLog 是牌局過程中的事件紀錄（audit trail），只增不改，
// 所以這裡沒有 EditOneById／RemoveOneById，跟 TableRecordModel 不同。
type TableRecordLogModel interface {
	AddOne(oTableRecordLog *domain.TableRecordLogValue) error
	ShowOneById(iId uint64) (*domain.TableRecordLog, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecordLog, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
