package input_application_resource_model

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
)

type GameHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameModelServer
	usecasePortAnyModel.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameUsecase usecasePortAnyModel.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler: oAbstractHandler,
		GameUsecase:     oGameUsecase,
	}
}

func gameToPb(oGame *domain.Game) *pbResourceModel.Game {
	return &pbResourceModel.Game{
		Id:          uint32(oGame.Id),
		GameTypeId:  uint32(oGame.GameTypeId),
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
		CreatedAt:   timestamppb.New(oGame.CreatedAt),
		UpdatedAt:   timestamppb.New(oGame.UpdatedAt),
		DeletedAt:   timestamppb.New(oGame.DeletedAt),
		GameType: &pbResourceModel.GameType{
			Id:        uint32(oGame.GameType.Id),
			Key:       oGame.GameType.Key,
			Name:      oGame.GameType.Name,
			CreatedAt: timestamppb.New(oGame.GameType.CreatedAt),
			UpdatedAt: timestamppb.New(oGame.GameType.UpdatedAt),
			DeletedAt: timestamppb.New(oGame.GameType.DeletedAt),
		},
	}
}

func (oSelf *GameHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameAddOneInput) (*pbResourceModel.GameAddOneOutput, error) {
	fmt.Println(runtime.Caller(0))

	var oGameValue domain.GameValue

	if oReq.GameTypeId != nil {
		iGameTypeId := uint(*oReq.GameTypeId)
		oGameValue.GameTypeId = &iGameTypeId
	}

	oGameValue.Key = oReq.Key
	oGameValue.Name = oReq.Name
	oGameValue.Description = oReq.Description

	_, oErr := oSelf.GameUsecase.AddOne(&oGameValue)

	if oErr != nil {
		return nil, oErr
	}

	var iGameTypeId uint32

	if oGameValue.GameTypeId != nil {
		iGameTypeId = uint32(*oGameValue.GameTypeId)
	}

	var sKey, sName, sDescription string

	if oGameValue.Key != nil {
		sKey = *oGameValue.Key
	}
	if oGameValue.Name != nil {
		sName = *oGameValue.Name
	}
	if oGameValue.Description != nil {
		sDescription = *oGameValue.Description
	}

	return &pbResourceModel.GameAddOneOutput{
		GameTypeId:  iGameTypeId,
		Key:         sKey,
		Name:        sName,
		Description: sDescription,
	}, nil
}

func (oSelf *GameHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.GameShowOneByIdInput) (*pbResourceModel.GameShowOneByIdOutput, error) {

	oGame, oErr := oSelf.GameUsecase.ShowOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	if oGame == nil {
		return nil, nil
	}

	oPbGame := gameToPb(oGame)

	return &pbResourceModel.GameShowOneByIdOutput{
		Id:          oPbGame.Id,
		GameTypeId:  oPbGame.GameTypeId,
		Key:         oPbGame.Key,
		Name:        oPbGame.Name,
		Description: oPbGame.Description,
		CreatedAt:   oPbGame.CreatedAt,
		UpdatedAt:   oPbGame.UpdatedAt,
		DeletedAt:   oPbGame.DeletedAt,
		GameType:    oPbGame.GameType,
	}, nil

}

func (oSelf *GameHandler) ShowOnesByWheresWithOrdersLimit(oContext context.Context, oReq *pbResourceModel.GameShowOnesByWheresWithOrdersLimitInput) (*pbResourceModel.GameShowOnesByWheresWithOrdersLimitOutput, error) {

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

	iOffset := uint(oReq.GetLimit().GetOffset())
	iCount := uint(oReq.GetLimit().GetCount())
	oLimit := &pkg.Limit{
		Offset: &iOffset,
		Count:  &iCount,
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

	aGames, oErr := oSelf.GameUsecase.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	if oErr != nil {
		return nil, oErr
	}

	aPbGames := make([]*pbResourceModel.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, gameToPb(oGame))
	}

	return &pbResourceModel.GameShowOnesByWheresWithOrdersLimitOutput{
		Games: aPbGames,
	}, nil

}

func (oSelf *GameHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameEditOneByIdInput) (*pbResourceModel.GameEditOneByIdOutput, error) {

	var oGameValue domain.GameValue

	if oReq.GameTypeId != nil {
		iGameTypeId := uint(*oReq.GameTypeId)
		oGameValue.GameTypeId = &iGameTypeId
	}

	oGameValue.Key = oReq.Key
	oGameValue.Name = oReq.Name
	oGameValue.Description = oReq.Description

	_, oErr := oSelf.GameUsecase.EditOneById(&oGameValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameRemoveOneByIdInput) (*pbResourceModel.GameRemoveOneByIdOutput, error) {

	_, oErr := oSelf.GameUsecase.RemoveOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameHandler) TotalByWheres(oContext context.Context, oReq *pbResourceModel.GameTotalByWheresInput) (*pbResourceModel.GameTotalByWheresOutput, error) {

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

	iTotal, oErr := oSelf.GameUsecase.TotalByWheres(aWheres)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTotalByWheresOutput{
		Total: iTotal,
	}, nil
}
