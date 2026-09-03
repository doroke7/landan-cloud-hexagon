package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionModel interface {
	AddOne(oAdminPermission *domain.AdminPermissionValue) (bool, error)
	ShowOneById(iId uint) (*domain.AdminPermission, error)
	EditOneById(oAdminPermission *domain.AdminPermissionValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermission, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
