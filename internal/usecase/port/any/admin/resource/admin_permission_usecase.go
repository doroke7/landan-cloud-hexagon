package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionUsecase interface {
	AddOne(oAdminPermission *domain.AdminPermissionVariable) error
	ShowOne(iId uint64) (*domain.AdminPermission, error)
	EditOne(oAdminPermission *domain.AdminPermissionVariable, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, uint64, error)
}
