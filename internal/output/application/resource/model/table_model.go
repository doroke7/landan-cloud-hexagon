package outputApplicationResourceModel

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type TableModel struct {
	*resourceBase.AbstractResource
}

func NewTableModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractResource: oAbstractModel,
	}
}

func protoTableToDomainTable(oTable *pbResource.Table) domain.Table {
	if oTable == nil {
		return domain.Table{}
	}

	return domain.Table{
		Id:          uint(oTable.GetId()),
		No:          oTable.GetNo(),
		GameId:      uint(oTable.GetGameId()),
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
		iGameId := uint32(*oTable.GameId)
		oVariable.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint32(*oTable.State)
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
		&pbResourceModel.TableShowOneByIdInput{Id: uint32(iId)},
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

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.TableEditOneByIdInput{
		Id:       uint32(iId),
		Variable: domainTableValueToProtoTableVariable(oTable),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.EditOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *TableModel) RemoveOneById(iId uint) (bool, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Table.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.TableRemoveOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {

	oRequest := &pbResourceModel.TableShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aTables := make([]*domain.Table, 0, len(oResponse.GetTables()))
	for _, oOne := range oResponse.GetTables() {
		oTable := protoTableToDomainTable(oOne)
		aTables = append(aTables, &oTable)
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.TableTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.TotalByFilters(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
}

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) (bool, error) {

	oRequest := &pbResourceModel.TableAddOneInput{
		Variable: domainTableValueToProtoTableVariable(oTable),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.AddOne(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}
