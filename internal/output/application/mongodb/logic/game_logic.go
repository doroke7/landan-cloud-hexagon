package outputApplicationMongodbLogic

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

// oDeletedAtZero 跟 mongodb/model adapter 用同一個「未刪除」標記值，
// 讓軟刪除的語意在各個 adapter 間保持一致。
var oDeletedAtZero = time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC)

type GameLogic struct {
	*AbstractLogic
	Collection *mongo.Collection
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	oCollection := oAbstractLogic.Database.Collection("games")

	return &GameLogic{
		AbstractLogic: oAbstractLogic,
		Collection:    oCollection,
	}
}

func (oSelf *GameLogic) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	var aGames []*domain.Game
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

		oFindErr = oCursor.All(oSelf.Context, &aGames)
	}()

	go func() {
		defer oWaitGroup.Done()

		iTotal, oCountErr = oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aGames, 0, oFindErr
	}

	if oCountErr != nil {
		return aGames, 0, oCountErr
	}

	return aGames, uint64(iTotal), nil
}

func (oSelf *GameLogic) ShowGameById(iId uint64) (*domain.Game, error) {
	var oGame domain.Game

	oResult := oSelf.Collection.FindOne(oSelf.Context, bson.M{
		"_id":        iId,
		"deleted_at": oDeletedAtZero,
	})
	oErr := oResult.Decode(&oGame)

	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, oErr
	}

	return &oGame, nil
}

func (oSelf *GameLogic) ShowGamesByGameTypeId(iGameTypeId uint64) ([]*domain.Game, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{
		"game_type_id": iGameTypeId,
		"deleted_at":   oDeletedAtZero,
	})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aGames []*domain.Game
	if oErr := oCursor.All(oSelf.Context, &aGames); oErr != nil {
		return nil, oErr
	}

	return aGames, nil
}
