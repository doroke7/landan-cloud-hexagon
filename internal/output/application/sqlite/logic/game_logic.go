package outputApplicationSqliteLogic

import (
	"strings"
	"sync"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type GameLogic struct {
	*AbstractLogic
}

func NewGameLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.GameLogic {
	return &GameLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *GameLogic) ShowGamesByGameTypeId(iGameTypeId uint64) ([]*domain.Game, error) {
	var aGames []*domain.Game

	if oErr := oSelf.DB.WithContext(oSelf.Context).
		Preload("GameType").
		Preload("GameType.Parent").
		Model(&domain.Game{}).
		Where("game_type_id = ?", iGameTypeId).
		Where("deleted_at = ?", "2038-01-19 03:14:07").
		Find(&aGames).Error; oErr != nil {
		return nil, oErr
	}

	return aGames, nil
}

func (oSelf *GameLogic) ShowGamesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Game, uint64, error) {
	aWheres := oSelf.AbstractSqlite.FiltersToWheres(aFilters)
	aOrders := oSelf.AbstractSqlite.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aGames []*domain.Game
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
			Preload("GameType").
			Preload("GameType.Parent").
			Model(&domain.Game{}).
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
			Find(&aGames).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.Game{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
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
