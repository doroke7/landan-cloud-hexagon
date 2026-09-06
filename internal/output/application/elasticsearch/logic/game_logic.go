package outputApplicationElasticsearchLogic

import (
	"encoding/json"
	"time"

	domain "example/internal/domain"
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

// oDeletedAtZero 跟 elasticsearch/model adapter 用同一個「未刪除」標記值。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameLogic struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
	Index string
}

func NewGameLogic(oAbstractLogic *outputApplicationElasticsearch.AbstractElasticsearch) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractElasticsearch: oAbstractLogic,
		Index:                 oAbstractLogic.IndexName("games"),
	}
}

// ShowGamesTotalByFiltersWithSortersPagination 一次 search 就同時拿到「這一頁」跟「符合條件的總數」
// （options 裡有 track_total_hits），不用再另外打一次 count。
func (oSelf *GameLogic) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, 0, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, 0, oErr
	}

	aGames := make([]*domain.Game, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oGame domain.Game
		if oErr := json.Unmarshal(oHit.Source, &oGame); oErr != nil {
			return nil, 0, oErr
		}
		aGames[i] = &oGame
	}

	return aGames, uint64(oResult.Total), nil
}
