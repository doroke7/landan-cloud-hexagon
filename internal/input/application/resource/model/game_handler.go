package inputApplicationResourceModel

import (
	"context"
	"fmt"
	"runtime"

	pbResource "example/pb"
	pbResourceModel "example/pb/resource/model"

	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
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

func (oSelf *GameHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameAddOneInput) (*pbResourceModel.GameAddOneOutput, error) {
	fmt.Println(runtime.Caller(0))

	oGameVariable := pkgProtoToDomain.GameVariable(oReq.GetVariable())

	oErr := oSelf.ModelGameUsecase.AddOne(&oGameVariable)

	if oErr != nil {
		return nil, oErr
	}

	oProtoGame := &pbResource.Game{}

	if oGameVariable.GameTypeId != nil {
		oProtoGame.GameTypeId = uint64(*oGameVariable.GameTypeId)
	}
	if oGameVariable.Key != nil {
		oProtoGame.Key = *oGameVariable.Key
	}
	if oGameVariable.Name != nil {
		oProtoGame.Name = *oGameVariable.Name
	}
	if oGameVariable.Description != nil {
		oProtoGame.Description = *oGameVariable.Description
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

	oProtoGame := pkgDomainToProto.Game(oGame)

	return &pbResourceModel.GameShowOneByKeyOutput{
		Game: oProtoGame,
	}, nil

}

func (oSelf *GameHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameEditOneByIdInput) (*pbResourceModel.GameEditOneByIdOutput, error) {

	oGameValue := pkgProtoToDomain.GameVariable(oReq.GetVariable())

	oErr := oSelf.ModelGameUsecase.EditOneById(&oGameValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameEditOneByIdOutput{}, nil
}

func (oSelf *GameHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameRemoveOneByIdInput) (*pbResourceModel.GameRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelGameUsecase.RemoveOneById(uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameRemoveOneByIdOutput{}, nil
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
