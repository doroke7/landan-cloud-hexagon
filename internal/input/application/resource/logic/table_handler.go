package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResource "example/pb/resource"
	pbResourceLogic "example/pb/resource/logic"
	pkgInput "example/pkg/input"
)

type TableHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedTableLogicServer
	LogicTableUsecase usecasePortAnyLogic.TableUsecase
}

func NewTableHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oTableUsecase usecasePortAnyLogic.TableUsecase) *TableHandler {
	return &TableHandler{
		AbstractHandler:   oAbstractHandler,
		LogicTableUsecase: oTableUsecase,
	}
}

func domainTableToProtoTable(oTable *domain.Table) *pbResource.Table {
	if oTable == nil {
		return nil
	}

	return &pbResource.Table{
		Id:          uint64(oTable.Id),
		No:          oTable.No,
		GameId:      uint64(oTable.GameId),
		Key:         oTable.Key,
		State:       uint64(oTable.State),
		Description: oTable.Description,
		Result:      oTable.Result,
		StartedAt:   timestamppb.New(oTable.StartedAt),
		EndedAt:     timestamppb.New(oTable.EndedAt),
		CreatedAt:   timestamppb.New(oTable.CreatedAt),
		UpdatedAt:   timestamppb.New(oTable.UpdatedAt),
		DeletedAt:   timestamppb.New(oTable.DeletedAt),
		Game:        domainGameToProtoGame(&oTable.Game),
	}
}

func (oSelf *TableHandler) ShowTableById(oContext context.Context, oReq *pbResourceLogic.TableShowTableByIdInput) (*pbResourceLogic.TableShowTableByIdOutput, error) {

	oTable, oErr := oSelf.LogicTableUsecase.ShowTableById(oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	if oTable == nil {
		return nil, nil
	}

	return &pbResourceLogic.TableShowTableByIdOutput{
		Table: domainTableToProtoTable(oTable),
	}, nil
}

func (oSelf *TableHandler) ShowTablesTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.TableShowTablesTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.TableShowTablesTotalByFiltersWithSortersPaginationOutput, error) {

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aFilters := make([]*pkgInput.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aFilters = append(aFilters, &pkgInput.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aSorters := make([]*pkgInput.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		aSorters = append(aSorters, &pkgInput.Sorter{
			Field: &sField,
			Order: &sOrder,
		})
	}

	aTables, iTotal, oErr := oSelf.LogicTableUsecase.ShowTablesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	aPbTables := make([]*pbResource.Table, 0, len(aTables))
	for _, oTable := range aTables {
		aPbTables = append(aPbTables, &pbResource.Table{
			Id:          uint64(oTable.Id),
			No:          oTable.No,
			GameId:      uint64(oTable.GameId),
			Key:         oTable.Key,
			State:       uint64(oTable.State),
			Description: oTable.Description,
			Result:      oTable.Result,
			StartedAt:   timestamppb.New(oTable.StartedAt),
			EndedAt:     timestamppb.New(oTable.EndedAt),
			CreatedAt:   timestamppb.New(oTable.CreatedAt),
			UpdatedAt:   timestamppb.New(oTable.UpdatedAt),
			DeletedAt:   timestamppb.New(oTable.DeletedAt),
			Game:        domainGameToProtoGame(&oTable.Game),
		})
	}

	return &pbResourceLogic.TableShowTablesTotalByFiltersWithSortersPaginationOutput{
		Total:  iTotal,
		Tables: aPbTables,
	}, oErr
}
