package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeVariable) error
	ShowOne(iId uint64) (*domain.GameType, error)
	ShowTree() (*domain.GameType, error)
	EditOne(oValue *domain.GameTypeVariable, iId uint64) error
	RemoveOne(iId uint64) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error)
}
