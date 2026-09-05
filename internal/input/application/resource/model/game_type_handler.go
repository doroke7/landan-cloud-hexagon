package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResource "example/pb/resource"
	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameTypeModelServer
	ModelGameTypeUsecase usecasePortAnyModel.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyModel.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler:      oAbstractHandler,
		ModelGameTypeUsecase: oGameTypeUsecase,
	}
}

// domainGameTypeToProtoGameType 遞迴帶出 Parent / Children，讓巢狀的 game type 一起過 gRPC。
func domainGameTypeToProtoGameType(oGameType *domain.GameType) *pbResource.GameType {
	if oGameType == nil {
		return nil
	}

	oPb := &pbResource.GameType{
		Id:        uint32(oGameType.Id),
		ParentId:  uint32(oGameType.ParentId),
		Key:       oGameType.Key,
		Name:      oGameType.Name,
		CreatedAt: timestamppb.New(oGameType.CreatedAt),
		UpdatedAt: timestamppb.New(oGameType.UpdatedAt),
		DeletedAt: timestamppb.New(oGameType.DeletedAt),
		Parent:    domainGameTypeToProtoGameType(oGameType.Parent),
	}

	for i := range oGameType.Children {
		oPb.Children = append(oPb.Children, domainGameTypeToProtoGameType(&oGameType.Children[i]))
	}

	return oPb
}

func (oSelf *GameTypeHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameTypeAddOneInput) (*pbResourceModel.GameTypeAddOneOutput, error) {

	oVariable := oReq.GetVariable()

	var oGameTypeValue domain.GameTypeValue
	if oVariable != nil {
		oGameTypeValue.Key = oVariable.Key
		oGameTypeValue.Name = oVariable.Name
	}

	oErr := oSelf.ModelGameTypeUsecase.AddOne(&oGameTypeValue)

	if oErr != nil {
		return nil, oErr
	}

	var sKey, sName string

	if oGameTypeValue.Key != nil {
		sKey = *oGameTypeValue.Key
	}
	if oGameTypeValue.Name != nil {
		sName = *oGameTypeValue.Name
	}

	return &pbResourceModel.GameTypeAddOneOutput{
		GameType: &pbResource.GameType{
			Key:  sKey,
			Name: sName,
		},
	}, nil
}

func (oSelf *GameTypeHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.GameTypeShowOneByIdInput) (*pbResourceModel.GameTypeShowOneByIdOutput, error) {

	oGameType, oErr := oSelf.ModelGameTypeUsecase.ShowOneById(uint64(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	if oGameType == nil {
		return nil, nil
	}

	oProtoGameType := domainGameTypeToProtoGameType(oGameType)

	return &pbResourceModel.GameTypeShowOneByIdOutput{
		GameType: oProtoGameType,
	}, nil
}

func (oSelf *GameTypeHandler) ShowOnesByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationInput) (*pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationOutput, error) {

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

	aGameTypes, oErr := oSelf.ModelGameTypeUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, oErr
	}

	aPbGameTypes := make([]*pbResource.GameType, 0, len(aGameTypes))
	for _, oGameType := range aGameTypes {
		aPbGameTypes = append(aPbGameTypes, domainGameTypeToProtoGameType(oGameType))
	}

	return &pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationOutput{
		GameTypes: aPbGameTypes,
	}, nil
}

func (oSelf *GameTypeHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameTypeEditOneByIdInput) (*pbResourceModel.GameTypeEditOneByIdOutput, error) {

	oVariable := oReq.GetVariable()

	var oGameTypeValue domain.GameTypeValue
	if oVariable != nil {
		oGameTypeValue.Key = oVariable.Key
		oGameTypeValue.Name = oVariable.Name
	}

	oErr := oSelf.ModelGameTypeUsecase.EditOneById(&oGameTypeValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameTypeHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameTypeRemoveOneByIdInput) (*pbResourceModel.GameTypeRemoveOneByIdOutput, error) {

	oErr := oSelf.ModelGameTypeUsecase.RemoveOneById(uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeRemoveOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameTypeHandler) TotalByFilters(oContext context.Context, oReq *pbResourceModel.GameTypeTotalByFiltersInput) (*pbResourceModel.GameTypeTotalByFiltersOutput, error) {

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

	iTotal, oErr := oSelf.ModelGameTypeUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
