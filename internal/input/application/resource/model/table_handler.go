package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "example/pb"
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

// protoTableValueToDomainTableValue 把 gRPC 帶進來的 variable 攤成 domain.TableVariable（指標欄位可選）。
func protoTableValueToDomainTableValue(oVariable *pb.TableVariable) domain.TableVariable {
	var oValue domain.TableVariable
	if oVariable == nil {
		return oValue
	}

	oValue.No = oVariable.No
	oValue.Key = oVariable.Key
	oValue.Description = oVariable.Description

	if oVariable.GameId != nil {
		iGameId := uint64(*oVariable.GameId)
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

func (oSelf *TableHandler) AddOne(oContext context.Context, oReq *pbResourceModel.TableAddOneInput) (*pbResourceModel.TableAddOneOutput, error) {

	oTableValue := protoTableValueToDomainTableValue(oReq.GetVariable())

	oErr := oSelf.ModelTableUsecase.AddOne(&oTableValue)

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableAddOneOutput{}, nil
}

func (oSelf *TableHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.TableEditOneByIdInput) (*pbResourceModel.TableEditOneByIdOutput, error) {

	oTableValue := protoTableValueToDomainTableValue(oReq.GetVariable())

	oErr := oSelf.ModelTableUsecase.EditOneById(&oTableValue, uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableEditOneByIdOutput{}, nil
}

func (oSelf *TableHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.TableRemoveOneByIdInput) (*pbResourceModel.TableRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelTableUsecase.RemoveOneById(uint64(oReq.Id))

	if oErr != nil {
		return nil, status.Error(codes.Aborted, oErr.Error())
	}

	return &pbResourceModel.TableRemoveOneByIdOutput{}, nil
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
