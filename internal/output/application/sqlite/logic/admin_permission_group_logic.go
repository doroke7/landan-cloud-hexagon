package outputApplicationSqliteLogic

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputApplicationSqlite "example/internal/output/application/sqlite"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionGroupLogic struct {
	*outputApplicationSqlite.AbstractSqlite
}

func NewAdminPermissionGroupLogic(oAbstractSqlite *outputApplicationSqlite.AbstractSqlite) outputPortAnyLogic.AdminPermissionGroupLogic {
	return &AdminPermissionGroupLogic{
		AbstractSqlite: oAbstractSqlite,
	}
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupVariable) error {
	oAdminPermissionGroupColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	// admin_permissions 不是 admin_permission_groups 的欄位，拆出來另外插
	delete(oAdminPermissionGroupColumns, "admin_permissions")

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		// 1. 先插 admin_permission_group
		oResult := oTx.
			Model(&domain.AdminPermissionGroup{}).
			Create(oAdminPermissionGroupColumns)

		if oResult.Error != nil {
			if errors.Is(oResult.Error, gorm.ErrDuplicatedKey) {
				var sKey string
				if oValue.Key != nil {
					sKey = *oValue.Key
				}
				sError := fmt.Sprintf("admin_permission_group key=%s already exists", sKey)

				oDuplicateError := pkgUtility.NewDefaultError(sError, -2, 200)

				return oDuplicateError
			}

			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			oZeroRowsError := errors.New("0 rows inserted")

			return oZeroRowsError
		}

		var iAdminPermissionGroupId uint64
		oResult = oTx.Raw("SELECT last_insert_rowid()").Scan(&iAdminPermissionGroupId)

		if oResult.Error != nil {
			return oResult.Error
		}

		// 2. 沒帶 admin_permissions 就結束
		if len(oValue.AdminPermissions) == 0 {
			return nil
		}

		// 逐筆插 admin_permissions（不用批次，撞唯一索引時才能指出是哪個 type / key）
		for _, oAdminPermissionValue := range oValue.AdminPermissions {

			if oAdminPermissionValue == nil {
				continue
			}

			oAdminPermissionColumns, oErr := pkgUtility.StructToMap(oAdminPermissionValue)
			if oErr != nil {
				return oErr
			}

			delete(oAdminPermissionColumns, "id")
			oAdminPermissionColumns["admin_permission_group_id"] = iAdminPermissionGroupId

			oResult = oTx.Model(&domain.AdminPermission{}).Create(oAdminPermissionColumns)

			if oResult.Error != nil {
				if errors.Is(oResult.Error, gorm.ErrDuplicatedKey) {
					var iType uint8
					if oAdminPermissionValue.Type != nil {
						iType = *oAdminPermissionValue.Type
					}

					var sKey string
					if oAdminPermissionValue.Key != nil {
						sKey = *oAdminPermissionValue.Key
					}
					sError := fmt.Sprintf("admin_permission type=%d key=%s already exists", iType, sKey)

					oDuplicateError := pkgUtility.NewDefaultError(sError, -2, 200)

					return oDuplicateError
				}

				return oResult.Error
			}
		}

		return nil
	})

	return oError
}

// EditAdminPermissionGroupById 更新 group 本身欄位，並同步 admin_permissions：
// 傳進來沒帶 id 的就新增，帶 id 的就修改；db 裡面現有、但沒出現在傳進來 id 清單裡的就刪除。
func (oSelf *AdminPermissionGroupLogic) EditAdminPermissionGroupById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oAdminPermissionGroupColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}

	delete(oAdminPermissionGroupColumns, "admin_permissions")

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		// 1. 更新 admin_permission_group 本身欄位
		oResult := oTx.
			Model(&domain.AdminPermissionGroup{}).
			Where("id = ?", iId).
			UpdateColumns(oAdminPermissionGroupColumns)

		if oResult.Error != nil {
			return oResult.Error
		}

		if oResult.RowsAffected == 0 {
			oZeroRowsError := errors.New("0 rows updated")

			return oZeroRowsError
		}

		// 2. 撈出 db 現有未刪除的 admin_permissions，準備跟傳進來的 id 清單比對
		var aExistingAdminPermissions []*domain.AdminPermission
		oResult = oTx.
			Model(&domain.AdminPermission{}).
			Where("admin_permission_group_id = ?", iId).
			Where("deleted_at = ?", "2038-01-19 03:14:07").
			Find(&aExistingAdminPermissions)

		if oResult.Error != nil {
			return oResult.Error
		}

		// 3. 逐筆處理傳進來的 admin_permissions：沒帶 id 就新增，帶 id 就修改
		aKeptAdminPermissionIds := make(map[uint64]bool, len(oValue.AdminPermissions))
		for _, oAdminPermissionValue := range oValue.AdminPermissions {
			if oAdminPermissionValue == nil {
				continue
			}

			oAdminPermissionColumns, oErr := pkgUtility.StructToMap(oAdminPermissionValue)
			if oErr != nil {
				return oErr
			}

			delete(oAdminPermissionColumns, "id")

			if oAdminPermissionValue.Id == nil {
				oAdminPermissionColumns["admin_permission_group_id"] = iId

				oResult = oTx.Model(&domain.AdminPermission{}).Create(oAdminPermissionColumns)
			} else {
				iAdminPermissionId := *oAdminPermissionValue.Id
				aKeptAdminPermissionIds[iAdminPermissionId] = true

				oResult = oTx.
					Model(&domain.AdminPermission{}).
					Where("id = ?", iAdminPermissionId).
					UpdateColumns(oAdminPermissionColumns)
			}

			if oResult.Error != nil {
				if errors.Is(oResult.Error, gorm.ErrDuplicatedKey) {
					var iType uint8
					if oAdminPermissionValue.Type != nil {
						iType = *oAdminPermissionValue.Type
					}

					var sKey string
					if oAdminPermissionValue.Key != nil {
						sKey = *oAdminPermissionValue.Key
					}
					sError := fmt.Sprintf("admin_permission type=%d key=%s already exists", iType, sKey)

					oDuplicateError := pkgUtility.NewDefaultError(sError, -2, 200)

					return oDuplicateError
				}

				return oResult.Error
			}
		}

		// 4. db 裡面有、但沒出現在傳進來 id 清單裡的，軟刪除
		oNow := time.Now()
		for _, oExistingAdminPermission := range aExistingAdminPermissions {
			if aKeptAdminPermissionIds[oExistingAdminPermission.Id] {
				continue
			}

			oResult = oTx.
				Model(&domain.AdminPermission{}).
				Where("id = ?", oExistingAdminPermission.Id).
				UpdateColumn("deleted_at", oNow)

			if oResult.Error != nil {
				return oResult.Error
			}
		}

		return nil
	})

	return oError
}

// ShowTree 先把所有未刪除的 admin_permission_group 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
	var aFlat []*domain.AdminPermissionGroup

	oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Order("id ASC").
		Find(&aFlat).Error

	if oErr != nil {
		return nil, oErr
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

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&oAdminPermissionGroup).
		Preload("Parent").
		Preload("Children").
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		First(&oAdminPermissionGroup, iId)

	if oResult.Error != nil {
		if errors.Is(oResult.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, oResult.Error
	}

	return &oAdminPermissionGroup, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroups() ([]*domain.AdminPermissionGroup, error) {
	var aAdminPermissionGroups []*domain.AdminPermissionGroup

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Order("id ASC").
		Find(&aAdminPermissionGroups)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)
	aOrders := oSelf.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.
			DB.
			WithContext(oSelf.Context).
			Preload("Parent").
			Preload("Children").
			Model(&domain.AdminPermissionGroup{}).
			Where("deleted_at = ?", "2038-01-19 03:14:07")

		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		for _, oOrder := range aOrders {
			if oOrder == nil || oOrder.Field == nil {
				continue
			}

			sDirection := "ASC"
			if oOrder.Value != nil && strings.EqualFold(*oOrder.Value, "desc") {
				sDirection = "DESC"
			}

			oQuery = oQuery.Order(*oOrder.Field + " " + sDirection)
		}

		oFindErr = oQuery.
			Limit(int(*oLimit.Count)).
			Offset(int(*oLimit.Offset)).
			Find(&aAdminPermissionGroups).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminPermissionGroup{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aAdminPermissionGroups, 0, oFindErr
	}

	if oCountErr != nil {
		return aAdminPermissionGroups, 0, oCountErr
	}

	return aAdminPermissionGroups, uint64(iTotal), nil
}
