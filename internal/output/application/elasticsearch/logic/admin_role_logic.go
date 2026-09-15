package outputApplicationElasticsearchLogic

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminRoleLogic struct {
	*AbstractLogic
	Index                             string
	AdminRolesToAdminPermissionsIndex string
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	return &AdminRoleLogic{
		AbstractLogic:                     oAbstractLogic,
		Index:                             oAbstractLogic.IndexName("admin_roles"),
		AdminRolesToAdminPermissionsIndex: oAbstractLogic.IndexName("admin_roles_to_admin_permissions"),
	}
}

// searchAdminRolesToAdminPermissions 查出某個 admin_role_id 底下現有的所有關聯。
func (oSelf *AdminRoleLogic) searchAdminRolesToAdminPermissions(iAdminRoleId uint64) ([]domain.AdminRolesToAdminPermission, error) {
	sAdminRoleIdField := "admin_role_id"
	iSearchSize := uint(10000)
	iSearchPage := uint(1)
	aFilters := []*pkgInput.Filter{{Field: &sAdminRoleIdField, Value: iAdminRoleId}}
	oPagination := &pkgInput.Pagination{Size: &iSearchSize, Page: &iSearchPage}

	aSearchOptions, oErr := oSelf.IndexFiltersSortersPaginationToOptions(oSelf.AdminRolesToAdminPermissionsIndex, aFilters, nil, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	oSearchResult, oErr := oSelf.SearchWithOptions(aSearchOptions)
	if oErr != nil {
		return nil, oErr
	}

	aExistingAdminRolesToAdminPermissions := make([]domain.AdminRolesToAdminPermission, 0, len(oSearchResult.Hits))
	for _, oHit := range oSearchResult.Hits {
		var oExisting domain.AdminRolesToAdminPermission
		if oErr := json.Unmarshal(oHit.Source, &oExisting); oErr != nil {
			return nil, oErr
		}

		aExistingAdminRolesToAdminPermissions = append(aExistingAdminRolesToAdminPermissions, oExisting)
	}

	return aExistingAdminRolesToAdminPermissions, nil
}

// deleteAdminRolesToAdminPermissions 清掉某個 admin_role_id 底下現有的所有關聯，
// 避免髒數據（例如 id 被重用）導致異常。
func (oSelf *AdminRoleLogic) deleteAdminRolesToAdminPermissions(iAdminRoleId uint64) error {
	aExistingAdminRolesToAdminPermissions, oErr := oSelf.searchAdminRolesToAdminPermissions(iAdminRoleId)
	if oErr != nil {
		return oErr
	}

	for _, oExisting := range aExistingAdminRolesToAdminPermissions {
		sRelationId := strconv.FormatUint(iAdminRoleId, 10) + "_" + strconv.FormatUint(oExisting.AdminPermissionId, 10)

		if oErr := oSelf.DeleteOne(oSelf.AdminRolesToAdminPermissionsIndex, sRelationId); oErr != nil {
			return oErr
		}
	}

	return nil
}

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	iId, oErr := oSelf.NextId("admin_role")
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminRole{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
	}

	if oVariable.Key != nil {
		oDoc.Key = *oVariable.Key
	}
	if oVariable.Name != nil {
		oDoc.Name = *oVariable.Name
	}

	if oErr := oSelf.IndexOne(oSelf.Index, strconv.FormatUint(uint64(iId), 10), oDoc); oErr != nil {
		return oErr
	}

	// 避免髒數據（例如 id 被重用）導致異常，插入前先把這個 admin_role_id 底下的關聯清乾淨
	if oErr := oSelf.deleteAdminRolesToAdminPermissions(uint64(iId)); oErr != nil {
		return oErr
	}

	if oVariable.AdminPermissionIds == nil || len(*oVariable.AdminPermissionIds) == 0 {
		return nil
	}

	for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
		oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
			AdminRoleId:       uint64(iId),
			AdminPermissionId: iAdminPermissionId,
		}

		sRelationId := strconv.FormatUint(uint64(iId), 10) + "_" + strconv.FormatUint(iAdminPermissionId, 10)

		if oErr := oSelf.IndexOne(oSelf.AdminRolesToAdminPermissionsIndex, sRelationId, oAdminRolesToAdminPermission); oErr != nil {
			return oErr
		}
	}

	return nil
}

func (oSelf *AdminRoleLogic) EditAdminRoleById(oVariable *domain.AdminRoleVariable, iId uint64) error {
	oPartial := map[string]any{"updated_at": time.Now()}

	if oVariable.Key != nil {
		oPartial["key"] = *oVariable.Key
	}
	if oVariable.Name != nil {
		oPartial["name"] = *oVariable.Name
	}

	bOk, oErr := oSelf.UpdateOne(oSelf.Index, strconv.FormatUint(iId, 10), oPartial)
	if oErr != nil {
		return oErr
	}

	if !bOk {
		oZeroRowsError := errors.New("0 rows updated")

		return oZeroRowsError
	}

	if oVariable.AdminPermissionIds == nil {
		return nil
	}

	aExistingAdminRolesToAdminPermissions, oErr := oSelf.searchAdminRolesToAdminPermissions(iId)
	if oErr != nil {
		return oErr
	}

	oMapExistingAdminPermissionIds := make(map[uint64]bool, len(aExistingAdminRolesToAdminPermissions))
	for _, oExisting := range aExistingAdminRolesToAdminPermissions {
		oMapExistingAdminPermissionIds[oExisting.AdminPermissionId] = true
	}

	oInputtingAdminPermissionIds := make(map[uint64]bool, len(*oVariable.AdminPermissionIds))
	for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
		oInputtingAdminPermissionIds[iAdminPermissionId] = true
	}

	// a. 傳進來的 id 在 ES 不存在 -> 插入
	for iAdminPermissionId := range oInputtingAdminPermissionIds {
		if oMapExistingAdminPermissionIds[iAdminPermissionId] {
			continue
		}

		oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
			AdminRoleId:       iId,
			AdminPermissionId: iAdminPermissionId,
		}

		sRelationId := strconv.FormatUint(iId, 10) + "_" + strconv.FormatUint(iAdminPermissionId, 10)

		if oErr := oSelf.IndexOne(oSelf.AdminRolesToAdminPermissionsIndex, sRelationId, oAdminRolesToAdminPermission); oErr != nil {
			return oErr
		}
	}

	// b. ES 存在但不在傳進來的 id 裡面 -> 刪除
	for iAdminPermissionId := range oMapExistingAdminPermissionIds {
		if oInputtingAdminPermissionIds[iAdminPermissionId] {
			continue
		}

		sRelationId := strconv.FormatUint(iId, 10) + "_" + strconv.FormatUint(iAdminPermissionId, 10)

		if oErr := oSelf.DeleteOne(oSelf.AdminRolesToAdminPermissionsIndex, sRelationId); oErr != nil {
			return oErr
		}
	}

	// 交集的部分不動，維持原樣

	return nil
}
