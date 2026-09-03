package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase interface {
	AddOne(oAdminUser *domain.AdminUserValue) (bool, error)
	EditOne(oAdminUser *domain.AdminUserValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error)
}
