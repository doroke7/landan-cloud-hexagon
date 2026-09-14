package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResource "example/pb"
	pbResourceLogic "example/pb/resource/logic"
	pkgDomainToProto "example/pkg/domain_to_proto"
	pkgInput "example/pkg/input"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedGameTypeLogicServer
	LogicGameTypeUsecase usecasePortAnyLogic.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyLogic.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler:      oAbstractHandler,
		LogicGameTypeUsecase: oGameTypeUsecase,
	}
}

// ShowTree 從 usecase 拿到組好的 tree，攤平成平的 GameType 回傳，
// client 端再自己組回 tree（組 tree 每個 adapter 各寫一份）。
func (oSelf *GameTypeHandler) ShowTree(oContext context.Context, oReq *pbResourceLogic.GameTypeShowTreeInput) (*pbResourceLogic.GameTypeShowTreeOutput, error) {

	aRoots, oErr := oSelf.LogicGameTypeUsecase.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.GameType, 0, len(aRoots))

	var fnFlatten func(oNode domain.GameType)
	fnFlatten = func(oNode domain.GameType) {
		aNodes = append(aNodes, &pbResource.GameType{
			Id:        uint64(oNode.Id),
			ParentId:  uint64(oNode.ParentId),
			Key:       oNode.Key,
			Name:      oNode.Name,
			CreatedAt: timestamppb.New(oNode.CreatedAt),
			UpdatedAt: timestamppb.New(oNode.UpdatedAt),
			DeletedAt: timestamppb.New(oNode.DeletedAt),
		})

		for _, oChild := range oNode.Children {
			fnFlatten(oChild)
		}
	}

	for _, oRoot := range aRoots {
		fnFlatten(*oRoot)
	}

	return &pbResourceLogic.GameTypeShowTreeOutput{GameTypes: aNodes}, nil
}

func (oSelf *GameTypeHandler) ShowGameTypeById(oContext context.Context, oReq *pbResourceLogic.GameTypeShowGameTypeByIdInput) (*pbResourceLogic.GameTypeShowGameTypeByIdOutput, error) {

	oGameType, oErr := oSelf.LogicGameTypeUsecase.ShowGameTypeById(oReq.GetId())
	if oErr != nil {
		return nil, oErr
	}

	if oGameType == nil {
		return nil, nil
	}

	oProtoGameType := pkgDomainToProto.GameType(oGameType)

	return &pbResourceLogic.GameTypeShowGameTypeByIdOutput{
		GameType: oProtoGameType,
	}, nil
}

func (oSelf *GameTypeHandler) ShowGameTypes(oContext context.Context, oReq *pbResourceLogic.GameTypeShowGameTypesInput) (*pbResourceLogic.GameTypeShowGameTypesOutput, error) {

	aGameTypes, oErr := oSelf.LogicGameTypeUsecase.ShowGameTypes()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.GameType, 0, len(aGameTypes))
	for _, oOne := range aGameTypes {
		oProtoGameType := pkgDomainToProto.GameType(oOne)
		aNodes = append(aNodes, oProtoGameType)
	}

	return &pbResourceLogic.GameTypeShowGameTypesOutput{GameTypes: aNodes}, nil
}

func (oSelf *GameTypeHandler) ShowGameTypesTotalByFiltersWithSortersPagination(oContext context.Context, oReq *pbResourceLogic.GameTypeShowGameTypesTotalByFiltersWithSortersPaginationInput) (*pbResourceLogic.GameTypeShowGameTypesTotalByFiltersWithSortersPaginationOutput, error) {

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

	aGameTypes, iTotal, oErr := oSelf.LogicGameTypeUsecase.ShowGameTypesTotalByFiltersWithSortersPagination(aFilters, aSorters, oPagination)
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResource.GameType, 0, len(aGameTypes))
	for _, oOne := range aGameTypes {
		aNodes = append(aNodes, &pbResource.GameType{
			Id:        uint64(oOne.Id),
			ParentId:  uint64(oOne.ParentId),
			Key:       oOne.Key,
			Name:      oOne.Name,
			CreatedAt: timestamppb.New(oOne.CreatedAt),
			UpdatedAt: timestamppb.New(oOne.UpdatedAt),
			DeletedAt: timestamppb.New(oOne.DeletedAt),
		})
	}

	return &pbResourceLogic.GameTypeShowGameTypesTotalByFiltersWithSortersPaginationOutput{
		GameTypes: aNodes,
		Total:     iTotal,
	}, nil
}
