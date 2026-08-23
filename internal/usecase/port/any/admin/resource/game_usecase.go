package usecasePortAnyAdminResource

import (
	domain "example/internal/domain"
	pkg "example/pkg"
)

type GameUsecase interface {
	AddOne(oValue *domain.GameValue) (bool, error)
	ShowOne(iId uint) (*domain.Game, error)
	EditOne(oValue *domain.GameValue, iId uint) (bool, error)
	RemoveOne(iId uint) (bool, error)
	ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, uint64, error)
}
