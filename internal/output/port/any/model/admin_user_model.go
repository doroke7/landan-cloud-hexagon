package port

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type AdminUserModel interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.AdminUser, error)
	TotalByWheres(aWheres []*pkg.Where) (uint64, error)
}
