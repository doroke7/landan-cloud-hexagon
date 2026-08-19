package usecase

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
)

type TableUsecase struct {
	outputPortAnyModel.TableModel
}

func NewTableUsecase(oTableModel outputPortAnyModel.TableModel) usecasePortAnyModel.TableUsecase {
	return &TableUsecase{
		TableModel: oTableModel,
	}
}

func (oSelf *TableUsecase) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, error) {

	return oSelf.TableModel.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)
}

func (oSelf *TableUsecase) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {

	return oSelf.TableModel.TotalByWheres(aWheres)
}
