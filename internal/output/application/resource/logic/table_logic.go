package resource

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type TableLogic struct {
	*AbstractLogic
}

func NewTableLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *TableLogic) ShowTablesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, int64, error) {

	oRequest := &pbResourceLogic.TableShowTablesTotalByWheresWithOrdersLimitInput{
		Limit: &pbResourceLogic.TableLimit{
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

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceLogic.TableWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil || oOrder.Value == nil {
			continue
		}

		oRequest.Orders = append(oRequest.Orders, &pbResourceLogic.TableOrder{
			Field: *oOrder.Field,
			Value: *oOrder.Value,
		})
	}

	oResponse, oErr := oSelf.ResourceLogicClient.Table.ShowTablesTotalByWheresWithOrdersLimit(oSelf.Context, oRequest)

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

	iTotal := oResponse.GetTotal()
	return aTables, int64(iTotal), oErr
}
