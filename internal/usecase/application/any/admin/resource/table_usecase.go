package usecaseApplicationAnyAdminResource

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	outputPortAnyModel "example/internal/output/port/any/model"
	usecaseApplicationAnyAdmin "example/internal/usecase/application/any/admin"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pkgInput "example/pkg/input"
)

type TableUsecase struct {
	*usecaseApplicationAnyAdmin.AbstractUsecase
	TableModel outputPortAnyModel.TableModel
	TableLogic outputPortAnyLogic.TableLogic
}

func NewTableUsecase(oTableModel outputPortAnyModel.TableModel, oTableLogic outputPortAnyLogic.TableLogic, oAbstractUsecase *usecaseApplicationAnyAdmin.AbstractUsecase) usecasePortAnyAdminResource.TableUsecase {
	return &TableUsecase{
		AbstractUsecase: oAbstractUsecase,
		TableModel:      oTableModel,
		TableLogic:      oTableLogic,
	}
}

func (oSelf *TableUsecase) AddOne(oValue *domain.TableValue) error {

	return oSelf.TableModel.AddOne(oValue)
}

func (oSelf *TableUsecase) ShowOne(iId uint64) (*domain.Table, error) {
	oTable, oErr := oSelf.TableModel.ShowOneById(uint(iId))
	return oTable, oErr
}

func (oSelf *TableUsecase) EditOne(oValue *domain.TableValue, iId uint64) error {

	return oSelf.TableModel.EditOneById(oValue, iId)
}

func (oSelf *TableUsecase) RemoveOne(iId uint64) error {

	return oSelf.TableModel.RemoveOneById(iId)
}

func (oSelf *TableUsecase) ShowOnes(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error) {

	aTables, iTotal, oErr := oSelf.TableLogic.ShowTablesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aTables, uint64(iTotal), oErr
}
