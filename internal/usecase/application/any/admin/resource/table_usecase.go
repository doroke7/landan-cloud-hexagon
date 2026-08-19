package resource

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkg "example/pkg"
)

type TableUsecase struct {
	*AbstractUsecase
	TableModel outputPortAnyModel.TableModel
	TableLogic outputPortAnyLogic.TableLogic
}

func NewTableUsecase(oTableModel outputPortAnyModel.TableModel, oTableLogic outputPortAnyLogic.TableLogic, oAbstractUsecase *AbstractUsecase) usecasePortAnyAdminResource.TableUsecase {
	return &TableUsecase{
		AbstractUsecase: oAbstractUsecase,
		TableModel:      oTableModel,
		TableLogic:      oTableLogic,
	}
}

func (oSelf *TableUsecase) AddOne(oValue *domain.TableValue) (bool, error) {

	_, oErr := oSelf.TableModel.AddOne(oValue)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) ShowOne(iId uint) (*domain.Table, error) {
	oTable, oErr := oSelf.TableModel.ShowOneById(iId)
	return oTable, oErr
}

func (oSelf *TableUsecase) EditOne(oValue *domain.TableValue, iId uint) (bool, error) {

	_, oErr := oSelf.TableModel.EditOneById(oValue, iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) RemoveOne(iId uint) (bool, error) {

	_, oErr := oSelf.TableModel.RemoveOneById(iId)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *TableUsecase) ShowOnes(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, uint64, error) {

	aWheres := pkg.FiltersToWheres([]string{}, aFilters)
	aOrders := pkg.SortersToOrders([]string{}, aSorters)

	oLimit := pkg.PaginationToLimit(oPagination)

	aTables, iTotal, oErr := oSelf.TableLogic.ShowTablesTotalByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	return aTables, uint64(iTotal), oErr
}
