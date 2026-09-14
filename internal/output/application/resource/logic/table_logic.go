package outputApplicationResourceLogic

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
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

	oTable := pkgProtoToDomain.Table(oProtoTable)

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
		oGame := pkgProtoToDomain.Game(oOne.GetGame())
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
