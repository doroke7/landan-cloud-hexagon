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

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.GameType, error) {

	oRequest := &pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationInput{}

	if oPagination != nil {
		oPbPagination := &pbResourceModel.GameTypePagination{}
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

		oPbFilter := &pbResourceModel.GameTypeFilter{
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

		oRequest.Sorters = append(oRequest.Sorters, &pbResourceModel.GameTypeSorter{
			Field: *oSorter.Field,
			Order: *oSorter.Order,
		})
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.ShowOnesByFiltersWithSortersPagination(oSelf.Context, oRequest)

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

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTypeTotalByFiltersInput{}

	for _, oFilter := range aFilters {
		if oFilter == nil || oFilter.Field == nil {
			continue
		}

		oValue, oErr := structpb.NewValue(oFilter.Value)
		if oErr != nil {
			continue
		}

		oPbFilter := &pbResourceModel.GameTypeFilter{
			Field: *oFilter.Field,
			Value: oValue,
		}
		if oFilter.Operator != nil {
			oPbFilter.Operator = *oFilter.Operator
		}

		oRequest.Filters = append(oRequest.Filters, oPbFilter)
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.TotalByFilters(oSelf.Context, oRequest)

	iTotal := oResponse.GetTotal()

	return iTotal, oErr
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

