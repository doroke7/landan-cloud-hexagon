package outputPortAnyModel

import (
	domain "example/internal/domain"
	pkgInput "example/pkg/input"
)

type TableModel interface {
	AddOne(oTable *domain.TableVariable) error
	EditOneById(oTable *domain.TableVariable, iId uint64) error
	RemoveOneById(iId uint64) error
	TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error)
}
