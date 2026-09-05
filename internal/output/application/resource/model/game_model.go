package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

func domainGameValueToProtoGameVariable(oGame *domain.GameValue) *pbResourceModel.GameVariable {
	oVariable := &pbResourceModel.GameVariable{
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
	}

	if oGame.GameTypeId != nil {
		iGameTypeId := uint64(*oGame.GameTypeId)
		oVariable.GameTypeId = &iGameTypeId
	}

	return oVariable
}

func protoGameToDomainGame(oProtoGame *pbResource.Game) domain.Game {
	if oProtoGame == nil {
		return domain.Game{}
	}

	return domain.Game{
		Id:          uint(oProtoGame.GetId()),
		GameTypeId:  uint(oProtoGame.GetGameTypeId()),
		Key:         oProtoGame.GetKey(),
		Name:        oProtoGame.GetName(),
		Description: oProtoGame.GetDescription(),
		CreatedAt:   oProtoGame.GetCreatedAt().AsTime(),
		UpdatedAt:   oProtoGame.GetUpdatedAt().AsTime(),
		DeletedAt:   oProtoGame.GetDeletedAt().AsTime(),
		GameType:    protoGameTypeToDomainGameType(oProtoGame.GetGameType()),
	}
}

type GameModel struct {
	*resourceBase.AbstractResource
}

func NewGameModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractResource: oAbstractModel,
	}
}

func (oSelf *GameModel) AddOne(oGameParm *domain.GameValue) error {

	oRequest := &pbResourceModel.GameAddOneInput{
		Variable: domainGameValueToProtoGameVariable(oGameParm),
	}

	_, oErr := oSelf.ResourceModelClient.Game.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameModel) ShowOneById(iId uint64) (*domain.Game, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Game.ShowOneById(
		oSelf.Context,
		&pbResourceModel.GameShowOneByIdInput{Id: uint64(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoGame := oResponse.GetGame()
	if oProtoGame.GetId() == 0 {
		return nil, nil
	}

	oGame := protoGameToDomainGame(oProtoGame)

	return &oGame, nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Game.ShowOneByKey(
		oSelf.Context,
		&pbResourceModel.GameShowOneByKeyInput{Key: sKey},
	)

	if oErr != nil {
		return nil, oErr
	}

	oProtoGame := oResponse.GetGame()
	if oProtoGame.GetId() == 0 {
		return nil, nil
	}

	oGame := protoGameToDomainGame(oProtoGame)

	return &oGame, nil
}

func (oSelf *GameModel) ShowOnesByGameTypeId(iGameTypeId uint64) ([]*domain.Game, error) {
	sField := "game_type_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iGameTypeId},
	}

	iSize := uint(10000)
	iPage := uint(1)
	aGames, oErr := oSelf.ShowOnesByFiltersWithOrdersPagination(aFilters, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
	return aGames, oErr
}

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint64) (uint, error) {
	sField := "game_type_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iGameTypeId},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return uint(iTotal), oErr
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
		oGame := protoGameToDomainGame(oOne)
		aGames = append(aGames, &oGame)
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Game.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameValue, iId uint64) error {

	oRequest := &pbResourceModel.GameEditOneByIdInput{
		Id:       uint32(iId),
		Variable: domainGameValueToProtoGameVariable(oGame),
	}

	_, oErr := oSelf.ResourceModelClient.Game.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.Game.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.GameRemoveOneByIdInput{Id: uint32(iId)},
	)

	return oErr
}
