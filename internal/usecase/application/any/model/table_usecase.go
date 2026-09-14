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

func (oSelf *TableUsecase) AddOne(oValue *domain.TableVariable) error {

	oErr := oSelf.TableModel.AddOne(oValue)

	return oErr
}

func (oSelf *TableUsecase) EditOneById(oValue *domain.TableVariable, iId uint64) error {

	oErr := oSelf.TableModel.EditOneById(oValue, iId)

	return oErr
}

func (oSelf *TableUsecase) RemoveOneById(iId uint64) error {

	oErr := oSelf.TableModel.RemoveOneById(iId)

	return oErr
}

func (oSelf *TableUsecase) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	iTotal, oErr := oSelf.TableModel.TotalByFilters(aFilters)
	return uint64(iTotal), oErr
}
