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
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.TableRecord, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
}
