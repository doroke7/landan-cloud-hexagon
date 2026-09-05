package usecasePortAnyAdminOption

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type AdminRoleUsecase interface {
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminRole, uint, error)
}
