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

func (oSelf *TableUsecase) AddOne(oValue *domain.TableValue) (bool, error) {

	oErr := oSelf.TableModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) EditOneById(oValue *domain.TableValue, iId uint) (bool, error) {

	oErr := oSelf.TableModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) RemoveOneById(iId uint) (bool, error) {

	oErr := oSelf.TableModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.TableModel.TotalByFilters(aFilters)
	return uint64(iTotal), oErr
}
