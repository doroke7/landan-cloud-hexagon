package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type GameHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedGameLogicServer
	LogicGameUsecase usecasePortAnyLogic.GameUsecase
}

func NewGameHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameUsecase usecasePortAnyLogic.GameUsecase) *GameHandler {
	return &GameHandler{
		AbstractHandler:  oAbstractHandler,
		LogicGameUsecase: oGameUsecase,
	}
}

// domainGameTypeToProtoGameType 遞迴帶出 Parent / Children，讓巢狀的 game type 一起過 gRPC。
func domainGameTypeToProtoGameType(oGameType *domain.GameType) *pbResource.GameType {
	if oGameType == nil {
		return nil
	}

	oPb := &pbResource.GameType{
		Id:        uint32(oGameType.Id),
		ParentId:  uint32(oGameType.ParentId),
		Key:       oGameType.Key,
		Name:      oGameType.Name,
		CreatedAt: timestamppb.New(oGameType.CreatedAt),
		UpdatedAt: timestamppb.New(oGameType.UpdatedAt),
		DeletedAt: timestamppb.New(oGameType.DeletedAt),
		Parent:    domainGameTypeToProtoGameType(oGameType.Parent),
	}

	for i := range oGameType.Children {
		oPb.Children = append(oPb.Children, domainGameTypeToProtoGameType(&oGameType.Children[i]))
	}

	return oPb
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

func (oSelf *GameHandler) ShowGamesTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationOutput, error) {

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

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

	aGames, iTotal, oErr := oSelf.LogicGameUsecase.ShowGamesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	aPbGames := make([]*pbResource.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aPbGames = append(aPbGames, domainGameToProtoGame(oGame))
	}

	return &pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationOutput{
		Total: iTotal,
		Games: aPbGames,
	}, oErr
}
