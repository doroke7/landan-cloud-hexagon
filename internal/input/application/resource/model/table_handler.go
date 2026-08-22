package input_application_resource_model

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkg "example/pkg"
)

type TableHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedTableModelServer
	usecasePortAnyModel.TableUsecase
}

func NewTableHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oTableUsecase usecasePortAnyModel.TableUsecase) *TableHandler {
	return &TableHandler{
		AbstractHandler: oAbstractHandler,
		TableUsecase:    oTableUsecase,
	}
}

func (oSelf *TableHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.TableShowOneByIdInput) (*pbResourceModel.TableShowOneByIdOutput, error) {

	oTable, oErr := oSelf.TableUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	if oTable == nil {
		return nil, nil
	}

	return &pbResourceModel.TableShowOneByIdOutput{
		Id:          uint32(oTable.Id),
		No:          oTable.No,
		GameId:      uint32(oTable.GameId),
		Key:         oTable.Key,
		State:       uint32(oTable.State),
		Description: oTable.Description,
		Result:      string(oTable.Result),
		StartedAt:   timestamppb.New(oTable.StartedAt),
		EndedAt:     timestamppb.New(oTable.EndedAt),
		CreatedAt:   timestamppb.New(oTable.CreatedAt),
		UpdatedAt:   timestamppb.New(oTable.UpdatedAt),
		DeletedAt:   timestamppb.New(oTable.DeletedAt),
		Game:        gameToPb(&oTable.Game),
	}, nil
}

func (oSelf *TableHandler) AddOne(oContext context.Context, oReq *pbResourceModel.TableAddOneInput) (*pbResourceModel.TableAddOneOutput, error) {

	var oTableValue domain.TableValue

	oTableValue.No = oReq.No
	oTableValue.Key = oReq.Key
	oTableValue.Description = oReq.Description

	if oReq.GameId != nil {
		iGameId := uint(*oReq.GameId)
		oTableValue.GameId = &iGameId
	}

	if oReq.State != nil {
		iState := uint8(*oReq.State)
		oTableValue.State = &iState
	}

	if oReq.Result != nil {
		oRawMessage := json.RawMessage(*oReq.Result)
		oTableValue.Result = &oRawMessage
	}

	if oReq.StartedAt != nil {
		oStartedAt := oReq.StartedAt.AsTime()
		oTableValue.StartedAt = &oStartedAt
	}

	if oReq.EndedAt != nil {
		oEndedAt := oReq.EndedAt.AsTime()
		oTableValue.EndedAt = &oEndedAt
	}

	bResult, oErr := oSelf.TableUsecase.AddOne(&oTableValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableAddOneOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.TableEditOneByIdInput) (*pbResourceModel.TableEditOneByIdOutput, error) {

	var oTableValue domain.TableValue

	oTableValue.No = oReq.No
	oTableValue.Key = oReq.Key
	oTableValue.Description = oReq.Description

	if oReq.GameId != nil {
		iGameId := uint(*oReq.GameId)
		oTableValue.GameId = &iGameId
	}

	if oReq.State != nil {
		iState := uint8(*oReq.State)
		oTableValue.State = &iState
	}

	if oReq.Result != nil {
		oRawMessage := json.RawMessage(*oReq.Result)
		oTableValue.Result = &oRawMessage
	}

	if oReq.StartedAt != nil {
		oStartedAt := oReq.StartedAt.AsTime()
		oTableValue.StartedAt = &oStartedAt
	}

	if oReq.EndedAt != nil {
		oEndedAt := oReq.EndedAt.AsTime()
		oTableValue.EndedAt = &oEndedAt
	}

	bResult, oErr := oSelf.TableUsecase.EditOneById(&oTableValue, uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableEditOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.TableRemoveOneByIdInput) (*pbResourceModel.TableRemoveOneByIdOutput, error) {

	bResult, oErr := oSelf.TableUsecase.RemoveOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableRemoveOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.TableShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.TableShowOnesByFiltersWithSortersPaginationOutput, error) {

	aFilters := make([]*pkg.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aFilters = append(aFilters, &pkg.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aSorters := make([]*pkg.Sorter, 0, len(oReq.GetSorters()))
	for _, oOne := range oReq.GetSorters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		aSorters = append(aSorters, &pkg.Sorter{
			Field: &sField,
			Order: &sOrder,
		})
	}

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkg.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aTables, oErr := oSelf.TableUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aPbTables := make([]*pbResourceModel.TableShowOneByIdOutput, 0, len(aTables))
	for _, oTable := range aTables {
		aPbTables = append(aPbTables, &pbResourceModel.TableShowOneByIdOutput{
			Id:          uint32(oTable.Id),
			No:          oTable.No,
			GameId:      uint32(oTable.GameId),
			Key:         oTable.Key,
			State:       uint32(oTable.State),
			Description: oTable.Description,
			Result:      string(oTable.Result),
			StartedAt:   timestamppb.New(oTable.StartedAt),
			EndedAt:     timestamppb.New(oTable.EndedAt),
			CreatedAt:   timestamppb.New(oTable.CreatedAt),
			UpdatedAt:   timestamppb.New(oTable.UpdatedAt),
			DeletedAt:   timestamppb.New(oTable.DeletedAt),
			Game:        gameToPb(&oTable.Game),
		})
	}

	return &pbResourceModel.TableShowOnesByFiltersWithSortersPaginationOutput{
		Tables: aPbTables,
	}, nil
}

func (oSelf *TableHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.TableTotalByFiltersInput) (*pbResourceModel.TableTotalByFiltersOutput, error) {

	aFilters := make([]*pkg.Filter, 0, len(oReq.GetFilters()))
	for _, oOne := range oReq.GetFilters() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aFilters = append(aFilters, &pkg.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	iTotal, oErr := oSelf.TableUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.TableTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
