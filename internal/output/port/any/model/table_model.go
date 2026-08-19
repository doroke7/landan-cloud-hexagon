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
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
}
