package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkgInput "example/pkg/input"
)

type TableUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.TableLogic
}

func NewTableUsecase(oAbstractUsecase *AbstractUsecase, oTableLogic outputPortAnyLogic.TableLogic) usecasePortAnyLogic.TableUsecase {
	return &TableUsecase{
		AbstractUsecase: oAbstractUsecase,
		TableLogic:      oTableLogic,
	}
}

func (oSelf *TableUsecase) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error) {

	aTables, iTotal, oErr := oSelf.TableLogic.ShowTablesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aTables, uint64(iTotal), oErr
}
