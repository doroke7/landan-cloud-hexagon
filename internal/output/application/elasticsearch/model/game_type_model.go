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

type GameTypeModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewGameTypeModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("game_types"),
	}
}

func (oSelf *GameTypeModel) ShowOneById(iId uint) (*domain.GameType, error) {
	var oGameType domain.GameType

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oGameType)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oGameType.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oGameType, nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)
	aWheres = append(aWheres, map[string]any{"term": map[string]any{"deleted_at": oDeletedAtZero}})
	aOrders := oSelf.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	oResult, oErr := oSelf.Search(oSelf.Index, aWheres, aOrders, oLimit)
	if oErr != nil {
		return nil, oErr
	}

	aGameTypes := make([]*domain.GameType, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oGameType domain.GameType
		if oErr := json.Unmarshal(oHit.Source, &oGameType); oErr != nil {
			return nil, oErr
		}
		aGameTypes[i] = &oGameType
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	return oSelf.Count(oSelf.Index, oSelf.FiltersToWheres(aFilters))
}

func (oSelf *GameTypeModel) AddOne(oValue *domain.GameTypeValue) (bool, error) {
	iId, oErr := oSelf.NextId("game_type")
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.GameType{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeModel) EditOneById(oValue *domain.GameTypeValue, iId uint) (bool, error) {
	oColumns, oErr := pkg.StructToMap(oValue)
	if oErr != nil {
		return false, oErr
	}
	oColumns["updated_at"] = time.Now()

	return oSelf.UpdateOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oColumns)
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint) (bool, error) {
	return oSelf.UpdateOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), map[string]any{"deleted_at": time.Now()})
}
