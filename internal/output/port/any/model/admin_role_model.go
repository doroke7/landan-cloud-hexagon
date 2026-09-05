package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleModel interface {
	AddOne(oAdminRole *domain.AdminRoleValue) error
	ShowOneById(iId uint) (*domain.AdminRole, error)
	EditOneById(oAdminRole *domain.AdminRoleValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
