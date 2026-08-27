package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type GameModel struct {
	*resourceBase.AbstractResource
}

func NewGameModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractResource: oAbstractModel,
	}
}

func gameTypeFromPb(oPbGameType *pbResourceModel.GameType) domain.GameType {
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

func gameFromPb(oPbGame *pbResourceModel.Game) domain.Game {
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

func (oSelf *GameModel) AddOne(oGameParm *domain.GameValue) (bool, error) {

	oRequest := &pbResourceModel.GameAddOneInput{}

	if oGameParm.GameTypeId != nil {
		iGameTypeId := uint32(*oGameParm.GameTypeId)
		oRequest.GameTypeId = &iGameTypeId
	}

	oRequest.Key = oGameParm.Key
	oRequest.Name = oGameParm.Name
	oRequest.Description = oGameParm.Description

	bResult, oErr := oSelf.ResourceModelClient.Game.AddOne(oSelf.Context, oRequest)
	_ = bResult

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Game.ShowOneById(
		oSelf.Context,
		&pbResourceModel.GameShowOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	if oResponse.GetId() == 0 {
		return nil, nil
	}

	return &domain.Game{
		Id:          uint(oResponse.GetId()),
		GameTypeId:  uint(oResponse.GetGameTypeId()),
		Key:         oResponse.GetKey(),
		Name:        oResponse.GetName(),
		Description: oResponse.GetDescription(),
		CreatedAt:   oResponse.GetCreatedAt().AsTime(),
		UpdatedAt:   oResponse.GetUpdatedAt().AsTime(),
		DeletedAt:   oResponse.GetDeletedAt().AsTime(),
		GameType:    gameTypeFromPb(oResponse.GetGameType()),
	}, nil
}

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, error) {

	oRequest := &pbResourceModel.GameShowOnesByFiltersWithOrdersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Game.ShowOnesByFiltersWithOrdersPagination(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

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

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Game.TotalByFilters(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.GameEditOneByIdInput{Id: uint32(iId)}

	if oGame.GameTypeId != nil {
		iGameTypeId := uint32(*oGame.GameTypeId)
		oRequest.GameTypeId = &iGameTypeId
	}

	oRequest.Key = oGame.Key
	oRequest.Name = oGame.Name
	oRequest.Description = oGame.Description

	oResponse, oErr := oSelf.ResourceModelClient.Game.EditOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *GameModel) RemoveOneById(iId uint) (bool, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Game.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.GameRemoveOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}
