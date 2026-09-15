package outputApplicationMysqlLogic

import (
	"errors"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgUtility "example/pkg/utility"

	"gorm.io/gorm"
)

type AdminRoleLogic struct {
	*AbstractLogic
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	oLogic := &AdminRoleLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	oColumns, oErr := pkgUtility.StructToMap(oVariable)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_permission_ids")

	if len(oColumns) == 0 {
		oZeroColumnsError := errors.New("0 columns to insert")

		return oZeroColumnsError
	}

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		oResult := oTx.
			Model(&domain.AdminRole{}).
			Create(oColumns)

		if oResult.Error != nil {
			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			oZeroRowsError := errors.New("0 rows inserted")

			return oZeroRowsError
		}

		var iAdminRoleId uint64
		oResult = oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminRoleId)

		if oResult.Error != nil {
			return oResult.Error
		}

		// 避免髒數據（例如 id 被重用）導致異常，插入前先把這個 admin_role_id 底下的關聯清乾淨
		oResult = oTx.Where("admin_role_id = ?", iAdminRoleId).Delete(&domain.AdminRolesToAdminPermission{})

		if oResult.Error != nil {
			return oResult.Error
		}

		if oVariable.AdminPermissionIds == nil || len(*oVariable.AdminPermissionIds) == 0 {
			return nil
		}

		aAdminRolesToAdminPermissions := make([]domain.AdminRolesToAdminPermission, 0, len(*oVariable.AdminPermissionIds))
		for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
			oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
				AdminRoleId:       iAdminRoleId,
				AdminPermissionId: iAdminPermissionId,
			}
			aAdminRolesToAdminPermissions = append(aAdminRolesToAdminPermissions, oAdminRolesToAdminPermission)
		}

		oResult = oTx.Create(&aAdminRolesToAdminPermissions)

		return oResult.Error
	})

	return oError
}

func (oSelf *AdminRoleLogic) EditAdminRoleById(oVariable *domain.AdminRoleVariable, iId uint64) error {
	oColumns, oErr := pkgUtility.StructToMap(oVariable)
	if oErr != nil {
		return oErr
	}

	delete(oColumns, "admin_permission_ids")

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		if len(oColumns) > 0 {
			oResult := oTx.
				Model(&domain.AdminRole{}).
				Where("id = ?", iId).
				UpdateColumns(oColumns)

			if oResult.Error != nil {
				return oResult.Error
			}

			if oResult.RowsAffected == 0 {
				oZeroRowsError := errors.New("0 rows updated")

				return oZeroRowsError
			}
		}

		if oVariable.AdminPermissionIds == nil {
			return nil
		}

		var aExistingAdminRolesToAdminPermissions []domain.AdminRolesToAdminPermission
		oResult := oTx.Where("admin_role_id = ?", iId).Find(&aExistingAdminRolesToAdminPermissions)

		if oResult.Error != nil {
			return oResult.Error
		}

		/*
		   AdminRoleToAdminPermission 更新策略
		   1. 取出 輸入的 AdminPermissionIds + 數據庫裡面 AdminPermissionIds，並轉成 Map
		   2. for-loop 輸入的 AdminPermissionIds ，如果 數據庫存在 -> continue, 如果不存在 數據寫入
		   3. for-loop 數據庫的 AdminPermissionIds ，如果 輸入存在 -> continue, 如果不存在 數據刪除

		*/

		oMapExistingAdminPermissionIds := make(map[uint64]bool, len(aExistingAdminRolesToAdminPermissions))
		for _, oExisting := range aExistingAdminRolesToAdminPermissions {
			oMapExistingAdminPermissionIds[oExisting.AdminPermissionId] = true
		}

		oInputtingAdminPermissionIds := make(map[uint64]bool, len(*oVariable.AdminPermissionIds))
		for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
			oInputtingAdminPermissionIds[iAdminPermissionId] = true
		}

		// a. 傳進來的 id 在 DB 不存在 -> 插入
		aAdminPermissionIdsToInsert := make([]domain.AdminRolesToAdminPermission, 0, len(oInputtingAdminPermissionIds))
		for iAdminPermissionId := range oInputtingAdminPermissionIds {
			if oMapExistingAdminPermissionIds[iAdminPermissionId] {
				continue
			}

			oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
				AdminRoleId:       iId,
				AdminPermissionId: iAdminPermissionId,
			}
			aAdminPermissionIdsToInsert = append(aAdminPermissionIdsToInsert, oAdminRolesToAdminPermission)
		}

		// b. DB 存在但不在傳進來的 id 裡面 -> 刪除
		aAdminPermissionIdsToDelete := make([]uint64, 0, len(oMapExistingAdminPermissionIds))
		for iAdminPermissionId := range oMapExistingAdminPermissionIds {
			if oInputtingAdminPermissionIds[iAdminPermissionId] {
				continue
			}

			aAdminPermissionIdsToDelete = append(aAdminPermissionIdsToDelete, iAdminPermissionId)
		}

		// c. 交集的部分不動，維持原樣

		if len(aAdminPermissionIdsToDelete) > 0 {
			oResult = oTx.
				Where("admin_role_id = ?", iId).
				Where("admin_permission_id IN ?", aAdminPermissionIdsToDelete).
				Delete(&domain.AdminRolesToAdminPermission{})

			if oResult.Error != nil {
				return oResult.Error
			}
		}

		if len(aAdminPermissionIdsToInsert) > 0 {
			oResult = oTx.Create(&aAdminPermissionIdsToInsert)

			if oResult.Error != nil {
				return oResult.Error
			}
		}

		return nil
	})

	return oError
}
