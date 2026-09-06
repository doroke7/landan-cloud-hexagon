package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputApplicationResource "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type GameTypeLogic struct {
	*outputApplicationResource.AbstractResource
}

func NewGameTypeLogic(oAbstractLogic *outputApplicationResource.AbstractResource) outputPortAnyLogic.GameTypeLogic {
	return &GameTypeLogic{
		AbstractResource: oAbstractLogic,
	}
}

// ShowTree 走 gRPC 拿回「平的」全部 game_type，再自己用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root（組 tree 每個 adapter 各寫一份）。
func (oSelf *GameTypeLogic) ShowTree() ([]*domain.GameType, error) {
	oResponse, oErr := oSelf.ResourceLogicClient.GameType.ShowTree(oSelf.Context, &pbResourceLogic.GameTypeShowTreeInput{})
	if oErr != nil {
		return nil, oErr
	}

	aFlat := make([]*domain.GameType, 0, len(oResponse.GetGameTypes()))
	for _, oOne := range oResponse.GetGameTypes() {
		aFlat = append(aFlat, &domain.GameType{
			Id:        uint64(oOne.GetId()),
			ParentId:  uint64(oOne.GetParentId()),
			Key:       oOne.GetKey(),
			Name:      oOne.GetName(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
	}

	aByParent := make(map[uint64][]*domain.GameType, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.GameType)
	fnAttach = func(oNode *domain.GameType) {
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

func (oSelf *GameTypeLogic) ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {

	oRequest := &pbResourceLogic.GameTypeShowGameTypesTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.GameType.ShowGameTypesTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, 0, oErr
	}

	aGameTypes := make([]*domain.GameType, 0, len(oResponse.GetGameTypes()))
	for _, oOne := range oResponse.GetGameTypes() {
		aGameTypes = append(aGameTypes, &domain.GameType{
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
	return aGameTypes, iTotal, nil
}
