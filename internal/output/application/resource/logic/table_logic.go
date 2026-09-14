package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type TableLogic struct {
	*AbstractLogic
}

func NewTableLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.TableLogic {
	oLogic := &TableLogic{
		AbstractLogic: oAbstractLogic,
	}

	return oLogic
}

func protoTableToDomainTable(oProtoTable *pbResource.Table) domain.Table {
	if oProtoTable == nil {
		return domain.Table{}
	}

	oGame := protoGameToDomainGame(oProtoTable.GetGame())
	oTable := domain.Table{
		Id:          uint64(oProtoTable.GetId()),
		No:          oProtoTable.GetNo(),
		GameId:      uint64(oProtoTable.GetGameId()),
		Key:         oProtoTable.GetKey(),
		State:       uint8(oProtoTable.GetState()),
		Description: oProtoTable.GetDescription(),
		Result:      oProtoTable.GetResult(),
		StartedAt:   oProtoTable.GetStartedAt().AsTime(),
		EndedAt:     oProtoTable.GetEndedAt().AsTime(),
		CreatedAt:   oProtoTable.GetCreatedAt().AsTime(),
		UpdatedAt:   oProtoTable.GetUpdatedAt().AsTime(),
		DeletedAt:   oProtoTable.GetDeletedAt().AsTime(),
		Game:        oGame,
	}

	return oTable
}

func (oSelf *TableLogic) ShowTableById(iId uint64) (*domain.Table, error) {

	oRequest := &pbResourceLogic.TableShowTableByIdInput{Id: iId}

	oResponse, oErr := oSelf.ResourceLogicClient.Table.ShowTableById(oSelf.Context, oRequest)
	if oErr != nil {
		return nil, oErr
	}

	oProtoTable := oResponse.GetTable()
	if oProtoTable.GetId() == 0 {
		return nil, nil
	}

	oTable := protoTableToDomainTable(oProtoTable)

	return &oTable, nil
}

func (oSelf *TableLogic) ShowTablesTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, uint64, error) {

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
		oGame := protoGameToDomainGame(oOne.GetGame())
		oTable := &domain.Table{
			Id:          uint64(oOne.GetId()),
			No:          oOne.GetNo(),
			GameId:      uint64(oOne.GetGameId()),
			Key:         oOne.GetKey(),
			State:       uint8(oOne.GetState()),
			Description: oOne.GetDescription(),
			Result:      oOne.GetResult(),
			StartedAt:   oOne.GetStartedAt().AsTime(),
			EndedAt:     oOne.GetEndedAt().AsTime(),
			CreatedAt:   oOne.GetCreatedAt().AsTime(),
			UpdatedAt:   oOne.GetUpdatedAt().AsTime(),
			DeletedAt:   oOne.GetDeletedAt().AsTime(),
			Game:        oGame,
		}
		aTables = append(aTables, oTable)
	}

	iTotal := uint64(oResponse.GetTotal())
	return aTables, iTotal, oErr
}

// ShowTablesByFiltersWithSortersPagination 直接複用 logic rpc（ShowTablesTotalBy…），把 total 丟掉。
func (oSelf *TableLogic) ShowTablesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.Table, error) {
	aTables, _, oErr := oSelf.ShowTablesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	return aTables, oErr
}
