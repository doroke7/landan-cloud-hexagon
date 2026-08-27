package usecaseApplicationAnyModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type TableUsecase struct {
	outputPortAnyModel.TableModel
}

func NewTableUsecase(oTableModel outputPortAnyModel.TableModel) usecasePortAnyModel.TableUsecase {
	return &TableUsecase{
		TableModel: oTableModel,
	}
}

func (oSelf *TableUsecase) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {

	return oSelf.TableModel.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
}

func (oSelf *TableUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	return oSelf.TableModel.TotalByFilters(aFilters)
}
