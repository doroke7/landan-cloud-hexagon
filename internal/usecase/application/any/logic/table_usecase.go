package usecase

import (
	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pkg "example/pkg"
)

type TableUsecase struct {
	*AbstractUsecase
	outputPortAnyLogic.TableLogic
}

func NewTableUsecase(oAbstractUsecase *AbstractUsecase, oTableLogic outputPortAnyLogic.TableLogic) usecasePortAnyLogic.TableUsecase {
	return &TableUsecase{
		AbstractUsecase: oAbstractUsecase,
		TableLogic:      oTableLogic,
	}
}

func (oSelf *TableUsecase) ShowTablesTotalByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.Table, int64, error) {

	aTables, iTotal, oErr := oSelf.TableLogic.ShowTablesTotalByWheresWithOrdersLimit(aWheres, aOrders, oLimit)

	return aTables, iTotal, oErr
}
