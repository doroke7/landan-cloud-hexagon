package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameValue) (bool, error)
	ShowOne(iId uint) (*domain.Game, error)
	EditOne(oValue *domain.GameValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error)
}
