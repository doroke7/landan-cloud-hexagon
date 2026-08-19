package input_application_resource_model

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameTypeModelServer
	usecasePortAnyModel.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyModel.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler: oAbstractHandler,
		GameTypeUsecase: oGameTypeUsecase,
	}
}

func gameTypeToPb(oGameType *domain.GameType) *pbResourceModel.GameType {
	return &pbResourceModel.GameType{
		Id:        uint32(oGameType.Id),
		Key:       oGameType.Key,
		Name:      oGameType.Name,
		CreatedAt: timestamppb.New(oGameType.CreatedAt),
		UpdatedAt: timestamppb.New(oGameType.UpdatedAt),
		DeletedAt: timestamppb.New(oGameType.DeletedAt),
	}
}

func (oSelf *GameTypeHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameTypeAddOneInput) (*pbResourceModel.GameTypeAddOneOutput, error) {

	var oGameTypeValue domain.GameTypeValue

	oGameTypeValue.Key = oReq.Key
	oGameTypeValue.Name = oReq.Name

	_, oErr := oSelf.GameTypeUsecase.AddOne(&oGameTypeValue)

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
		Key:  sKey,
		Name: sName,
	}, nil
}

func (oSelf *GameTypeHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.GameTypeShowOneByIdInput) (*pbResourceModel.GameTypeShowOneByIdOutput, error) {

	oGameType, oErr := oSelf.GameTypeUsecase.ShowOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	if oGameType == nil {
		return nil, nil
	}

	oPbGameType := gameTypeToPb(oGameType)

	return &pbResourceModel.GameTypeShowOneByIdOutput{
		Id:        oPbGameType.Id,
		Key:       oPbGameType.Key,
		Name:      oPbGameType.Name,
		CreatedAt: oPbGameType.CreatedAt,
		UpdatedAt: oPbGameType.UpdatedAt,
		DeletedAt: oPbGameType.DeletedAt,
	}, nil
}

func (oSelf *GameTypeHandler) ShowOnesByWheresWithOrdersLimit(oContext context.Context, oReq *pbResourceModel.GameTypeShowOnesByWheresWithOrdersLimitInput) (*pbResourceModel.GameTypeShowOnesByWheresWithOrdersLimitOutput, error) {

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

	aGameTypes, oErr := oSelf.GameTypeUsecase.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	if oErr != nil {
		return nil, oErr
	}

	aPbGameTypes := make([]*pbResourceModel.GameType, 0, len(aGameTypes))
	for _, oGameType := range aGameTypes {
		aPbGameTypes = append(aPbGameTypes, gameTypeToPb(oGameType))
	}

	return &pbResourceModel.GameTypeShowOnesByWheresWithOrdersLimitOutput{
		GameTypes: aPbGameTypes,
	}, nil
}

func (oSelf *GameTypeHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameTypeEditOneByIdInput) (*pbResourceModel.GameTypeEditOneByIdOutput, error) {

	var oGameTypeValue domain.GameTypeValue

	oGameTypeValue.Key = oReq.Key
	oGameTypeValue.Name = oReq.Name

	_, oErr := oSelf.GameTypeUsecase.EditOneById(&oGameTypeValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameTypeHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameTypeRemoveOneByIdInput) (*pbResourceModel.GameTypeRemoveOneByIdOutput, error) {

	_, oErr := oSelf.GameTypeUsecase.RemoveOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameTypeHandler) TotalByWheres(oContext context.Context, oReq *pbResourceModel.GameTypeTotalByWheresInput) (*pbResourceModel.GameTypeTotalByWheresOutput, error) {

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

	iTotal, oErr := oSelf.GameTypeUsecase.TotalByWheres(aWheres)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeTotalByWheresOutput{
		Total: iTotal,
	}, nil
}
