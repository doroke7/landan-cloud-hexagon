package inputApplicationResourceModel

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/protobuf/types/known/timestamppb"

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

func domainGameToProtoGame(oGame *domain.Game) *pbResourceModel.Game {
	return &pbResourceModel.Game{
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

	oPbGame := domainGameToProtoGame(oGame)

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

	aPbGames := make([]*pbResourceModel.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, domainGameToProtoGame(oGame))
	}

	return &pbResourceModel.GameShowOnesByFiltersWithOrdersPaginationOutput{
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
