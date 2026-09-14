package outputApplicationElasticsearchLogic

import (
	"encoding/json"
	"strconv"
	"time"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
	Index string
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	return &AdminPermissionGroupLogic{
		AbstractLogic: oAbstractLogic,
		Index:         oAbstractLogic.IndexName("admin_permission_groups"),
	}
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupVariable) error {
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

// ShowTree 先把所有未刪除的 admin_permission_group 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
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

	aFlat := make([]*domain.AdminPermissionGroup, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oAdminPermissionGroup domain.AdminPermissionGroup
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermissionGroup); oErr != nil {
			return nil, oErr
		}
		aFlat = append(aFlat, &oAdminPermissionGroup)
	}

	aByParent := make(map[uint64][]*domain.AdminPermissionGroup, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.AdminPermissionGroup)
	fnAttach = func(oNode *domain.AdminPermissionGroup) {
		for _, oChild := range aByParent[oNode.Id] {
			fnAttach(oChild)
			oNode.Children = append(oNode.Children, oChild)
		}
	}

	aRoots := aByParent[0]
	for _, oRoot := range aRoots {
		fnAttach(oRoot)
	}

	return aRoots, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {
	var oAdminPermissionGroup domain.AdminPermissionGroup

	bFound, oErr := oSelf.GetById(oSelf.Index, strconv.FormatUint(iId, 10), &oAdminPermissionGroup)
	if oErr != nil {
		return nil, oErr
	}

	if !bFound || !oAdminPermissionGroup.DeletedAt.Equal(oDeletedAtZero) {
		return nil, nil
	}

	return &oAdminPermissionGroup, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error) {
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

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oAdminPermissionGroup domain.AdminPermissionGroup
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermissionGroup); oErr != nil {
			return nil, oErr
		}
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {
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

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oAdminPermissionGroup domain.AdminPermissionGroup
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermissionGroup); oErr != nil {
			return nil, 0, oErr
		}
		aAdminPermissionGroups = append(aAdminPermissionGroups, &oAdminPermissionGroup)
	}

	return aAdminPermissionGroups, uint64(oResult.Total), nil
}
