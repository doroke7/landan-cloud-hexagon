package resource

import (
	"google.golang.org/protobuf/types/known/structpb"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkg "example/pkg"
)

type GameTypeModel struct {
	*AbstractModel
}

func NewGameTypeModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) AddOne(oGameTypeParm *domain.GameTypeValue) (bool, error) {

	oRequest := &pbResourceModel.GameTypeAddOneInput{}

	oRequest.Key = oGameTypeParm.Key
	oRequest.Name = oGameTypeParm.Name

	_, oErr := oSelf.ResourceModelClient.GameType.AddOne(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return true, nil
}

func (oSelf *GameTypeModel) ShowOneById(iId uint) (*domain.GameType, error) {

	oResponse, oErr := oSelf.ResourceModelClient.GameType.ShowOneById(
		oSelf.Context,
		&pbResourceModel.GameTypeShowOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return nil, oErr
	}

	if oResponse.GetId() == 0 {
		return nil, nil
	}

	return &domain.GameType{
		Id:        uint(oResponse.GetId()),
		Key:       oResponse.GetKey(),
		Name:      oResponse.GetName(),
		CreatedAt: oResponse.GetCreatedAt().AsTime(),
		UpdatedAt: oResponse.GetUpdatedAt().AsTime(),
		DeletedAt: oResponse.GetDeletedAt().AsTime(),
	}, nil
}

func (oSelf *GameTypeModel) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.GameType, error) {

	oRequest := &pbResourceModel.GameTypeShowOnesByWheresWithOrdersLimitInput{
		Limit: &pbResourceModel.GameTypeLimit{
			Offset: uint64(*oLimit.Offset),
			Count:  uint64(*oLimit.Count),
		},
	}

	for _, oWhere := range aWheres {
		if oWhere == nil || oWhere.Field == nil || oWhere.Operator == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oWhere.Value)
		if oErr != nil {
			continue
		}

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceModel.GameTypeWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	for _, oOrder := range aOrders {
		if oOrder == nil || oOrder.Field == nil || oOrder.Value == nil {
			continue
		}

		oRequest.Orders = append(oRequest.Orders, &pbResourceModel.GameTypeOrder{
			Field: *oOrder.Field,
			Value: *oOrder.Value,
		})
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.ShowOnesByWheresWithOrdersLimit(oSelf.Context, oRequest)

	if oErr != nil {
		return nil, oErr
	}

	aGameTypes := make([]*domain.GameType, 0, len(oResponse.GetGameTypes()))
	for _, oOne := range oResponse.GetGameTypes() {
		aGameTypes = append(aGameTypes, &domain.GameType{
			Id:        uint(oOne.GetId()),
			Key:       oOne.GetKey(),
			Name:      oOne.GetName(),
			CreatedAt: oOne.GetCreatedAt().AsTime(),
			UpdatedAt: oOne.GetUpdatedAt().AsTime(),
			DeletedAt: oOne.GetDeletedAt().AsTime(),
		})
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error) {
	aWheres := pkg.FiltersToMysqlWheres([]string{}, aFilters)
	aOrders := pkg.SortersToMysqlOrders([]string{}, aSorters)
	oLimit := pkg.PaginationToMysqlLimit(oPagination)

	return oSelf.ShowOnesByWheresWithOrdersLimit(aWheres, aOrders, oLimit)
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	aWheres := pkg.FiltersToMysqlWheres([]string{}, aFilters)

	return oSelf.TotalByWheres(aWheres)
}

func (oSelf *GameTypeModel) EditOneById(oGameType *domain.GameTypeValue, iId uint) (bool, error) {

	oRequest := &pbResourceModel.GameTypeEditOneByIdInput{Id: uint32(iId)}

	oRequest.Key = oGameType.Key
	oRequest.Name = oGameType.Name

	oResponse, oErr := oSelf.ResourceModelClient.GameType.EditOneById(oSelf.Context, oRequest)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint) (bool, error) {

	oResponse, oErr := oSelf.ResourceModelClient.GameType.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.GameTypeRemoveOneByIdInput{Id: uint32(iId)},
	)

	if oErr != nil {
		return false, oErr
	}

	return oResponse.GetStatus(), nil
}

func (oSelf *GameTypeModel) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {

	oRequest := &pbResourceModel.GameTypeTotalByWheresInput{}

	for _, oWhere := range aWheres {
		if oWhere == nil || oWhere.Field == nil || oWhere.Operator == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oWhere.Value)
		if oErr != nil {
			continue
		}

		oRequest.Wheres = append(oRequest.Wheres, &pbResourceModel.GameTypeWhere{
			Field:    *oWhere.Field,
			Operator: *oWhere.Operator,
			Value:    oValue,
		})
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.TotalByWheres(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
}
