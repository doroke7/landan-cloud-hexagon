package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	bootstrap "example/bootstrap"
	pkg "example/pkg"
)

// Context 是程序等級的全局 ctx（來源是 cmd/xx.go），跟 pkg.Aop、cache/memory 的
// AbstractRepository 做法一致。
type AbstractElasticsearch struct {
	Client  *elasticsearch.Client
	Context context.Context
}

func NewAbstractElasticsearch(oContext context.Context, oClient *elasticsearch.Client) *AbstractElasticsearch {
	return &AbstractElasticsearch{
		Client:  oClient,
		Context: oContext,
	}
}

// IndexName 沿用 CONFIG.DATABASE.PREFIX，跟 mysql table 前綴一致。
func (oSelf *AbstractElasticsearch) IndexName(sName string) string {
	return bootstrap.CONFIG.DATABASE.PREFIX + sName
}

var oOperatorMap = map[string]bool{
	"eq": true, "ne": true, "gt": true, "gte": true, "lt": true, "lte": true,
	"contains": true, "notContains": true, "startsWith": true, "endsWith": true,
	"in": true, "notIn": true, "between": true,
}

func (oSelf *AbstractElasticsearch) FiltersToWheres(aFilters []*pkg.Filter) []map[string]any {
	aWheres := make([]map[string]any, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		sField := *oFilter.Field
		oValue := oFilter.Value

		sOperator := "eq"
		if oFilter.Operator != nil && oOperatorMap[*oFilter.Operator] {
			sOperator = *oFilter.Operator
		}

		switch sOperator {
		case "ne":
			aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"term": map[string]any{sField: oValue}}}})
		case "gt", "gte", "lt", "lte":
			aWheres = append(aWheres, map[string]any{"range": map[string]any{sField: map[string]any{sOperator: oValue}}})
		case "contains":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}})
			}
		case "notContains":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}}}})
			}
		case "startsWith":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"prefix": map[string]any{sField: map[string]any{"value": sValue, "case_insensitive": true}}})
			}
		case "endsWith":
			if sValue, bOk := oValue.(string); bOk {
				aWheres = append(aWheres, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue, "case_insensitive": true}}})
			}
		case "in":
			aWheres = append(aWheres, map[string]any{"terms": map[string]any{sField: oValue}})
		case "notIn":
			aWheres = append(aWheres, map[string]any{"bool": map[string]any{"must_not": map[string]any{"terms": map[string]any{sField: oValue}}}})
		case "between":
			if aRange, bOk := oValue.([]any); bOk && len(aRange) == 2 {
				aWheres = append(aWheres, map[string]any{"range": map[string]any{sField: map[string]any{"gte": aRange[0], "lte": aRange[1]}}})
			}
		default:
			aWheres = append(aWheres, map[string]any{"term": map[string]any{sField: oValue}})
		}
	}

	return aWheres
}

func (oSelf *AbstractElasticsearch) SortersToOrders(aSorters []*pkg.Sorter) []map[string]any {
	aOrders := make([]map[string]any, 0, len(aSorters))

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil {
			continue
		}

		sDirection := "asc"
		if oSorter.Order != nil && strings.EqualFold(*oSorter.Order, "desc") {
			sDirection = "desc"
		}

		aOrders = append(aOrders, map[string]any{*oSorter.Field: map[string]any{"order": sDirection}})
	}

	return aOrders
}

type ElasticsearchLimit struct {
	From *int
	Size *int
}

func (oSelf *AbstractElasticsearch) PaginationToLimit(oPagination *pkg.Pagination) *ElasticsearchLimit {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iFrom := int((iPage - 1) * iSize)
	iSizeInt := int(iSize)

	return &ElasticsearchLimit{
		From: &iFrom,
		Size: &iSizeInt,
	}
}

// ---- 查詢/寫入的共用底層操作，取代 mysql 版本裡 *gorm.DB 幫忙做的事 ----

type ElasticsearchHit struct {
	Id     string          `json:"_id"`
	Source json.RawMessage `json:"_source"`
}

type ElasticsearchSearchResult struct {
	Total int64
	Hits  []ElasticsearchHit
}

func wheresToQuery(aWheres []map[string]any) map[string]any {
	if len(aWheres) == 0 {
		return map[string]any{"match_all": map[string]any{}}
	}

	return map[string]any{"bool": map[string]any{"must": aWheres}}
}

func (oSelf *AbstractElasticsearch) IndexWheresOrdersLimitToOptions(sIndex string, aWheres []map[string]any, aOrders []map[string]any, oLimit *ElasticsearchLimit) ([]func(*esapi.SearchRequest), error) {
	oBody := map[string]any{"query": wheresToQuery(aWheres)}
	if len(aOrders) > 0 {
		oBody["sort"] = aOrders
	}

	aBodyBytes, oErr := json.Marshal(oBody)
	if oErr != nil {
		return nil, oErr
	}

	aOptions := []func(*esapi.SearchRequest){
		oSelf.Client.Search.WithContext(oSelf.Context),
		oSelf.Client.Search.WithIndex(sIndex),
		oSelf.Client.Search.WithBody(bytes.NewReader(aBodyBytes)),
		oSelf.Client.Search.WithTrackTotalHits(true),
	}

	if oLimit != nil {
		if oLimit.From != nil {
			aOptions = append(aOptions, oSelf.Client.Search.WithFrom(*oLimit.From))
		}
		if oLimit.Size != nil {
			aOptions = append(aOptions, oSelf.Client.Search.WithSize(*oLimit.Size))
		}
	}
	return aOptions, nil
}

func (oSelf *AbstractElasticsearch) IndexFiltersSortersPaginationToOptions(sIndex string, aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]func(*esapi.SearchRequest), error) {

	aWheres := oSelf.FiltersToWheres(aFilters)
	aOrders := oSelf.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	aOptions, oErr := oSelf.IndexWheresOrdersLimitToOptions(sIndex, aWheres, aOrders, oLimit)

	return aOptions, oErr
}

func (oSelf *AbstractElasticsearch) Search(sIndex string, aWheres []map[string]any, aOrders []map[string]any, oLimit *ElasticsearchLimit) (*ElasticsearchSearchResult, error) {

	aOptions, oErr := oSelf.IndexWheresOrdersLimitToOptions(sIndex, aWheres, aOrders, oLimit)
	if oErr != nil {
		return nil, oErr
	}

	oResponse, oErr := oSelf.Client.Search(aOptions...)
	if oErr != nil {
		return nil, oErr
	}
	defer oResponse.Body.Close()

	aResponseBody, oErr := io.ReadAll(oResponse.Body)
	if oErr != nil {
		return nil, oErr
	}

	if oResponse.IsError() {
		return nil, fmt.Errorf("elasticsearch search %s failed: %s", sIndex, string(aResponseBody))
	}

	var oResult struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []ElasticsearchHit `json:"hits"`
		} `json:"hits"`
	}

	if oErr := json.Unmarshal(aResponseBody, &oResult); oErr != nil {
		return nil, oErr
	}

	return &ElasticsearchSearchResult{
		Total: oResult.Hits.Total.Value,
		Hits:  oResult.Hits.Hits,
	}, nil
}

func (oSelf *AbstractElasticsearch) Count(sIndex string, aWheres []map[string]any) (uint64, error) {
	aBodyBytes, oErr := json.Marshal(map[string]any{"query": wheresToQuery(aWheres)})
	if oErr != nil {
		return 0, oErr
	}

	oResponse, oErr := oSelf.Client.Count(
		oSelf.Client.Count.WithContext(oSelf.Context),
		oSelf.Client.Count.WithIndex(sIndex),
		oSelf.Client.Count.WithBody(bytes.NewReader(aBodyBytes)),
	)
	if oErr != nil {
		return 0, oErr
	}
	defer oResponse.Body.Close()

	aResponseBody, oErr := io.ReadAll(oResponse.Body)
	if oErr != nil {
		return 0, oErr
	}

	if oResponse.IsError() {
		return 0, fmt.Errorf("elasticsearch count %s failed: %s", sIndex, string(aResponseBody))
	}

	var oResult struct {
		Count uint64 `json:"count"`
	}
	if oErr := json.Unmarshal(aResponseBody, &oResult); oErr != nil {
		return 0, oErr
	}

	return oResult.Count, nil
}

func (oSelf *AbstractElasticsearch) GetById(sIndex string, sId string, oOut any) (bool, error) {
	oResponse, oErr := oSelf.Client.Get(
		sIndex,
		sId,
		oSelf.Client.Get.WithContext(oSelf.Context),
	)
	if oErr != nil {
		return false, oErr
	}
	defer oResponse.Body.Close()

	if oResponse.StatusCode == 404 {
		return false, nil
	}

	aResponseBody, oErr := io.ReadAll(oResponse.Body)
	if oErr != nil {
		return false, oErr
	}

	if oResponse.IsError() {
		return false, fmt.Errorf("elasticsearch get %s/%s failed: %s", sIndex, sId, string(aResponseBody))
	}

	var oResult struct {
		Found  bool            `json:"found"`
		Source json.RawMessage `json:"_source"`
	}
	if oErr := json.Unmarshal(aResponseBody, &oResult); oErr != nil {
		return false, oErr
	}

	if !oResult.Found {
		return false, nil
	}

	if oErr := json.Unmarshal(oResult.Source, oOut); oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *AbstractElasticsearch) IndexOne(sIndex string, sId string, oDoc any) error {
	aBodyBytes, oErr := json.Marshal(oDoc)
	if oErr != nil {
		return oErr
	}

	oResponse, oErr := oSelf.Client.Index(
		sIndex,
		bytes.NewReader(aBodyBytes),
		oSelf.Client.Index.WithContext(oSelf.Context),
		oSelf.Client.Index.WithDocumentID(sId),
		oSelf.Client.Index.WithRefresh("true"),
	)
	if oErr != nil {
		return oErr
	}
	defer oResponse.Body.Close()

	if oResponse.IsError() {
		aResponseBody, _ := io.ReadAll(oResponse.Body)
		return fmt.Errorf("elasticsearch index %s/%s failed: %s", sIndex, sId, string(aResponseBody))
	}

	return nil
}

func (oSelf *AbstractElasticsearch) UpdateOne(sIndex string, sId string, oPartial map[string]any) (bool, error) {
	aBodyBytes, oErr := json.Marshal(map[string]any{"doc": oPartial})
	if oErr != nil {
		return false, oErr
	}

	oResponse, oErr := oSelf.Client.Update(
		sIndex,
		sId,
		bytes.NewReader(aBodyBytes),
		oSelf.Client.Update.WithContext(oSelf.Context),
		oSelf.Client.Update.WithRefresh("true"),
	)
	if oErr != nil {
		return false, oErr
	}
	defer oResponse.Body.Close()

	if oResponse.StatusCode == 404 {
		return false, nil
	}

	if oResponse.IsError() {
		aResponseBody, _ := io.ReadAll(oResponse.Body)
		return false, fmt.Errorf("elasticsearch update %s/%s failed: %s", sIndex, sId, string(aResponseBody))
	}

	return true, nil
}

// NextId 用 painless script 對 counters index 做原子遞增，取得跟其他 adapter
// 一致的 uint 自增 id（模仿 mongodb adapter 用 counters collection 累加的做法）。
func (oSelf *AbstractElasticsearch) NextId(sName string) (uint, error) {
	sBody := `{"script":{"source":"ctx._source.seq += 1","lang":"painless"},"upsert":{"seq":1}}`

	oResponse, oErr := oSelf.Client.Update(
		oSelf.IndexName("counters"),
		sName,
		strings.NewReader(sBody),
		oSelf.Client.Update.WithContext(oSelf.Context),
		oSelf.Client.Update.WithSource("true"),
		oSelf.Client.Update.WithRetryOnConflict(5),
	)
	if oErr != nil {
		return 0, oErr
	}
	defer oResponse.Body.Close()

	aResponseBody, oErr := io.ReadAll(oResponse.Body)
	if oErr != nil {
		return 0, oErr
	}

	if oResponse.IsError() {
		return 0, fmt.Errorf("elasticsearch next id %s failed: %s", sName, string(aResponseBody))
	}

	var oResult struct {
		Get struct {
			Source struct {
				Seq uint `json:"seq"`
			} `json:"_source"`
		} `json:"get"`
	}
	if oErr := json.Unmarshal(aResponseBody, &oResult); oErr != nil {
		return 0, oErr
	}

	return oResult.Get.Source.Seq, nil
}
