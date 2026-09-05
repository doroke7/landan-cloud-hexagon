package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameValue) error
	ShowOne(iId uint) (*domain.Game, error)
	EditOne(oValue *domain.GameValue, iId uint) error
	RemoveOne(iId uint) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint, error)
}
