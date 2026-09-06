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

func domainTableValueToProtoTableValue(oTable *domain.TableValue) *pbResourceModel.TableValue {
	oValue := &pbResourceModel.TableValue{
		No:          oTable.No,
		Key:         oTable.Key,
		Description: oTable.Description,
	}

	if oTable.GameId != nil {
		iGameId := uint64(*oTable.GameId)
		oValue.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint64(*oTable.State)
		oValue.State = &iState
	}

	if oTable.Result != nil {
		sResult := *oTable.Result
		oValue.Result = &sResult
	}

	if oTable.StartedAt != nil {
		oValue.StartedAt = timestamppb.New(*oTable.StartedAt)
	}

	if oTable.EndedAt != nil {
		oValue.EndedAt = timestamppb.New(*oTable.EndedAt)
	}

	return oValue
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint64) error {

	oRequest := &pbResourceModel.TableEditOneByIdInput{
		Id:    uint64(iId),
		Value: domainTableValueToProtoTableValue(oTable),
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
		Value: domainTableValueToProtoTableValue(oTable),
	}

	_, oErr := oSelf.ResourceModelClient.Table.AddOne(oSelf.Context, oRequest)

	return oErr
}
