package inputApplicationResourceModel

import (
	"context"

	pbResource "example/pb"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameTypeModelServer
	ModelGameTypeUsecase usecasePortAnyModel.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyModel.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler:      oAbstractHandler,
		ModelGameTypeUsecase: oGameTypeUsecase,
	}
}

func (oSelf *GameTypeHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameTypeAddOneInput) (*pbResourceModel.GameTypeAddOneOutput, error) {

	oValue := oReq.GetVariable()

	var oGameTypeValue domain.GameTypeVariable
	if oValue != nil {
		oGameTypeValue.Key = oValue.Key
		oGameTypeValue.Name = oValue.Name
	}

	oErr := oSelf.ModelGameTypeUsecase.AddOne(&oGameTypeValue)

	if oErr != nil {
		return nil, oErr
	}

	var sKey, sName string

	if oGameTypeValue.Key != nil {
		sKey = *oGameTypeValue.Key
	}
	if oGameTypeValue.Name != nil {
		sName = *oGameTypeValue.Name
	}

	return &pbResourceModel.GameTypeAddOneOutput{
		GameType: &pbResource.GameType{
			Key:  sKey,
			Name: sName,
		},
	}, nil
}

func (oSelf *GameTypeHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameTypeEditOneByIdInput) (*pbResourceModel.GameTypeEditOneByIdOutput, error) {

	oValue := oReq.GetVariable()

	var oGameTypeValue domain.GameTypeVariable
	if oValue != nil {
		oGameTypeValue.Key = oValue.Key
		oGameTypeValue.Name = oValue.Name
	}

	oErr := oSelf.ModelGameTypeUsecase.EditOneById(&oGameTypeValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeEditOneByIdOutput{}, nil
}

func (oSelf *GameTypeHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameTypeRemoveOneByIdInput) (*pbResourceModel.GameTypeRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelGameTypeUsecase.RemoveOneById(uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeRemoveOneByIdOutput{}, nil
}

func (oSelf *GameTypeHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.GameTypeTotalByFiltersInput) (*pbResourceModel.GameTypeTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelGameTypeUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
