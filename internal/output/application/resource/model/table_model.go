package resource

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
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

func (oSelf *TableModel) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, error) {

	oRequest := &pbResourceModel.TableShowOnesByWheresWithOrdersLimitInput{
		Limit: &pbResourceModel.TableLimit{
			Offset: uint64(*oLimit.Offset),
			Count:  uint64(*oLimit.Count),
		},
	}

	for _, oWhere := range aWheres {
		if oWhere == nil || oWhere.Field == nil || oWhere.Operator == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oWhere.Value)
		if oErr != nil {
			continue
		}

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceModel.TableWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil || oOrder.Value == nil {
			continue
		}

		oRequest.Orders = append(oRequest.Orders, &pbResourceModel.TableOrder{
			Field: *oOrder.Field,
			Value: *oOrder.Value,
		})
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.ShowOnesByWheresWithOrdersLimit(oSelf.Context, oRequest)

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

func (oSelf *TableModel) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {

	oRequest := &pbResourceModel.TableTotalByWheresInput{}

	for _, oWhere := range aWheres {
		if oWhere == nil || oWhere.Field == nil || oWhere.Operator == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oWhere.Value)
		if oErr != nil {
			continue
		}

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceModel.TableWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	oResponse, oErr := oSelf.ResourceModelClient.Table.TotalByWheres(oSelf.Context, oRequest)

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
