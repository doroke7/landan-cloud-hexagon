package outputApplicationElasticsearchLogic

import (
	"encoding/json"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type GameTypeLogic struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewGameTypeLogic(oAbstractLogic *elasticsearchBase.AbstractElasticsearch) outputPortAnyLogic.GameTypeLogic {
	return &GameTypeLogic{
		AbstractElasticsearch: oAbstractLogic,
		Index:                 oAbstractLogic.IndexName("game_types"),
	}
}

// ShowTree 先把所有未刪除的 game_type 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *GameTypeLogic) ShowTree() ([]*domain.GameType, error) {
	sDeletedAtField := "deleted_at"
	aFilters := []*pkgInput.Filter{{Field: &sDeletedAtField, Value: oDeletedAtZero}}

	iSize := uint(10000)
	iPage := uint(1)
	oPagination := &pkgInput.Pagination{Size: &iSize, Page: &iPage}

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, nil, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aFlat := make([]*domain.GameType, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oGameType domain.GameType
		if oErr := json.Unmarshal(oHit.Source, &oGameType); oErr != nil {
			return nil, oErr
		}
		aFlat = append(aFlat, &oGameType)
	}

	aByParent := make(map[uint][]*domain.GameType, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.GameType)
	fnAttach = func(oNode *domain.GameType) {
		for _, oChild := range aByParent[oNode.Id] {
			fnAttach(oChild)
			oNode.Children = append(oNode.Children, *oChild)
		}
	}

	aRoots := aByParent[0]
	for _, oRoot := range aRoots {
		fnAttach(oRoot)
	}

	return aRoots, nil
}

func (oSelf *GameTypeLogic) ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {
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

	aGameTypes := make([]*domain.GameType, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oGameType domain.GameType
		if oErr := json.Unmarshal(oHit.Source, &oGameType); oErr != nil {
			return nil, 0, oErr
		}
		aGameTypes = append(aGameTypes, &oGameType)
	}

	return aGameTypes, uint64(oResult.Total), nil
}
