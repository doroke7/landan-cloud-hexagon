package outputApplicationMongodbLogic

import (
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
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
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	var aGameTypes []*domain.GameType
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
		if oErr != nil {
			oFindErr = oErr
			return
		}
		defer oCursor.Close(oSelf.Context)

		oFindErr = oCursor.All(oSelf.Context, &aGameTypes)
	}()

	go func() {
		defer oWaitGroup.Done()

		iTotal, oCountErr = oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aGameTypes, 0, oFindErr
	}

	if oCountErr != nil {
		return aGameTypes, 0, oCountErr
	}

	return aGameTypes, uint64(iTotal), nil
}
