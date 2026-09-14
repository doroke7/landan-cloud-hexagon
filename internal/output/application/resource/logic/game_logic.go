package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
)

type GameLogic struct {
	*AbstractLogic
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	oLogic := &GameLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
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
		oGameType := pkgProtoToDomain.GameType(oOne.GetGameType())
		oGame := &domain.Game{
			Id:          uint64(oOne.GetId()),
			GameTypeId:  uint64(oOne.GetGameTypeId()),
			Key:         oOne.GetKey(),
			Name:        oOne.GetName(),
			Description: oOne.GetDescription(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			GameType:    oGameType,
		}
		aGames = append(aGames, oGame)
	}

	iTotal := uint64(oResponse.GetTotal())
	return aGames, iTotal, oErr
}

// ShowGameById proto GameLogic 沒有對應 rpc，改用既有的 list rpc 帶 id 過濾（size 1）。
func (oSelf *GameLogic) ShowGameById(iId uint64) (*domain.Game, error) {
	sField := "id"
	sOperator := "eq"
	iSize := uint(1)
	iPage := uint(1)

	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: float64(iId)},
	}
	oPagination := &pkgInput.Pagination{Size: &iSize, Page: &iPage}

	aGames, _, oErr := oSelf.ShowGamesTotalByFiltersWithSortersPagination(aFilters, nil, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	if len(aGames) == 0 {
		return nil, nil
	}

	return aGames[0], nil
}

// ShowGamesByGameTypeId proto GameLogic 沒有對應 rpc，改用既有的 list rpc 帶 game_type_id 過濾。
func (oSelf *GameLogic) ShowGamesByGameTypeId(iGameTypeId uint64) ([]*domain.Game, error) {
	sField := "game_type_id"
	sOperator := "eq"
	iSize := uint(10000)
	iPage := uint(1)

	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: float64(iGameTypeId)},
	}
	oPagination := &pkgInput.Pagination{Size: &iSize, Page: &iPage}

	aGames, _, oErr := oSelf.ShowGamesTotalByFiltersWithSortersPagination(aFilters, nil, oPagination)

	return aGames, oErr
}
