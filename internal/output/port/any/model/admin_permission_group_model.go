package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupModel interface {
	AddOne(oAdminPermissionGroupValue *domain.AdminPermissionGroupValue) error
	EditOneById(oAdminPermissionGroupValue *domain.AdminPermissionGroupValue, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
