package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase interface {
	AddOne(oAdminUser *domain.AdminUserValue) error
	EditOne(oAdminUser *domain.AdminUserValue, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOne(iId uint64) (*domain.AdminUser, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error)
}
