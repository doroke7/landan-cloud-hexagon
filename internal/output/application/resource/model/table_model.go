package resource

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkg "example/pkg"
)

type TableModel struct {
	*AbstractModel
}

func NewTableModel(oAbstractModel *AbstractModel) outputPortAnyModel.TableModel {
	return &TableModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Table.ShowOneById(
		oSelf.Context,
		&pbResourceModel.TableShowOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	if oResponse.GetId() == 0 {
		return nil, nil
	}

	return &domain.Table{
		Id:          uint(oResponse.GetId()),
		No:          oResponse.GetNo(),
		GameId:      uint(oResponse.GetGameId()),
		Key:         oResponse.GetKey(),
		State:       uint8(oResponse.GetState()),
		Description: oResponse.GetDescription(),
		Result:      json.RawMessage(oResponse.GetResult()),
		StartedAt:   oResponse.GetStartedAt().AsTime(),
		EndedAt:     oResponse.GetEndedAt().AsTime(),
		CreatedAt:   oResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   oResponse.GetUpdatedAt().AsTime(),
		DeletedAt:   oResponse.GetDeletedAt().AsTime(),
		Game:        gameFromPb(oResponse.GetGame()),
	}, nil
}

func (oSelf *TableModel) EditOneById(oTable *domain.TableValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.TableEditOneByIdInput{Id: uint32(iId)}

	oRequest.No = oTable.No
	oRequest.Key = oTable.Key
	oRequest.Description = oTable.Description

	if oTable.GameId != nil {
		iGameId := uint32(*oTable.GameId)
		oRequest.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint32(*oTable.State)
		oRequest.State = &iState
	}

	if oTable.Result != nil {
		sResult := string(*oTable.Result)
		oRequest.Result = &sResult
	}

	if oTable.StartedAt != nil {
		oRequest.StartedAt = timestamppb.New(*oTable.StartedAt)
	}

	if oTable.EndedAt != nil {
		oRequest.EndedAt = timestamppb.New(*oTable.EndedAt)
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

func (oSelf *TableModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, error) {

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
		aTables = append(aTables, &domain.Table{
			Id:          uint(oOne.GetId()),
			No:          oOne.GetNo(),
			GameId:      uint(oOne.GetGameId()),
			Key:         oOne.GetKey(),
			State:       uint8(oOne.GetState()),
			Description: oOne.GetDescription(),
			Result:      json.RawMessage(oOne.GetResult()),
			StartedAt:   oOne.GetStartedAt().AsTime(),
			EndedAt:     oOne.GetEndedAt().AsTime(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			Game:        gameFromPb(oOne.GetGame()),
		})
	}

	return aTables, nil
}

func (oSelf *TableModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {

	oRequest := &pbResourceModel.TableTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.TotalByFilters(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
}

func (oSelf *TableModel) AddOne(oTable *domain.TableValue) (bool, error) {

	oRequest := &pbResourceModel.TableAddOneInput{}

	oRequest.No = oTable.No
	oRequest.Key = oTable.Key
	oRequest.Description = oTable.Description

	if oTable.GameId != nil {
		iGameId := uint32(*oTable.GameId)
		oRequest.GameId = &iGameId
	}

	if oTable.State != nil {
		iState := uint32(*oTable.State)
		oRequest.State = &iState
	}

	if oTable.Result != nil {
		sResult := string(*oTable.Result)
		oRequest.Result = &sResult
	}

	if oTable.StartedAt != nil {
		oRequest.StartedAt = timestamppb.New(*oTable.StartedAt)
	}

	if oTable.EndedAt != nil {
		oRequest.EndedAt = timestamppb.New(*oTable.EndedAt)
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.AddOne(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}
