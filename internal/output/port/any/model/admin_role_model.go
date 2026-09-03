package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleModel interface {
	AddOne(oAdminRole *domain.AdminRoleValue) (bool, error)
	ShowOneById(iId uint) (*domain.AdminRole, error)
	EditOneById(oAdminRole *domain.AdminRoleValue, iId uint) (bool, error)
	RemoveOneById(iId uint) (bool, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
