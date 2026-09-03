package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase interface {
	AddOne(oAdminPermission *domain.AdminPermissionValue) (bool, error)
	ShowOne(iId uint) (*domain.AdminPermission, error)
	EditOne(oAdminPermission *domain.AdminPermissionValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint64, error)
}
