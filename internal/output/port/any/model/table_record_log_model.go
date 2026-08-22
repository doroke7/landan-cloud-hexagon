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
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.TableRecordLog, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
}
