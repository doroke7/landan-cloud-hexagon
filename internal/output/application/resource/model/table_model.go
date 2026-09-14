package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
)

type TableModel struct {
	*AbstractModel
}

func NewTableModel(oAbstractModel *AbstractModel) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableVariable, iId uint64) error {

	oRequest := &pbResourceModel.TableEditOneByIdInput{
		Id:    uint64(iId),
		Value: pkgDomainToProto.TableValue(oTable),
	}

	_, oErr := oSelf.ResourceModelClient.Table.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *TableModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.Table.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.TableRemoveOneByIdInput{Id: iId},
	)

	return oErr
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.TableTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *TableModel) AddOne(oTable *domain.TableVariable) error {

	oRequest := &pbResourceModel.TableAddOneInput{
		Value: pkgDomainToProto.TableValue(oTable),
	}

	_, oErr := oSelf.ResourceModelClient.Table.AddOne(oSelf.Context, oRequest)

	return oErr
}
