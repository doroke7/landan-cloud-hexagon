package inputApplicationResourceModel

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbResourceModel "example/pb/resource/model"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyModel "example/internal/usecase/port/any/model"
	pkgInput "example/pkg/input"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceModel.UnimplementedGameTypeModelServer
	usecasePortAnyModel.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyModel.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler: oAbstractHandler,
		GameTypeUsecase: oGameTypeUsecase,
	}
}

func domainGameTypeToProtoGameType(oGameType *domain.GameType) *pbResourceModel.GameType {
	return &pbResourceModel.GameType{
		Id:        uint32(oGameType.Id),
		Key:       oGameType.Key,
		Name:      oGameType.Name,
		CreatedAt: timestamppb.New(oGameType.CreatedAt),
		UpdatedAt: timestamppb.New(oGameType.UpdatedAt),
		DeletedAt: timestamppb.New(oGameType.DeletedAt),
	}
}

func (oSelf *GameTypeHandler) AddOne(oContext context.Context, oReq *pbResourceModel.GameTypeAddOneInput) (*pbResourceModel.GameTypeAddOneOutput, error) {

	var oGameTypeValue domain.GameTypeValue

	oGameTypeValue.Key = oReq.Key
	oGameTypeValue.Name = oReq.Name

	_, oErr := oSelf.GameTypeUsecase.AddOne(&oGameTypeValue)

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
		Key:  sKey,
		Name: sName,
	}, nil
}

func (oSelf *GameTypeHandler) ShowOneById(oContext context.Context, oReq *pbResourceModel.GameTypeShowOneByIdInput) (*pbResourceModel.GameTypeShowOneByIdOutput, error) {

	oGameType, oErr := oSelf.GameTypeUsecase.ShowOneById(uint(oReq.Id))

	if oErr != nil {
		return nil, oErr
	}

	if oGameType == nil {
		return nil, nil
	}

	oPbGameType := domainGameTypeToProtoGameType(oGameType)

	return &pbResourceModel.GameTypeShowOneByIdOutput{
		Id:        oPbGameType.Id,
		Key:       oPbGameType.Key,
		Name:      oPbGameType.Name,
		CreatedAt: oPbGameType.CreatedAt,
		UpdatedAt: oPbGameType.UpdatedAt,
		DeletedAt: oPbGameType.DeletedAt,
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

	aGameTypes, oErr := oSelf.GameTypeUsecase.ShowOnesByFiltersWithSortersPagination(aFilters, aSorters, oPagination)

	if oErr != nil {
		return nil, oErr
	}

	aPbGameTypes := make([]*pbResourceModel.GameType, 0, len(aGameTypes))
	for _, oGameType := range aGameTypes {
		aPbGameTypes = append(aPbGameTypes, domainGameTypeToProtoGameType(oGameType))
	}

	return &pbResourceModel.GameTypeShowOnesByFiltersWithSortersPaginationOutput{
		GameTypes: aPbGameTypes,
	}, nil
}

func (oSelf *GameTypeHandler) EditOneById(oContext context.Context, oReq *pbResourceModel.GameTypeEditOneByIdInput) (*pbResourceModel.GameTypeEditOneByIdOutput, error) {

	var oGameTypeValue domain.GameTypeValue

	oGameTypeValue.Key = oReq.Key
	oGameTypeValue.Name = oReq.Name

	_, oErr := oSelf.GameTypeUsecase.EditOneById(&oGameTypeValue, uint64(oReq.Id))
	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeEditOneByIdOutput{
		Status: true,
	}, nil
}

func (oSelf *GameTypeHandler) RemoveOneById(oContext context.Context, oReq *pbResourceModel.GameTypeRemoveOneByIdInput) (*pbResourceModel.GameTypeRemoveOneByIdOutput, error) {

	_, oErr := oSelf.GameTypeUsecase.RemoveOneById(uint(oReq.Id))
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

	iTotal, oErr := oSelf.GameTypeUsecase.TotalByFilters(aFilters)

	if oErr != nil {
		return nil, oErr
	}

	return &pbResourceModel.GameTypeTotalByFiltersOutput{
		Total: iTotal,
	}, nil
}
