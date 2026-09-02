package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
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

func protoGameTypeToDomainGameType(oProtoGameType *pbResource.GameType) domain.GameType {
	if oProtoGameType == nil {
		return domain.GameType{}
	}

	oGameType := domain.GameType{
		Id:        uint(oProtoGameType.GetId()),
		ParentId:  uint(oProtoGameType.GetParentId()),
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
		GameType:    protoGameTypeToDomainGameType(oResponse.GetGameType()),
	}, nil
}

func (oSelf *GameModel) ShowOneByKey(sKey string) (*domain.Game, error) {

	oResponse, oErr := oSelf.ResourceModelClient.Game.ShowOneByKey(
		oSelf.Context,
		&pbResourceModel.GameShowOneByKeyInput{Key: sKey},
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
		GameType:    protoGameTypeToDomainGameType(oResponse.GetGameType()),
	}, nil
}

func (oSelf *GameModel) ShowOnesByGameTypeId(iGameTypeId uint) ([]*domain.Game, error) {
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

func (oSelf *GameModel) TotalByGameTypeId(iGameTypeId uint) (uint64, error) {
	sField := "game_type_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iGameTypeId},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
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
			GameType:    protoGameTypeToDomainGameType(oOne.GetGameType()),
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
