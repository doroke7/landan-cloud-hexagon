package outputApplicationElasticsearchModel

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionGroupModel struct {
	*AbstractModel
	Index string
}

func NewAdminPermissionGroupModel(oAbstractModel *AbstractModel) outputPortAnyModel.AdminPermissionGroupModel {
	return &AdminPermissionGroupModel{
		AbstractModel: oAbstractModel,
		Index:         oAbstractModel.IndexName("admin_permission_groups"),
	}
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByParentId(iParentId uint64) ([]*domain.AdminPermissionGroup, error) {
	sParentIdField := "parent_id"
	sOperator := "eq"
	sDeletedAtField := "deleted_at"

	aFilters := []*pkgInput.Filter{
		{Field: &sParentIdField, Operator: &sOperator, Value: iParentId},
		{Field: &sDeletedAtField, Value: oDeletedAtZero},
	}

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

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oAdminPermissionGroup domain.AdminPermissionGroup
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermissionGroup); oErr != nil {
			return nil, oErr
		}
		aAdminPermissionGroups[i] = &oAdminPermissionGroup
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error) {
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

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oAdminPermissionGroup domain.AdminPermissionGroup
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermissionGroup); oErr != nil {
			return nil, oErr
		}
		aAdminPermissionGroups[i] = &oAdminPermissionGroup
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	aFilterClauses := make([]map[string]any, 0, len(aFilters))

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		sField := *oFilter.Field
		oValue := oFilter.Value

		sOperator := "eq"
		if oFilter.Operator != nil {
			sOperator = *oFilter.Operator
		}

		switch sOperator {
		case "ne":
			aFilterClauses = append(aFilterClauses, map[string]any{"bool": map[string]any{"must_not": map[string]any{"term": map[string]any{sField: oValue}}}})
		case "gt", "gte", "lt", "lte":
			aFilterClauses = append(aFilterClauses, map[string]any{"range": map[string]any{sField: map[string]any{sOperator: oValue}}})
		case "in":
			aFilterClauses = append(aFilterClauses, map[string]any{"terms": map[string]any{sField: oValue}})
		default:
			aFilterClauses = append(aFilterClauses, map[string]any{"term": map[string]any{sField: oValue}})
		}
	}

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

func (oSelf *AdminPermissionGroupModel) AddOne(oValue *domain.AdminPermissionGroupVariable) error {
	iId, oErr := oSelf.NextId("admin_permission_group")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminPermissionGroup{
		Id:        uint64(iId),
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
		return oErr
	}

	return nil
}

func (oSelf *AdminPermissionGroupModel) EditOneById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(iId, 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *AdminPermissionGroupModel) RemoveOneById(iId uint64) error {
	sId := strconv.FormatUint(iId, 10)
	oPartial := map[string]any{"deleted_at": time.Now()}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oPartial)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows deleted")
	}

	return nil
}
