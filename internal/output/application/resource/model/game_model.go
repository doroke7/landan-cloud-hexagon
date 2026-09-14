package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type GameModel struct {
	*AbstractModel
}

func NewGameModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameModel) AddOne(oGameVariable *domain.GameVariable) error {

	oRequest := &pbResourceModel.GameAddOneInput{
		Value: pkgDomainToProto.GameValue(oGameVariable),
	}

	_, oErr := oSelf.ResourceModelClient.Game.AddOne(oSelf.Context, oRequest)

	return oErr
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

	oGame := pkgProtoToDomain.Game(oProtoGame)

	return &oGame, nil
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

func (oSelf *GameModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.Game.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *GameModel) EditOneById(oGame *domain.GameVariable, iId uint64) error {

	oRequest := &pbResourceModel.GameEditOneByIdInput{
		Id:    iId,
		Value: pkgDomainToProto.GameValue(oGame),
	}

	_, oErr := oSelf.ResourceModelClient.Game.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.Game.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.GameRemoveOneByIdInput{Id: iId},
	)

	return oErr
}
