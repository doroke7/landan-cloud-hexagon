package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

func protoGameTypeToDomainGameType(oProtoGameType *pbResource.GameType) domain.GameType {
	if oProtoGameType == nil {
		return domain.GameType{}
	}

	oGameType := domain.GameType{
		Id:        uint64(oProtoGameType.GetId()),
		ParentId:  uint64(oProtoGameType.GetParentId()),
		Key:       oProtoGameType.GetKey(),
		Name:      oProtoGameType.GetName(),
		CreatedAt: oProtoGameType.GetCreatedAt().AsTime(),
		UpdatedAt: oProtoGameType.GetUpdatedAt().AsTime(),
		DeletedAt: oProtoGameType.GetDeletedAt().AsTime(),
	}

	if oParent := oProtoGameType.GetParent(); oParent != nil {
		oParentDomain := protoGameTypeToDomainGameType(oParent)
		oGameType.Parent = &oParentDomain
	}

	for _, oChild := range oProtoGameType.GetChildren() {
		oGameType.Children = append(oGameType.Children, protoGameTypeToDomainGameType(oChild))
	}

	return oGameType
}

type GameTypeModel struct {
	*AbstractModel
}

func NewGameTypeModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) AddOne(oGameTypeParm *domain.GameTypeValue) error {

	oRequest := &pbResourceModel.GameTypeAddOneInput{
		Variable: &pbResourceModel.GameTypeVariable{
			Key:  oGameTypeParm.Key,
			Name: oGameTypeParm.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.GameType.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameTypeModel) ShowOnes() ([]*domain.GameType, error) {
	iSize := uint(10000)
	iPage := uint(1)
	return oSelf.ShowOnesByFiltersWithSortersPagination(nil, nil, &pkgInput.Pagination{Size: &iSize, Page: &iPage})
}

// ShowOnesByParentId 撈出指定父類型底下、尚未刪除的子類型（給刪除前的擋關檢查用）。
// 走既有的 filters 查詢，實際的 deleted_at 過濾由 gRPC 後端底層的持久化 adapter 負責。
func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	sField := "parent_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iParentId},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
}

func (oSelf *GameTypeModel) ShowOnesByParentId(iParentId uint64) ([]*domain.GameType, error) {
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

	for _, oProtoGameType := range oResponse.GetGameTypes() {
		oGameType := protoGameTypeToDomainGameType(oProtoGameType)
		aGameTypes = append(aGameTypes, &oGameType)
	}

	return aGameTypes, nil
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTypeTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *GameTypeModel) EditOneById(oGameType *domain.GameTypeValue, iId uint64) error {

	oRequest := &pbResourceModel.GameTypeEditOneByIdInput{
		Id: uint64(iId),
		Variable: &pbResourceModel.GameTypeVariable{
			Key:  oGameType.Key,
			Name: oGameType.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.GameType.EditOneById(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameTypeModel) RemoveOneById(iId uint64) error {

	_, oErr := oSelf.ResourceModelClient.GameType.RemoveOneById(
		oSelf.Context,
		&pbResourceModel.GameTypeRemoveOneByIdInput{Id: iId},
	)

	return oErr
}
