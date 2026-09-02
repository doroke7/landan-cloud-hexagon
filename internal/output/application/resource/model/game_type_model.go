package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	resourceBase "example/internal/output/application/resource"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type GameTypeModel struct {
	*resourceBase.AbstractResource
}

func NewGameTypeModel(oAbstractModel *resourceBase.AbstractResource) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractResource: oAbstractModel,
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

func (oSelf *GameTypeModel) ShowOnes() ([]*domain.GameType, error) {
	iSize := uint(10000)
	iPage := uint(1)
	return oSelf.ShowOnesByFiltersWithSortersPagination(nil, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
}

// ShowOnesByParentId 撈出指定父類型底下、尚未刪除的子類型（給刪除前的擋關檢查用）。
// 走既有的 filters 查詢，實際的 deleted_at 過濾由 gRPC 後端底層的持久化 adapter 負責。
func (oSelf *GameTypeModel) TotalByParentId(iParentId uint) (uint64, error) {
	sField := "parent_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iParentId},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
}

func (oSelf *GameTypeModel) ShowOnesByParentId(iParentId uint) ([]*domain.GameType, error) {
	sField := "parent_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iParentId},
	}

	iSize := uint(10000)
	iPage := uint(1)
	return oSelf.ShowOnesByFiltersWithSortersPagination(aFilters, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
}

func (oSelf *GameTypeModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.GameType, error) {

	oRequest := &pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationInput{
		Filters:    oSelf.ToFilters(aFilters),
		Sorters:    oSelf.ToSorters(aSorters),
		Pagination: oSelf.ToPagination(oPagination),
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

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTypeTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
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
