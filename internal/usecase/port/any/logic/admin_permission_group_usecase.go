package usecasePortAnyLogic

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase interface {
	AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupVariable) error
	EditAdminPermissionGroupById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error
	ShowTree() ([]*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error)
	ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error)
	ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error)
}
