package outputApplicationMysqlLogic

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	oLogic := &AdminPermissionGroupLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func (oSelf *AdminPermissionGroupLogic) AddAdminPermissionGroup(oValue *domain.AdminPermissionGroupValue) error {
	oAdminPermissionGroupColumns, oErr := pkgUtility.StructToMap(oValue)
	if oErr != nil {
		return oErr
	}
	fmt.Println("oAdminPermissionGroupColumns=", oAdminPermissionGroupColumns)

	delete(oAdminPermissionGroupColumns, "admin_permissions")

	oError := oSelf.DB.WithContext(oSelf.Context).Transaction(func(oTx *gorm.DB) error {

		// 1. 先插 admin_permission_group
		oResult := oTx.Model(&domain.AdminPermissionGroup{}).Create(oAdminPermissionGroupColumns)

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
			return errors.New("0 rows inserted")
		}

		var iAdminPermissionGroupId uint64
		oResult = oTx.Raw("SELECT LAST_INSERT_ID()").Scan(&iAdminPermissionGroupId)

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

func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
	var aFlat []*domain.AdminPermissionGroup

	oResult := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.AdminPermissionGroup{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Order("id ASC").
		Find(&aFlat)

	if oResult.Error != nil {
		return nil, oResult.Error
	}

	aByParent := make(map[uint64][]*domain.AdminPermissionGroup, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.AdminPermissionGroup)
	fnAttach = func(oNode *domain.AdminPermissionGroup) {
		for _, oChild := range aByParent[oNode.Id] {
			fnAttach(oChild)
			oNode.Children = append(oNode.Children, *oChild)
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
		Preload("AdminPermissions").
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
		Preload("Parent").
		Preload("Children").
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

		oResult := oQuery.
			Limit(int(*oLimit.Count)).
			Offset(int(*oLimit.Offset)).
			Find(&aAdminPermissionGroups)
		oFindErr = oResult.Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.AdminPermissionGroup{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oResult := oQuery.Count(&iTotal)
		oCountErr = oResult.Error
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
