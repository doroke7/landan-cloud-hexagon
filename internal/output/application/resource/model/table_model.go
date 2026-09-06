package outputApplicationResourceModel

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
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

func protoTableToDomainTable(oTable *pbResource.Table) domain.Table {
	if oTable == nil {
		return domain.Table{}
	}

	return domain.Table{
		Id:          uint64(oTable.GetId()),
		No:          oTable.GetNo(),
		GameId:      uint64(oTable.GetGameId()),
		Key:         oTable.GetKey(),
		State:       uint8(oTable.GetState()),
		Description: oTable.GetDescription(),
		Result:      oTable.GetResult(),
		StartedAt:   oTable.GetStartedAt().AsTime(),
		EndedAt:     oTable.GetEndedAt().AsTime(),
		CreatedAt:   oTable.GetCreatedAt().AsTime(),
		UpdatedAt:   oTable.GetUpdatedAt().AsTime(),
		DeletedAt:   oTable.GetDeletedAt().AsTime(),
		Game:        protoGameToDomainGame(oTable.GetGame()),
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

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Table.ShowOneById(
		oSelf.Context,
		&pbResourceModel.TableShowOneByIdInput{Id: uint64(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoTable := oResponse.GetTable()
	if oProtoTable.GetId() == 0 {
		return nil, nil
	}

	oTable := protoTableToDomainTable(oProtoTable)

	return &oTable, nil
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
