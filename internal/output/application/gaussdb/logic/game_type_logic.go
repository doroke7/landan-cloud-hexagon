package outputApplicationGaussdbLogic

import (
	domain "example/internal/domain"
	gaussdbBase "example/internal/output/application/gaussdb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type GameTypeLogic struct {
	*gaussdbBase.AbstractGaussdb
}

func NewGameTypeLogic(oAbstractLogic *gaussdbBase.AbstractGaussdb) outputPortAnyLogic.GameType {
	return &GameTypeLogic{
		AbstractGaussdb: oAbstractLogic,
	}
}

// ShowTree 先把所有未刪除的 game_type 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *GameTypeLogic) ShowTree() ([]*domain.GameType, error) {
	var aFlat []*domain.GameType

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Model(&domain.GameType{}).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Order("id ASC").
		Find(&aFlat).Error; oErr != nil {
		return nil, oErr
	}

	aByParent := make(map[uint][]*domain.GameType, len(aFlat))
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
