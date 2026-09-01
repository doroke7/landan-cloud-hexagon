package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
)

type GameTypeLogic struct {
	*resourceBase.AbstractResource
}

func NewGameTypeLogic(oAbstractLogic *resourceBase.AbstractResource) outputPortAnyLogic.GameTypeLogic {
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
			Id:        uint(oOne.GetId()),
			ParentId:  uint(oOne.GetParentId()),
			Key:       oOne.GetKey(),
			Name:      oOne.GetName(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
	}

	aByParent := make(map[uint][]*domain.GameType, len(aFlat))
	for _, oOne := range aFlat {
		aByParent[oOne.ParentId] = append(aByParent[oOne.ParentId], oOne)
	}

	var fnAttach func(oNode *domain.GameType)
	fnAttach = func(oNode *domain.GameType) {
		oNode.Children = []domain.GameType{}

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
