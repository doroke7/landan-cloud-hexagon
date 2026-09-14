package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic interface {
	AddAdminPermissionGroup(oAdminPermissionGroupVariable *domain.AdminPermissionGroupVariable) error
	EditAdminPermissionGroupById(oAdminPermissionGroupVariable *domain.AdminPermissionGroupVariable, iId uint64) error
	ShowTree() ([]*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error)
}
