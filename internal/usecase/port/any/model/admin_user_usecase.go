package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminUserUsecase interface {
	ShowOneByName(sName string) (*domain.AdminUser, error)
	ShowOneById(iId uint64) (*domain.AdminUser, error)
	AddOne(oAdminUser *domain.AdminUserVariable) error
	EditOneById(oAdminUser *domain.AdminUserVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
