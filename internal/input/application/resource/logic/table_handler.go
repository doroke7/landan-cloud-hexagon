package input_application_resource_logic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
	pkg "example/pkg"
)

type TableHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedTableLogicServer
	usecasePortAnyLogic.TableUsecase
}

func NewTableHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oTableUsecase usecasePortAnyLogic.TableUsecase) *TableHandler {
	return &TableHandler{
		AbstractHandler: oAbstractHandler,
		TableUsecase:    oTableUsecase,
	}
}

func (oSelf *TableHandler) ShowTablesTotalByWheresWithOrdersLimit(oContext context.Context, oReq *pbResourceLogic.TableShowTablesTotalByWheresWithOrdersLimitInput) (*pbResourceLogic.TableShowTablesTotalByWheresWithOrdersLimitOutput, error) {

	iOffset := uint(oReq.GetLimit().GetOffset())
	iCount := uint(oReq.GetLimit().GetCount())
	oLimit := &pkg.Limit{
		Offset: &iOffset,
		Count:  &iCount,
	}

	aWheres := make([]*pkg.Where, 0, len(oReq.GetWheres()))
	for _, oOne := range oReq.GetWheres() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		aWheres = append(aWheres, &pkg.Where{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oOne.GetValue().AsInterface(),
		})
	}

	aOrders := make([]*pkg.Order, 0, len(oReq.GetOrders()))
	for _, oOne := range oReq.GetOrders() {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sValue := oOne.GetValue()
		aOrders = append(aOrders, &pkg.Order{
			Field: &sField,
			Value: &sValue,
		})
	}

	aTables, iTotal, oErr := oSelf.TableUsecase.ShowTablesTotalByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	aPbTables := make([]*pbResourceLogic.Table, 0, len(aTables))
	for _, oTable := range aTables {
		aPbTables = append(aPbTables, &pbResourceLogic.Table{
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

	return &pbResourceLogic.TableShowTablesTotalByWheresWithOrdersLimitOutput{
		Total:  uint64(iTotal),
		Tables: aPbTables,
	}, oErr
}
