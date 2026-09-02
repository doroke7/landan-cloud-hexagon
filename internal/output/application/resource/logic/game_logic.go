package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type GameLogic struct {
	*resourceBase.AbstractResource
}

func NewGameLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractResource: oAbstractLogic,
	}
}

func protoGameTypeToDomainGameType(oPbGameType *pbResourceLogic.GameType) domain.GameType {
	if oPbGameType == nil {
		return domain.GameType{}
	}

	oGameType := domain.GameType{
		Id:        uint(oPbGameType.GetId()),
		ParentId:  uint(oPbGameType.GetParentId()),
		Key:       oPbGameType.GetKey(),
		Name:      oPbGameType.GetName(),
		CreatedAt: oPbGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oPbGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oPbGameType.GetDeletedAt().AsTime(),
	}

	if oParent := oPbGameType.GetParent(); oParent != nil {
		oParentDomain := protoGameTypeToDomainGameType(oParent)
		oGameType.Parent = &oParentDomain
	}

	for _, oChild := range oPbGameType.GetChildren() {
		oGameType.Children = append(oGameType.Children, protoGameTypeToDomainGameType(oChild))
	}

	return oGameType
}

func protoGameToDomainGame(oPbGame *pbResourceLogic.Game) domain.Game {
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
		GameType:    protoGameTypeToDomainGameType(oPbGame.GetGameType()),
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
			Id:          uint(oOne.GetId()),
			GameTypeId:  uint(oOne.GetGameTypeId()),
			Key:         oOne.GetKey(),
			Name:        oOne.GetName(),
			Description: oOne.GetDescription(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			GameType:    protoGameTypeToDomainGameType(oOne.GetGameType()),
		})
	}

	iTotal := oResponse.GetTotal()
	return aGames, iTotal, oErr
}
