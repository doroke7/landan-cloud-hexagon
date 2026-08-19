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
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, error)
}
