package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase interface {
	AddOne(oAdminRole *domain.AdminRoleValue) error
	ShowOne(iId uint) (*domain.AdminRole, error)
	EditOne(oAdminRole *domain.AdminRoleValue, iId uint) error
	RemoveOne(iId uint) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint, error)
}
