package elasticsearch

import (
	"encoding/json"
	"strconv"
	"time"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

// oDeletedAtZero 跟 mysql/mongodb adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewGameModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.GameModel {
	return &GameModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("games"),
	}
}

func (oSelf *GameModel) ShowOneById(iId uint) (*domain.Game, error) {
	var oGame domain.Game

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oGame)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oGame.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oGame, nil
}

func (oSelf *GameModel) ShowOnesByFiltersWithOrdersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Game, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkg.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aGames := make([]*domain.Game, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oGame domain.Game
		if oErr := json.Unmarshal(oHit.Source, &oGame); oErr != nil {
			return nil, oErr
		}
		aGames[i] = &oGame
	}

	return aGames, nil
}

func (oSelf *GameModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)

	iTotal, oErr := oSelf.Count(oSelf.Index, aWheres)

	return iTotal, oErr
}

func (oSelf *GameModel) AddOne(oValue *domain.GameValue) (bool, error) {
	iId, oErr := oSelf.NextId("game")
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.Game{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.GameTypeId != nil {
		oDoc.GameTypeId = *oValue.GameTypeId
	}
	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}
	if oValue.Description != nil {
		oDoc.Description = *oValue.Description
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameModel) EditOneById(oValue *domain.GameValue, iId uint) (bool, error) {
	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(uint64(iId), 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)

	return bOk, oErr
}

func (oSelf *GameModel) RemoveOneById(iId uint) (bool, error) {
	sId := strconv.FormatUint(uint64(iId), 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)

	return bOk, oErr
}
