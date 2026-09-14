package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase interface {
	AddOne(oAdminRole *domain.AdminRoleVariable) error
	ShowOne(iId uint64) (*domain.AdminRole, error)
	EditOne(oAdminRole *domain.AdminRoleVariable, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint64, error)
}
