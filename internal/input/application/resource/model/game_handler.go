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
	usecasePortAnyModel.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameUsecase usecasePortAnyModel.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler: oAbstractHandler,
		GameUsecase:     oGameUsecase,
	}
}

// protoGameVariableToDomainGameValue 把 gRPC 帶進來的 variable 攤成 domain.GameValue（指標欄位可選）。
func protoGameVariableToDomainGameValue(oVariable *pbResourceModel.GameVariable) domain.GameValue {
	var oValue domain.GameValue
	if oVariable == nil {
		return oValue
	}

	oValue.Key = oVariable.Key
	oValue.Name = oVariable.Name
	oValue.Description = oVariable.Description

	if oVariable.GameTypeId != nil {
		iGameTypeId := uint(*oVariable.GameTypeId)
		oValue.GameTypeId = &iGameTypeId
	}

	return oValue
}

func domainGameToProtoGame(oGame *domain.Game) *pbResource.Game {
	return &pbResource.Game{
		Id:          uint32(oGame.Id),
		GameTypeId:  uint32(oGame.GameTypeId),
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

	oGameValue := protoGameVariableToDomainGameValue(oReq.GetVariable())

	_, oErr := oSelf.GameUsecase.AddOne(&oGameValue)

	if oErr != nil {
		return nil, oErr
	}

	oProtoGame := &pbResource.Game{}

	if oGameValue.GameTypeId != nil {
		oProtoGame.GameTypeId = uint32(*oGameValue.GameTypeId)
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

func (oSelf *GameHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.GameShowOneByIdInput) (*pbResourceModel.GameShowOneByIdOutput, error) {

	oGame, oErr := oSelf.GameUsecase.ShowOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	if oGame == nil {
		return nil, nil
	}

	oProtoGame := domainGameToProtoGame(oGame)

	return &pbResourceModel.GameShowOneByIdOutput{
		Game: oProtoGame,
	}, nil

}

func (oSelf *GameHandler) ShowOneByKey(oContext context.Context, oReq *pbResourceModel.GameShowOneByKeyInput) (*pbResourceModel.GameShowOneByKeyOutput, error) {

	oGame, oErr := oSelf.GameUsecase.ShowOneByKey(oReq.Key)

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

func (oSelf *GameHandler) ShowOnesByFiltersWithOrdersPagination(oContext context.Context, oReq *pbResourceModel.GameShowOnesByFiltersWithOrdersPaginationInput) (*pbResourceModel.GameShowOnesByFiltersWithOrdersPaginationOutput, error) {

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

	aSorters := make([]*pkgInput.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		aSorters = append(aSorters, &pkgInput.Sorter{
			Field: &sField,
			Order: &sOrder,
		})
	}

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aGames, oErr := oSelf.GameUsecase.ShowOnesByFiltersWithOrdersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, oErr
	}

	aPbGames := make([]*pbResource.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, domainGameToProtoGame(oGame))
	}

	return &pbResourceModel.GameShowOnesByFiltersWithOrdersPaginationOutput{
		Games: aPbGames,
	}, nil

}

func (oSelf *GameHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameEditOneByIdInput) (*pbResourceModel.GameEditOneByIdOutput, error) {

	oGameValue := protoGameVariableToDomainGameValue(oReq.GetVariable())

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

	iTotal, oErr := oSelf.GameUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
