package resource

import (
	"encoding/json"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type TableLogic struct {
	*AbstractLogic
}

func NewTableLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.TableLogic {
	return &TableLogic{
		AbstractLogic: oAbstractLogic,
	}
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.Table, int64, error) {

	oRequest := &pbResourceLogic.TableShowTablesTotalByFiltersWithSortersPaginationInput{
		Filters: oSelf.ToFilters(aFilters),
		Sorters: oSelf.ToSorters(aSorters),
	}

	if oPagination != nil {
		oRequest.Pagination = oSelf.ToPagination(oPagination)
	}

	oResponse, oErr := oSelf.ResourceLogicClient.Table.ShowTablesTotalByFiltersWithSortersPagination(oSelf.Context, oRequest)

	aTables := make([]*domain.Table, 0, len(oResponse.GetTables()))
	for _, oOne := range oResponse.GetTables() {
		aTables = append(aTables, &domain.Table{
			Id:          uint(oOne.GetId()),
			No:          oOne.GetNo(),
			GameId:      uint(oOne.GetGameId()),
			Key:         oOne.GetKey(),
			State:       uint8(oOne.GetState()),
			Description: oOne.GetDescription(),
			Result:      json.RawMessage(oOne.GetResult()),
			StartedAt:   oOne.GetStartedAt().AsTime(),
			EndedAt:     oOne.GetEndedAt().AsTime(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			Game:        gameFromPb(oOne.GetGame()),
		})
	}

	iTotal := oResponse.GetTotal()
	return aTables, int64(iTotal), oErr
}
