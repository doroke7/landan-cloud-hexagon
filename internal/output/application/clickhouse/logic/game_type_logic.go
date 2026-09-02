package outputApplicationClickhouseLogic

import (
	"strings"
	"sync"

	domain "example/internal/domain"
	clickhouseBase "example/internal/output/application/clickhouse"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type GameTypeLogic struct {
	*clickhouseBase.AbstractClickhouse
}

func NewGameTypeLogic(oAbstractLogic *clickhouseBase.AbstractClickhouse) outputPortAnyLogic.GameTypeLogic {
	return &GameTypeLogic{
		AbstractClickhouse: oAbstractLogic,
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

func (oSelf *GameTypeLogic) ShowGameTypesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, uint64, error) {
	aWheres := oSelf.FiltersToWheres(aFilters)
	aOrders := oSelf.SortersToOrders(aSorters)
	oLimit := oSelf.PaginationToLimit(oPagination)

	var aGameTypes []*domain.GameType
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
			Model(&domain.GameType{}).
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
			Find(&aGameTypes).Error
	}()

	go func() {
		defer oWaitGroup.Done()

		oQuery := oSelf.DB.WithContext(oSelf.Context).Model(&domain.GameType{}).Where("deleted_at = ?", "2038-01-19 03:14:07")
		for _, oWhere := range aWheres {
			oQuery = oQuery.Where(*oWhere.Field+" "+*oWhere.Operator+" ?", oWhere.Value)
		}

		oCountErr = oQuery.Count(&iTotal).Error
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
