package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type GameTypeModel interface {
	AddOne(oGameType *domain.GameTypeValue) error
	EditOneById(oGameType *domain.GameTypeValue, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByParentId(iParentId uint64) (uint64, error)
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
