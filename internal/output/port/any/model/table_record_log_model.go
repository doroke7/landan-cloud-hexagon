package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

// TableRecordLog 是牌局過程中的事件紀錄（audit trail），只增不改，
// 所以這裡沒有 EditOneById／RemoveOneById，跟 TableRecordModel 不同。
type TableRecordLogModel interface {
	AddOne(oTableRecordLog *domain.TableRecordLogValue) (bool, error)
	ShowOneById(iId uint) (*domain.TableRecordLog, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.TableRecordLog, error)
	TotalByFilters(aFilters []*pkg.Filter) (uint64, error)
}
