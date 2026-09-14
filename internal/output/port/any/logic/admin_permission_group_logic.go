package outputPortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic interface {
	AddAdminPermissionGroup(oAdminPermissionGroupValue *domain.AdminPermissionGroupVariable) error
	EditAdminPermissionGroupById(oAdminPermissionGroupValue *domain.AdminPermissionGroupVariable, iId uint64) error
	ShowTree() ([]*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error)
}
