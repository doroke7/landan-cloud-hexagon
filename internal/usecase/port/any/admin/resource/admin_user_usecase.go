package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase interface {
	AddOne(oAdminUser *domain.AdminUserValue) error
	EditOne(oAdminUser *domain.AdminUserValue, iId uint) error
	RemoveOne(iId uint) error
	ShowOne(iId uint) (*domain.AdminUser, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint, error)
}
