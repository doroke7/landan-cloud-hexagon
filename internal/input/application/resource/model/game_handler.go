package inputApplicationResourceModel

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameModelServer
	ModelGameUsecase usecasePortAnyModel.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameUsecase usecasePortAnyModel.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler:  oAbstractHandler,
		ModelGameUsecase: oGameUsecase,
	}
}

// protoGameValueToDomainGameValue 把 gRPC 帶進來的 variable 攤成 domain.GameValue（指標欄位可選）。
func protoGameValueToDomainGameValue(oVariable *pbResourceModel.GameValue) domain.GameValue {
	var oValue domain.GameValue
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name
	oValue.Description = oVariable.Description

	if oVariable.GameTypeId != nil {
		iGameTypeId := uint64(*oVariable.GameTypeId)
		oValue.GameTypeId = &iGameTypeId
	}

	return oValue
}

func domainGameToProtoGame(oGame *domain.Game) *pbResource.Game {
	return &pbResource.Game{
		Id:          uint64(oGame.Id),
		GameTypeId:  uint64(oGame.GameTypeId),
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
		CreatedAt:   timestamppb.New(oGame.CreatedAt),
		UpdatedAt:   timestamppb.New(oGame.UpdatedAt),
		DeletedAt:   timestamppb.New(oGame.DeletedAt),
		GameType:    domainGameTypeToProtoGameType(&oGame.GameType),
	}
}

func (oSelf *GameHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameAddOneInput) (*pbResourceModel.GameAddOneOutput, error) {
	fmt.Println(runtime.Caller(0))

	oGameValue := protoGameValueToDomainGameValue(oReq.GetValue())

	oErr := oSelf.ModelGameUsecase.AddOne(&oGameValue)

	if oErr != nil {
		return nil, oErr
	}

	oProtoGame := &pbResource.Game{}

	if oGameValue.GameTypeId != nil {
		oProtoGame.GameTypeId = uint64(*oGameValue.GameTypeId)
	}
	if oGameValue.Key != nil {
		oProtoGame.Key = *oGameValue.Key
	}
	if oGameValue.Name != nil {
		oProtoGame.Name = *oGameValue.Name
	}
	if oGameValue.Description != nil {
		oProtoGame.Description = *oGameValue.Description
	}

	return &pbResourceModel.GameAddOneOutput{
		Game: oProtoGame,
	}, nil
}

func (oSelf *GameHandler) ShowOneByKey(oContext context.Context, oReq *pbResourceModel.GameShowOneByKeyInput) (*pbResourceModel.GameShowOneByKeyOutput, error) {

	oGame, oErr := oSelf.ModelGameUsecase.ShowOneByKey(oReq.Key)

	if oErr != nil {
		return nil, oErr
	}

	if oGame == nil {
		return nil, nil
	}

	oProtoGame := domainGameToProtoGame(oGame)

	return &pbResourceModel.GameShowOneByKeyOutput{
		Game: oProtoGame,
	}, nil

}

func (oSelf *GameHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameEditOneByIdInput) (*pbResourceModel.GameEditOneByIdOutput, error) {

	oGameValue := protoGameValueToDomainGameValue(oReq.GetValue())

	oErr := oSelf.ModelGameUsecase.EditOneById(&oGameValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameRemoveOneByIdInput) (*pbResourceModel.GameRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelGameUsecase.RemoveOneById(uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.GameTotalByFiltersInput) (*pbResourceModel.GameTotalByFiltersOutput, error) {

	aFilters := make([]*pkgInput.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aFilters = append(aFilters, &pkgInput.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	iTotal, oErr := oSelf.ModelGameUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
