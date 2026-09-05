package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase interface {
	AddOne(oAdminPermission *domain.AdminPermissionValue) error
	ShowOne(iId uint64) (*domain.AdminPermission, error)
	EditOne(oAdminPermission *domain.AdminPermissionValue, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint64, error)
}
