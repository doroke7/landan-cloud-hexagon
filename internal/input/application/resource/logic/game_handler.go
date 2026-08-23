package applicationResourceModelLogic

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

func (oSelf *GameHandler) ShowGamesTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationOutput, error) {

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkg.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aFilters := make([]*pkg.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aFilters = append(aFilters, &pkg.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aSorters := make([]*pkg.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		aSorters = append(aSorters, &pkg.Sorter{
			Field: &sField,
			Order: &sOrder,
		})
	}

	aGames, iTotal, oErr := oSelf.GameUsecase.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	aPbGames := make([]*pbResourceLogic.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, gameToPb(oGame))
	}

	return &pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationOutput{
		Total: uint64(iTotal),
		Games: aPbGames,
	}, oErr
}
