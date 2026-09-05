package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionModel interface {
	AddOne(oAdminPermission *domain.AdminPermissionValue) error
	ShowOneById(iId uint) (*domain.AdminPermission, error)
	EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
