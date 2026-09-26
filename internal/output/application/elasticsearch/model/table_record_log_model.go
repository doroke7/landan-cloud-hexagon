package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	domain "example/internal/domain"
	outputApplicationElasticsearch "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
)

type TableRecordLogModel struct {
	*outputApplicationElasticsearch.AbstractElasticsearch
	Index string
}

func NewTableRecordLogModel(oAbstractElasticsearch *outputApplicationElasticsearch.AbstractElasticsearch) outputPortAnyModel.TableRecordLogModel {
	return &TableRecordLogModel{
		AbstractElasticsearch: oAbstractElasticsearch,
		Index:                 oAbstractElasticsearch.IndexName("table_record_logs"),
	}
}

func (oSelf *TableRecordLogModel) AddOne(oValue *domain.TableRecordLogVariable) error {
	iId, oErr := oSelf.NextId("table_record_log")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.TableRecordLog{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.GameId != nil {
		oDoc.GameId = *oValue.GameId
	}
	if oValue.TableRecordId != nil {
		oDoc.TableRecordId = *oValue.TableRecordId
	}
	if oValue.State != nil {
		oDoc.State = *oValue.State
	}
	if oValue.Text != nil {
		oDoc.Text = *oValue.Text
	}
	if oValue.Image != nil {
		oDoc.Image = *oValue.Image
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *TableRecordLogModel) ShowOneById(iId uint64) (*domain.TableRecordLog, error) {
	var oTableRecordLog domain.TableRecordLog

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oTableRecordLog)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oTableRecordLog.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oTableRecordLog, nil
}

func (oSelf *TableRecordLogModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.TableRecordLog, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aTableRecordLogs := make([]*domain.TableRecordLog, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oTableRecordLog domain.TableRecordLog
		if oErr := json.Unmarshal(oHit.Source, &oTableRecordLog); oErr != nil {
			return nil, oErr
		}
		aTableRecordLogs[i] = &oTableRecordLog
	}

	return aTableRecordLogs, nil
}

func (oSelf *TableRecordLogModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	sDeletedAtField := "deleted_at"
	aFilters = append(aFilters, &pkgInput.Filter{Field: &sDeletedAtField, Value: oDeletedAtZero})

	_, aFilterClauses := oSelf.FiltersToMustFilter(aFilters)

	oQuery := map[string]any{"match_all": map[string]any{}}
	if len(aFilterClauses) > 0 {
		oQuery = map[string]any{"bool": map[string]any{"filter": aFilterClauses}}
	}

	aBodyBytes, oErr := json.Marshal(map[string]any{"query": oQuery})
	if oErr != nil {
		return 0, oErr
	}

	aOptions := []func(*esapi.CountRequest){
		oSelf.Client.Count.WithContext(oSelf.Context),
		oSelf.Client.Count.WithIndex(oSelf.Index),
		oSelf.Client.Count.WithBody(bytes.NewReader(aBodyBytes)),
	}

	iTotal, oErr := oSelf.CountWithOptions(aOptions)

	return uint64(iTotal), oErr
}
