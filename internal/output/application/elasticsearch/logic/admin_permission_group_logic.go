package outputApplicationElasticsearchLogic

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
	Index                 string
	AdminPermissionsIndex string
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	return &AdminPermissionGroupLogic{
		AbstractLogic:         oAbstractLogic,
		Index:                 oAbstractLogic.IndexName("admin_permission_groups"),
		AdminPermissionsIndex: oAbstractLogic.IndexName("admin_permissions"),
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

	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		if oAdminPermissionValue == nil {
			continue
		}

		if oErr := oSelf.insertAdminPermission(oAdminPermissionValue, uint64(iId)); oErr != nil {
			return oErr
		}
	}

	return nil
}

func (oSelf *AdminPermissionGroupLogic) insertAdminPermission(oValue *domain.AdminPermissionVariable, iAdminPermissionGroupId uint64) error {
	iId, oErr := oSelf.NextId("admin_permission")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminPermission{
		Id:                     uint64(iId),
		AdminPermissionGroupId: iAdminPermissionGroupId,
		CreatedAt:              oNow,
		UpdatedAt:              oNow,
		DeletedAt:              oDeletedAtZero,
	}

	if oValue.Type != nil {
		oDoc.Type = *oValue.Type
	}
	if oValue.Key != nil {
		oDoc.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oDoc.Name = *oValue.Name
	}

	sId := strconv.FormatUint(uint64(iId), 10)
	oErr = oSelf.IndexOne(oSelf.AdminPermissionsIndex, sId, oDoc)

	return oErr
}

// existingAdminPermissions 撈出某個 admin_permission_group 底下、目前未刪除的 admin_permissions。
func (oSelf *AdminPermissionGroupLogic) existingAdminPermissions(iAdminPermissionGroupId uint64) ([]*domain.AdminPermission, error) {
	sGroupIdField := "admin_permission_group_id"
	sDeletedAtField := "deleted_at"
	aFilters := []*pkgInput.Filter{
		{Field: &sGroupIdField, Value: iAdminPermissionGroupId},
		{Field: &sDeletedAtField, Value: oDeletedAtZero},
	}

	iSize := uint(10000)
	iPage := uint(1)
	oPagination := &pkgInput.Pagination{Size: &iSize, Page: &iPage}

	aOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.AdminPermissionsIndex, aFilters, nil, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oResult, oErr := oSelf.SearchWithOptions(aOptions)
	if oErr != nil {
		return nil, oErr
	}

	aAdminPermissions := make([]*domain.AdminPermission, 0, len(oResult.Hits))
	for _, oHit := range oResult.Hits {
		var oAdminPermission domain.AdminPermission
		if oErr := json.Unmarshal(oHit.Source, &oAdminPermission); oErr != nil {
			return nil, oErr
		}
		aAdminPermissions = append(aAdminPermissions, &oAdminPermission)
	}

	return aAdminPermissions, nil
}

// EditAdminPermissionGroupById 更新 group 本身欄位，並同步 admin_permissions：
// 傳進來沒帶 id 的就新增，帶 id 的就修改；db 裡面現有、但沒出現在傳進來 id 清單裡的就刪除。
func (oSelf *AdminPermissionGroupLogic) EditAdminPermissionGroupById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_permissions")
	oColumns["updated_at"] = time.Now()

	sId := strconv.FormatUint(iId, 10)

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, sId, oColumns)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		return errors.New("0 rows updated")
	}

	aExistingAdminPermissions, oErr := oSelf.existingAdminPermissions(iId)
	if oErr != nil {
		return oErr
	}

	aKeptAdminPermissionIds := make(map[uint64]bool, len(oValue.AdminPermissions))
	for _, oAdminPermissionValue := range oValue.AdminPermissions {
		if oAdminPermissionValue == nil {
			continue
		}

		if oAdminPermissionValue.Id == nil {
			if oErr := oSelf.insertAdminPermission(oAdminPermissionValue, iId); oErr != nil {
				return oErr
			}

			continue
		}

		iAdminPermissionId := *oAdminPermissionValue.Id
		aKeptAdminPermissionIds[iAdminPermissionId] = true

		oAdminPermissionColumns, oErr := pkgUtility.StructToMap(oAdminPermissionValue)
		if oErr != nil {
			return oErr
		}

		delete(oAdminPermissionColumns, "id")
		oAdminPermissionColumns["updated_at"] = time.Now()

		sAdminPermissionId := strconv.FormatUint(iAdminPermissionId, 10)

		if _, oErr := oSelf.UpdateOne(oSelf.AdminPermissionsIndex, sAdminPermissionId, oAdminPermissionColumns); oErr != nil {
			return oErr
		}
	}

	oNow := time.Now()
	for _, oExistingAdminPermission := range aExistingAdminPermissions {
		if aKeptAdminPermissionIds[oExistingAdminPermission.Id] {
			continue
		}

		sExistingAdminPermissionId := strconv.FormatUint(oExistingAdminPermission.Id, 10)
		oPartial := map[string]any{"deleted_at": oNow}

		if _, oErr := oSelf.UpdateOne(oSelf.AdminPermissionsIndex, sExistingAdminPermissionId, oPartial); oErr != nil {
			return oErr
		}
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
