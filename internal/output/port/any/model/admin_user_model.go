package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserModel interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserValue) error
	EditOneById(oAdminUser *domain.AdminUserValue, iId uint) error
	RemoveOneById(iId uint) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint, error)
}
