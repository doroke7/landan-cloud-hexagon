package outputApplicationElasticsearch

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
	pkgInput "example/pkg/input"
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
	"eq":          true,
	"ne":          true,
	"gt":          true,
	"gte":         true,
	"lt":          true,
	"lte":         true,
	"contains":    true,
	"notContains": true,
	"startsWith":  true,
	"endsWith":    true,
	"in":          true,
	"notIn":       true,
	"between":     true,
	"match":       true,
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

func (oSelf *AbstractElasticsearch) PaginationToFrom(oPagination *pkgInput.Pagination) int {
	iSize := uint(10)
	iPage := uint(1)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	if oPagination != nil && oPagination.Page != nil && *oPagination.Page != 0 {
		iPage = *oPagination.Page
	}

	iFrom := int((iPage - 1) * iSize)

	return iFrom
}

func (oSelf *AbstractElasticsearch) PaginationToSize(oPagination *pkgInput.Pagination) int {
	iSize := uint(10)

	if oPagination != nil && oPagination.Size != nil && *oPagination.Size != 0 {
		iSize = *oPagination.Size
	}

	return int(iSize)
}

func (oSelf *AbstractElasticsearch) FiltersToMustFilter(aFilters []*pkgInput.Filter) ([]map[string]any, []map[string]any) {
	aFilterClauses := make([]map[string]any, 0, len(aFilters))
	aMustClauses := make([]map[string]any, 0, len(aFilters))

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
			aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"term": map[string]any{sField: oValue}}}})
		case "gt", "gte", "lt", "lte":
			aFilterClauses = append(aFilterClauses, map[string]any{"range": map[string]any{sField: map[string]any{sOperator: oValue}}})
		case "contains":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}})
			}
		case "notContains":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue + "*", "case_insensitive": true}}}}})
			}
		case "startsWith":
			if sValue, bOk := oValue.(string); bOk {
				aFilterClauses = append(aFilterClauses, map[string]any{"prefix": map[string]any{sField: map[string]any{"value": sValue, "case_insensitive": true}}})
			}
		case "endsWith":
			if sValue, bOk := oValue.(string); bOk {
				// match     = 按照文字意思 搜尋
				// wildcard  = 按照字串格式 搜尋
				aFilterClauses = append(aFilterClauses, map[string]any{"wildcard": map[string]any{sField: map[string]any{"value": "*" + sValue, "case_insensitive": true}}})
			}
		case "in":
			aFilterClauses = append(aFilterClauses, map[string]any{"terms": map[string]any{sField: oValue}})
		case "notIn":
			aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"terms": map[string]any{sField: oValue}}}})
		case "between":
			if aRange, bOk := oValue.([]any); bOk && len(aRange) == 2 {
				aFilterClauses = append(aFilterClauses, map[string]any{"range": map[string]any{sField: map[string]any{"gte": aRange[0], "lte": aRange[1]}}})
			}
		case "match":

			if sValue, bOk := oValue.(string); bOk {
				aMustClauses = append(aMustClauses, map[string]any{"match": map[string]any{sField: sValue}})
			}
		default:
			aFilterClauses = append(aFilterClauses, map[string]any{"term": map[string]any{sField: oValue}})
		}
	}

	return aMustClauses, aFilterClauses
}

func (oSelf *AbstractElasticsearch) SortersToSort(aSorters []*pkgInput.Sorter) []map[string]any {
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

func (oSelf *AbstractElasticsearch) IndexFiltersSortersPaginationToOptions(sIndex string, aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]func(*esapi.SearchRequest), error) {

	aMustClauses, aFilterClauses := oSelf.FiltersToMustFilter(aFilters)
	aSorts := oSelf.SortersToSort(aSorters)
	iFrom := oSelf.PaginationToFrom(oPagination)
	iSize := oSelf.PaginationToSize(oPagination)

	oBoolQuery := map[string]any{}
	if len(aFilterClauses) > 0 {
		oBoolQuery["filter"] = aFilterClauses
	}
	if len(aMustClauses) > 0 {
		oBoolQuery["must"] = aMustClauses
	}

	oQuery := map[string]any{"match_all": map[string]any{}}
	if len(oBoolQuery) > 0 {
		oQuery = map[string]any{"bool": oBoolQuery}
	}

	oBody := map[string]any{"query": oQuery}
	if len(aSorts) > 0 {
		oBody["sort"] = aSorts
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
		oSelf.Client.Search.WithFrom(iFrom),
		oSelf.Client.Search.WithSize(iSize),
	}

	return aOptions, nil
}

func (oSelf *AbstractElasticsearch) SearchWithOptions(aOptions []func(*esapi.SearchRequest)) (*ElasticsearchSearchResult, error) {

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
		return nil, fmt.Errorf("elasticsearch search failed: %s", string(aResponseBody))
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
func (oSelf *AbstractElasticsearch) CountWithOptions(aOptions []func(*esapi.CountRequest)) (uint64, error) {
	oResponse, oErr := oSelf.Client.Count(aOptions...)
	if oErr != nil {
		return 0, oErr
	}
	defer oResponse.Body.Close()

	aResponseBody, oErr := io.ReadAll(oResponse.Body)
	if oErr != nil {
		return 0, oErr
	}

	if oResponse.IsError() {
		return 0, fmt.Errorf("elasticsearch count failed: %s", string(aResponseBody))
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
