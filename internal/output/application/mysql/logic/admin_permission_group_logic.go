package outputApplicationMysqlLogic

import (
	"strings"
	"sync"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminPermissionGroupLogic struct {
	*AbstractLogic
}

func NewAdminPermissionGroupLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminPermissionGroupLogic {
	return &AdminPermissionGroupLogic{
		AbstractLogic: oAbstractLogic,
	}
}

// ShowTree 先把所有未刪除的 admin_permission_group 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
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
