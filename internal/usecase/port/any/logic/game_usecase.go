package any

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameUsecase interface {
	ShowGamesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, int64, error)
}
