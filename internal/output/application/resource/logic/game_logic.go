package resource

import (
	"google.golang.org/protobuf/types/known/structpb"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type GameLogic struct {
	*AbstractLogic
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func gameTypeFromPb(oPbGameType *pbResourceLogic.GameType) domain.GameType {
	if oPbGameType == nil {
		return domain.GameType{}
	}

	return domain.GameType{
		Id:        uint(oPbGameType.GetId()),
		Key:       oPbGameType.GetKey(),
		Name:      oPbGameType.GetName(),
		CreatedAt: oPbGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oPbGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oPbGameType.GetDeletedAt().AsTime(),
	}
}

func gameFromPb(oPbGame *pbResourceLogic.Game) domain.Game {
	if oPbGame == nil {
		return domain.Game{}
	}

	return domain.Game{
		Id:          uint(oPbGame.GetId()),
		GameTypeId:  uint(oPbGame.GetGameTypeId()),
		Key:         oPbGame.GetKey(),
		Name:        oPbGame.GetName(),
		Description: oPbGame.GetDescription(),
		CreatedAt:   oPbGame.GetCreatedAt().AsTime(),
		UpdatedAt:   oPbGame.GetUpdatedAt().AsTime(),
		DeletedAt:   oPbGame.GetDeletedAt().AsTime(),
		GameType:    gameTypeFromPb(oPbGame.GetGameType()),
	}
}

func (oSelf *GameLogic) ShowGamesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Game, int64, error) {

	oRequest := &pbResourceLogic.GameShowGamesTotalByWheresWithOrdersLimitInput{
		Limit: &pbResourceLogic.GameLimit{
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

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceLogic.GameWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil || oOrder.Value == nil {
			continue
		}

		oRequest.Orders = append(oRequest.Orders, &pbResourceLogic.GameOrder{
			Field: *oOrder.Field,
			Value: *oOrder.Value,
		})
	}

	oResponse, oErr := oSelf.ResourceLogicClient.Game.ShowGamesTotalByWheresWithOrdersLimit(oSelf.Context, oRequest)

	aGames := make([]*domain.Game, 0, len(oResponse.GetGames()))
	for sKey, oOne := range oResponse.GetGames() {
		aGames = append(aGames, &domain.Game{
			Id:          uint(oOne.GetId()),
			GameTypeId:  uint(oOne.GetGameTypeId()),
			Key:         oOne.GetKey(),
			Name:        oOne.GetName(),
			Description: oOne.GetDescription(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			GameType:    gameTypeFromPb(oOne.GetGameType()),
		})
		_ = sKey
	}

	iTotal := oResponse.GetTotal()
	return aGames, int64(iTotal), oErr
}
