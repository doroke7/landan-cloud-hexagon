package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type TableHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedTableModelServer
	ModelTableUsecase usecasePortAnyModel.TableUsecase
}

func NewTableHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oTableUsecase usecasePortAnyModel.TableUsecase) *TableHandler {
	return &TableHandler{
		AbstractHandler:   oAbstractHandler,
		ModelTableUsecase: oTableUsecase,
	}
}

func domainTableToProtoTable(oTable *domain.Table) *pbResource.Table {
	if oTable == nil {
		return nil
	}

	return &pbResource.Table{
		Id:          uint32(oTable.Id),
		No:          oTable.No,
		GameId:      uint32(oTable.GameId),
		Key:         oTable.Key,
		State:       uint32(oTable.State),
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

// protoTableVariableToDomainTableValue 把 gRPC 帶進來的 variable 攤成 domain.TableValue（指標欄位可選）。
func protoTableVariableToDomainTableValue(oVariable *pbResourceModel.TableVariable) domain.TableValue {
	var oValue domain.TableValue
	if oVariable == nil {
		return oValue
	}

	oValue.No = oVariable.No
	oValue.Key = oVariable.Key
	oValue.Description = oVariable.Description

	if oVariable.GameId != nil {
		iGameId := uint(*oVariable.GameId)
		oValue.GameId = &iGameId
	}

	if oVariable.State != nil {
		iState := uint8(*oVariable.State)
		oValue.State = &iState
	}

	if oVariable.Result != nil {
		sResult := *oVariable.Result
		oValue.Result = &sResult
	}

	if oVariable.StartedAt != nil {
		oStartedAt := oVariable.StartedAt.AsTime()
		oValue.StartedAt = &oStartedAt
	}

	if oVariable.EndedAt != nil {
		oEndedAt := oVariable.EndedAt.AsTime()
		oValue.EndedAt = &oEndedAt
	}

	return oValue
}

func (oSelf *TableHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.TableShowOneByIdInput) (*pbResourceModel.TableShowOneByIdOutput, error) {

	oTable, oErr := oSelf.ModelTableUsecase.ShowOneById(uint(oReq.Id))
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	if oTable == nil {
		return nil, nil
	}

	return &pbResourceModel.TableShowOneByIdOutput{
		Table: domainTableToProtoTable(oTable),
	}, nil
}

func (oSelf *TableHandler) AddOne(oContext context.Context, oReq *pbResourceModel.TableAddOneInput) (*pbResourceModel.TableAddOneOutput, error) {

	oTableValue := protoTableVariableToDomainTableValue(oReq.GetVariable())

	bResult, oErr := oSelf.ModelTableUsecase.AddOne(&oTableValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableAddOneOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.TableEditOneByIdInput) (*pbResourceModel.TableEditOneByIdOutput, error) {

	oTableValue := protoTableVariableToDomainTableValue(oReq.GetVariable())

	bResult, oErr := oSelf.ModelTableUsecase.EditOneById(&oTableValue, uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableEditOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.TableRemoveOneByIdInput) (*pbResourceModel.TableRemoveOneByIdOutput, error) {

	bResult, oErr := oSelf.ModelTableUsecase.RemoveOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableRemoveOneByIdOutput{
		Status: bResult,
	}, nil
}

func (oSelf *TableHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.TableShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.TableShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	iSize := uint(oReq.GetPagination().GetSize())
	iPage := uint(oReq.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aTables, oErr := oSelf.ModelTableUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	aProtoTables := make([]*pbResource.Table, 0, len(aTables))
	for _, oTable := range aTables {
		aProtoTables = append(aProtoTables, domainTableToProtoTable(oTable))
	}

	return &pbResourceModel.TableShowOnesByFiltersWithSortersPaginationOutput{
		Tables: aProtoTables,
	}, nil
}

func (oSelf *TableHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.TableTotalByFiltersInput) (*pbResourceModel.TableTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelTableUsecase.TotalByFilters(aFilters)
	if oErr != nil {
		return nil, status.Error(codes.NotFound, oErr.Error())
	}

	return &pbResourceModel.TableTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
