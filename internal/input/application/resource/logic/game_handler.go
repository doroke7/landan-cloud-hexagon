package input_application_resource_logic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type GameHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedGameLogicServer
	usecasePortAnyLogic.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameUsecase usecasePortAnyLogic.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler: oAbstractHandler,
		GameUsecase:     oGameUsecase,
	}
}

func gameToPb(oGame *domain.Game) *pbResourceLogic.Game {
	return &pbResourceLogic.Game{
		Id:          uint32(oGame.Id),
		GameTypeId:  uint32(oGame.GameTypeId),
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
		CreatedAt:   timestamppb.New(oGame.CreatedAt),
		UpdatedAt:   timestamppb.New(oGame.UpdatedAt),
		DeletedAt:   timestamppb.New(oGame.DeletedAt),
		GameType: &pbResourceLogic.GameType{
			Id:        uint32(oGame.GameType.Id),
			Key:       oGame.GameType.Key,
			Name:      oGame.GameType.Name,
			CreatedAt: timestamppb.New(oGame.GameType.CreatedAt),
			UpdatedAt: timestamppb.New(oGame.GameType.UpdatedAt),
			DeletedAt: timestamppb.New(oGame.GameType.DeletedAt),
		},
	}
}

func (oSelf *GameHandler) ShowGamesTotalByWheresWithOrdersLimit(oContext context.Context, oReq *pbResourceLogic.GameShowGamesTotalByWheresWithOrdersLimitInput) (*pbResourceLogic.GameShowGamesTotalByWheresWithOrdersLimitOutput, error) {

	iOffset := uint(oReq.GetLimit().GetOffset())
	iCount := uint(oReq.GetLimit().GetCount())
	oLimit := &pkg.Limit{
		Offset: &iOffset,
		Count:  &iCount,
	}

	aWheres := make([]*pkg.Where, 0, len(oReq.GetWheres()))
	for _, oOne := range oReq.GetWheres() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aWheres = append(aWheres, &pkg.Where{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aOrders := make([]*pkg.Order, 0, len(oReq.GetOrders()))
	for _, oOne := range oReq.GetOrders() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sValue := oOne.GetValue()
		aOrders = append(aOrders, &pkg.Order{
			Field: &sField,
			Value: &sValue,
		})
	}

	aGames, iTotal, oErr := oSelf.GameUsecase.ShowGamesTotalByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	aPbGames := make([]*pbResourceLogic.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, gameToPb(oGame))
	}

	return &pbResourceLogic.GameShowGamesTotalByWheresWithOrdersLimitOutput{
		Total: uint64(iTotal),
		Games: aPbGames,
	}, oErr
}
