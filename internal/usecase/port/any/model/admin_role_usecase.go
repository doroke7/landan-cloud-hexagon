package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase interface {
	AddOne(oAdminRole *domain.AdminRoleVariable) error
	ShowOneById(iId uint) (*domain.AdminRole, error)
	EditOneById(oAdminRole *domain.AdminRoleVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	ShowOnes() ([]*domain.AdminRole, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
