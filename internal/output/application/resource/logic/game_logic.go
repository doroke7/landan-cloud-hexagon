package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type GameLogic struct {
	*AbstractLogic
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func protoGameTypeToDomainGameType(oProtoGameType *pbResource.GameType) domain.GameType {
	if oProtoGameType == nil {
		return domain.GameType{}
	}

	oGameType := domain.GameType{
		Id:        uint64(oProtoGameType.GetId()),
		ParentId:  uint64(oProtoGameType.GetParentId()),
		Key:       oProtoGameType.GetKey(),
		Name:      oProtoGameType.GetName(),
		CreatedAt: oProtoGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oProtoGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oProtoGameType.GetDeletedAt().AsTime(),
	}

	if oParent := oProtoGameType.GetParent(); oParent != nil {
		oParentDomain := protoGameTypeToDomainGameType(oParent)
		oGameType.Parent = &oParentDomain
	}

	for _, oChild := range oProtoGameType.GetChildren() {
		oGameType.Children = append(oGameType.Children, protoGameTypeToDomainGameType(oChild))
	}

	return oGameType
}

func protoGameToDomainGame(oProtoGame *pbResource.Game) domain.Game {
	if oProtoGame == nil {
		return domain.Game{}
	}

	return domain.Game{
		Id:          uint64(oProtoGame.GetId()),
		GameTypeId:  uint64(oProtoGame.GetGameTypeId()),
		Key:         oProtoGame.GetKey(),
		Name:        oProtoGame.GetName(),
		Description: oProtoGame.GetDescription(),
		CreatedAt:   oProtoGame.GetCreatedAt().AsTime(),
		UpdatedAt:   oProtoGame.GetUpdatedAt().AsTime(),
		DeletedAt:   oProtoGame.GetDeletedAt().AsTime(),
		GameType:    protoGameTypeToDomainGameType(oProtoGame.GetGameType()),
	}
}

func (oSelf *GameLogic) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {

	oRequest := &pbResourceLogic.GameShowGamesTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.Game.ShowGamesTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)

	aGames := make([]*domain.Game, 0, len(oResponse.GetGames()))
	for _, oOne := range oResponse.GetGames() {
		aGames = append(aGames, &domain.Game{
			Id:          uint64(oOne.GetId()),
			GameTypeId:  uint64(oOne.GetGameTypeId()),
			Key:         oOne.GetKey(),
			Name:        oOne.GetName(),
			Description: oOne.GetDescription(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			GameType:    protoGameTypeToDomainGameType(oOne.GetGameType()),
		})
	}

	iTotal := uint64(oResponse.GetTotal())
	return aGames, iTotal, oErr
}
