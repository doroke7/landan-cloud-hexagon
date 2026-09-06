package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
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

// ShowTree 走 gRPC 拿回「平的」全部 admin_permission_group，再自己用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root（組 tree 每個 adapter 各寫一份）。
func (oSelf *AdminPermissionGroupLogic) ShowTree() ([]*domain.AdminPermissionGroup, error) {
	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowTree(oSelf.Context, &pbResourceLogic.AdminPermissionGroupShowTreeInput{})
	if oErr != nil {
		return nil, oErr
	}

	aFlat := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		aFlat = append(aFlat, &domain.AdminPermissionGroup{
			Id:        uint64(oOne.GetId()),
			ParentId:  uint64(oOne.GetParentId()),
			Key:       oOne.GetKey(),
			Name:      oOne.GetName(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
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

// ShowAdminPermissionGroupById 尚未接 gRPC，先回 nil。
func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupById(iId uint64) (*domain.AdminPermissionGroup, error) {
	return nil, nil
}

func (oSelf *AdminPermissionGroupLogic) ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, uint64, error) {

	oRequest := &pbResourceLogic.AdminPermissionGroupShowAdminPermissionGroupsTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.AdminPermissionGroup.ShowAdminPermissionGroupsTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, 0, oErr
	}

	aAdminPermissionGroups := make([]*domain.AdminPermissionGroup, 0, len(oResponse.GetAdminPermissionGroups()))
	for _, oOne := range oResponse.GetAdminPermissionGroups() {
		aAdminPermissionGroups = append(aAdminPermissionGroups, &domain.AdminPermissionGroup{
			Id:        uint64(oOne.GetId()),
			ParentId:  uint64(oOne.GetParentId()),
			Key:       oOne.GetKey(),
			Name:      oOne.GetName(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
	}

	iTotal := uint64(oResponse.GetTotal())
	return aAdminPermissionGroups, iTotal, nil
}
