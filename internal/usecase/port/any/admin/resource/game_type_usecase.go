package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeUsecase interface {
	AddOne(oValue *domain.GameTypeValue) error
	ShowOne(iId uint) (*domain.GameType, error)
	ShowTree() (*domain.GameType, error)
	EditOne(oValue *domain.GameTypeValue, iId uint) error
	RemoveOne(iId uint) error
	ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint, error)
}
