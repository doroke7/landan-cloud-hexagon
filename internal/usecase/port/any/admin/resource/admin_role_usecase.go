package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase interface {
	AddOne(oAdminRole *domain.AdminRoleValue) (bool, error)
	ShowOne(iId uint) (*domain.AdminRole, error)
	EditOne(oAdminRole *domain.AdminRoleValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint64, error)
}
