package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameVariable) error
	ShowOne(iId uint64) (*domain.Game, error)
	EditOne(oValue *domain.GameVariable, iId uint64) error
	RemoveOne(iId uint) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error)
}
