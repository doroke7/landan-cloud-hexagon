package resource

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"

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

	oRequest := &pbResourceLogic.TableShowTablesTotalByFiltersWithSortersPaginationInput{}

	if oPagination != nil {
		oPbPagination := &pbResourceLogic.TablePagination{}
		if oPagination.Size != nil {
			oPbPagination.Size = uint64(*oPagination.Size)
		}
		if oPagination.Page != nil {
			oPbPagination.Page = uint64(*oPagination.Page)
		}
		oRequest.Pagination = oPbPagination
	}

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oFilter.Value)
		if oErr != nil {
			continue
		}

		oPbFilter := &pbResourceLogic.TableFilter{
			Field: *oFilter.Field,
			Value: oValue,
		}
		if oFilter.Operator != nil {
			oPbFilter.Operator = *oFilter.Operator
		}

		oRequest.Filters = append(oRequest.Filters, oPbFilter)
	}

	for _, oSorter := range aSorters {
		if oSorter == nil || oSorter.Field == nil || oSorter.Order == nil {
			continue
		}

		oRequest.Sorters = append(oRequest.Sorters, &pbResourceLogic.TableSorter{
			Field: *oSorter.Field,
			Order: *oSorter.Order,
		})
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
