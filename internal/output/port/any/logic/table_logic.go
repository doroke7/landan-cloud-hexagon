package logic

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type TableLogic interface {
	ShowTablesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, int64, error)
}
