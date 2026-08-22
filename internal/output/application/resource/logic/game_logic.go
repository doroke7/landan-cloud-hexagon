package resource

import (
	"google.golang.org/protobuf/types/known/structpb"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type GameLogic struct {
	*AbstractLogic
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func gameTypeFromPb(oPbGameType *pbResourceLogic.GameType) domain.GameType {
	if oPbGameType == nil {
		return domain.GameType{}
	}

	return domain.GameType{
		Id:        uint(oPbGameType.GetId()),
		Key:       oPbGameType.GetKey(),
		Name:      oPbGameType.GetName(),
		CreatedAt: oPbGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oPbGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oPbGameType.GetDeletedAt().AsTime(),
	}
}

func gameFromPb(oPbGame *pbResourceLogic.Game) domain.Game {
	if oPbGame == nil {
		return domain.Game{}
	}

	return domain.Game{
		Id:          uint(oPbGame.GetId()),
		GameTypeId:  uint(oPbGame.GetGameTypeId()),
		Key:         oPbGame.GetKey(),
		Name:        oPbGame.GetName(),
		Description: oPbGame.GetDescription(),
		CreatedAt:   oPbGame.GetCreatedAt().AsTime(),
		UpdatedAt:   oPbGame.GetUpdatedAt().AsTime(),
		DeletedAt:   oPbGame.GetDeletedAt().AsTime(),
		GameType:    gameTypeFromPb(oPbGame.GetGameType()),
	}
}

func (oSelf *GameLogic) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, int64, error) {

	oRequest := &pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationInput{}

	if oPagination != nil {
		oPbPagination := &pbResource.Pagination{}
		if oPagination.Size != nil {
			oPbPagination.Size = uint64(*oPagination.Size)
		}
		if oPagination.Page != nil {
			oPbPagination.Page = uint64(*oPagination.Page)
		}
		oRequest.Pagination = oPbPagination
	}

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oFilter.Value)
		if oErr != nil {
			continue
		}

		oPbFilter := &pbResource.Filter{
			Field: *oFilter.Field,
			Value: oValue,
		}
		if oFilter.Operator != nil {
			oPbFilter.Operator = *oFilter.Operator
		}

		oRequest.Filters = append(oRequest.Filters, oPbFilter)
	}

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil || oSorter.Order == nil {
			continue
		}

		oRequest.Sorters = append(oRequest.Sorters, &pbResource.Sorter{
			Field: *oSorter.Field,
			Order: *oSorter.Order,
		})
	}

	oResponse, oErr := oSelf.ResourceLogicClient.Game.ShowGamesTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)

	aGames := make([]*domain.Game, 0, len(oResponse.GetGames()))
	for _, oOne := range oResponse.GetGames() {
		aGames = append(aGames, &domain.Game{
			Id:          uint(oOne.GetId()),
			GameTypeId:  uint(oOne.GetGameTypeId()),
			Key:         oOne.GetKey(),
			Name:        oOne.GetName(),
			Description: oOne.GetDescription(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			GameType:    gameTypeFromPb(oOne.GetGameType()),
		})
	}

	iTotal := oResponse.GetTotal()
	return aGames, int64(iTotal), oErr
}
