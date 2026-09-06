package outputApplicationResourceModel

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
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

func domainTableValueToProtoTableVariable(oTable *domain.TableValue) *pbResourceModel.TableVariable {
	oVariable := &pbResourceModel.TableVariable{
		No:          oTable.No,
		Key:         oTable.Key,
		Description: oTable.Description,
	}

	if oTable.GameId != nil {
		iGameId := uint64(*oTable.GameId)
		oVariable.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint64(*oTable.State)
		oVariable.State = &iState
	}

	if oTable.Result != nil {
		sResult := *oTable.Result
		oVariable.Result = &sResult
	}

	if oTable.StartedAt != nil {
		oVariable.StartedAt = timestamppb.New(*oTable.StartedAt)
	}

	if oTable.EndedAt != nil {
		oVariable.EndedAt = timestamppb.New(*oTable.EndedAt)
	}

	return oVariable
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint64) error {

	oRequest := &pbResourceModel.TableEditOneByIdInput{
		Id:       uint64(iId),
		Variable: domainTableValueToProtoTableVariable(oTable),
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

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) error {

	oRequest := &pbResourceModel.TableAddOneInput{
		Variable: domainTableValueToProtoTableVariable(oTable),
	}

	_, oErr := oSelf.ResourceModelClient.Table.AddOne(oSelf.Context, oRequest)

	return oErr
}
