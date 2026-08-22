package elasticsearch

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	domain "example/internal/domain"
	elasticsearchBase "example/internal/output/application/elasticsearch"
	outputPortAnyModel "example/internal/output/port/any/model"
	pkg "example/pkg"
)

type AdminUserModel struct {
	*elasticsearchBase.AbstractElasticsearch
	Index string
}

func NewAdminUserModel(oAbstractModel *elasticsearchBase.AbstractElasticsearch) outputPortAnyModel.AdminUserModel {
	return &AdminUserModel{
		AbstractElasticsearch: oAbstractModel,
		Index:                 oAbstractModel.IndexName("admin_users"),
	}
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {
	iSize := 1

	aOptions, oErr := oSelf.IndexWheresOrdersLimitToOptions(
		oSelf.Index,
		[]map[string]any{{"term": map[string]any{"name": sName}}},
		nil,
		&elasticsearchBase.ElasticsearchLimit{Size: &iSize},
	)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	if len(oResult.Hits) == 0 {
		return nil, errors.New("資料不存在")
	}

	var oAdminUser domain.AdminUser
	if oErr := json.Unmarshal(oResult.Hits[0].Source, &oAdminUser); oErr != nil {
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(uint64(iId), 10), &oAdminUser)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound {
		return nil, errors.New("資料不存在")
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, error) {
	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.Index, aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aAdminUsers := make([]*domain.AdminUser, len(oResult.Hits))
	for i, oHit := range oResult.Hits {
		var oAdminUser domain.AdminUser
		if oErr := json.Unmarshal(oHit.Source, &oAdminUser); oErr != nil {
			return nil, oErr
		}
		aAdminUsers[i] = &oAdminUser
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)

	iTotal, oErr := oSelf.Count(oSelf.Index, aWheres)

	return iTotal, oErr
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {
	iId, oErr := oSelf.NextId("admin_user")
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminUser{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
	}

	if oAdminUser.Name != nil {
		oDoc.Name = *oAdminUser.Name
	}
	if oAdminUser.Password != nil {
		oDoc.Password = *oAdminUser.Password
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}
