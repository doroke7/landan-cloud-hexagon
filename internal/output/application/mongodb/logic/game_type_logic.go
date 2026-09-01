package outputApplicationMongodbLogic

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type GameTypeLogic struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
}

func NewGameTypeLogic(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) outputPortAnyLogic.GameTypeLogic {
	return &GameTypeLogic{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oAbstractMongodb.Database.Collection("game_types"),
	}
}

// ShowTree 先把所有未刪除的 game_type 一次撈成平的，再用 ParentId 掛 Children，
// 回傳 ParentId == 0 的 root。
func (oSelf *GameTypeLogic) ShowTree() ([]*domain.GameType, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{"deleted_at": oDeletedAtZero})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aFlat []*domain.GameType
	if oErr := oCursor.All(oSelf.Context, &aFlat); oErr != nil {
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
