package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type AdminUserUsecase interface {
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, uint64, error)
}
