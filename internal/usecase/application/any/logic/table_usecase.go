package usecaseApplicationAnyLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkg "example/pkg"
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

func (oSelf *TableUsecase) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, int64, error) {

	aTables, iTotal, oErr := oSelf.TableLogic.ShowTablesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aTables, iTotal, oErr
}
