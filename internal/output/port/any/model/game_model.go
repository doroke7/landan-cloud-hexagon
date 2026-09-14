package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameModel interface {
	AddOne(oGame *domain.GameVariable) error
	ShowOneByKey(sKey string) (*domain.Game, error)
	EditOneById(oGame *domain.GameVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByGameTypeId(iGameTypeId uint64) (uint, error)

	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
