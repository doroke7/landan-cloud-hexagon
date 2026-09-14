package usecasePortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupUsecase interface {
	AddOne(oValue *domain.AdminPermissionGroupVariable) error
	EditOneById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
	ShowOnesByParentId(iParentId uint64) ([]*domain.AdminPermissionGroup, error)
	ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error)
}
