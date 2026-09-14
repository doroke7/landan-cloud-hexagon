package inputApplicationFacadeAdminResource

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pb "example/pb"
	pbFacadeAdminResource "example/pb/facade/admin/resource"
	pkgInput "example/pkg/input"
	pkgProtoToDomain "example/pkg/proto_to_domain"
	pkgUtility "example/pkg/utility"
)

// GameHandler 是 http admin resource GameHandler 的 gRPC facade 版本，
// 走同一個 usecasePortAnyAdminResource.GameUsecase，只差在 in/out 是 proto。
type GameHandler struct {
	pbFacadeAdminResource.UnimplementedGameControllerServer
	*inputApplicationFacade.AbstractHandler
	AdminResourceGameUsecase usecasePortAnyAdminResource.GameUsecase
}

func NewGameHandler(oGameUsecase usecasePortAnyAdminResource.GameUsecase, oAbstractHandler *inputApplicationFacade.AbstractHandler) *GameHandler {
	return &GameHandler{
		AbstractHandler:          oAbstractHandler,
		AdminResourceGameUsecase: oGameUsecase,
	}
}

// toStatusError DefaultError（業務錯誤）用 Aborted，其餘用 Internal。
func toStatusError(oErr error) error {
	if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {
		sDefaultError := oDefaultError.Error()
		oStatusError := status.Error(codes.Aborted, sDefaultError)

		return oStatusError
	}

	sError := oErr.Error()
	oStatusError := status.Error(codes.Internal, sError)

	return oStatusError
}

// domainGameTypeToProtoGameType 遞迴帶出 Parent / Children。
func domainGameTypeToProtoGameType(oGameType *domain.GameType) *pb.GameType {
	if oGameType == nil {
		return nil
	}

	oProto := &pb.GameType{
		Id:        uint64(oGameType.Id),
		ParentId:  uint64(oGameType.ParentId),
		Key:       oGameType.Key,
		Name:      oGameType.Name,
		CreatedAt: timestamppb.New(oGameType.CreatedAt),
		UpdatedAt: timestamppb.New(oGameType.UpdatedAt),
		DeletedAt: timestamppb.New(oGameType.DeletedAt),
		Parent:    domainGameTypeToProtoGameType(oGameType.Parent),
	}

	for i := range oGameType.Children {
		oChildProto := domainGameTypeToProtoGameType(&oGameType.Children[i])
		oProto.Children = append(oProto.Children, oChildProto)
	}

	return oProto
}

func domainGameToProtoGame(oGame *domain.Game) *pb.Game {
	if oGame == nil {
		return nil
	}

	return &pb.Game{
		Id:          uint64(oGame.Id),
		GameTypeId:  uint64(oGame.GameTypeId),
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
		CreatedAt:   timestamppb.New(oGame.CreatedAt),
		UpdatedAt:   timestamppb.New(oGame.UpdatedAt),
		DeletedAt:   timestamppb.New(oGame.DeletedAt),
		GameType:    domainGameTypeToProtoGameType(&oGame.GameType),
	}
}

// idFromFilters 從 filters[0]（field 必須是 "id"）取出 id，比照 http admin resource handler 的做法。
func idFromFilters(aFilters []*pb.Filter) (uint, error) {
	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].GetField() != "id" {
		oStatusError := status.Error(codes.InvalidArgument, "filter position error")

		return 0, oStatusError
	}

	oValue := aFilters[0].GetValue()
	oInterfaceValue := oValue.AsInterface()
	fId, bOk := oInterfaceValue.(float64)
	if !bOk {
		oStatusError := status.Error(codes.InvalidArgument, "filter.id format error")

		return 0, oStatusError
	}

	return uint(fId), nil
}

func (oSelf *GameHandler) AddOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameAddOneRequest) (*pbFacadeAdminResource.GameAddOneResponse, error) {

	oRequestVariable := oRequest.GetVariable()
	oDomainGameVariable := pkgProtoToDomain.GameVariable(oRequestVariable)
	oValue := &oDomainGameVariable

	if oErr := oSelf.AdminResourceGameUsecase.AddOne(oValue); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameAddOneResponse{}, nil
}

func (oSelf *GameHandler) EditOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameEditOneRequest) (*pbFacadeAdminResource.GameEditOneResponse, error) {

	aRequestFilters := oRequest.GetFilters()
	iId, oErr := idFromFilters(aRequestFilters)
	if oErr != nil {
		return nil, oErr
	}

	oRequestVariable := oRequest.GetVariable()
	oDomainGameVariable := pkgProtoToDomain.GameVariable(oRequestVariable)
	oValue := &oDomainGameVariable

	if oErr := oSelf.AdminResourceGameUsecase.EditOne(oValue, uint64(iId)); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameEditOneResponse{}, nil
}

func (oSelf *GameHandler) RemoveOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameRemoveOneRequest) (*pbFacadeAdminResource.GameRemoveOneResponse, error) {

	aRequestFilters := oRequest.GetFilters()
	iId, oErr := idFromFilters(aRequestFilters)
	if oErr != nil {
		return nil, oErr
	}

	if oErr := oSelf.AdminResourceGameUsecase.RemoveOne(iId); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameRemoveOneResponse{}, nil
}

func (oSelf *GameHandler) ShowOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameShowOneRequest) (*pbFacadeAdminResource.GameShowOneResponse, error) {

	aRequestFilters := oRequest.GetFilters()
	iId, oErr := idFromFilters(aRequestFilters)
	if oErr != nil {
		return nil, oErr
	}

	oGame, oErr := oSelf.AdminResourceGameUsecase.ShowOne(uint64(iId))
	if oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	oProtoGame := domainGameToProtoGame(oGame)

	return &pbFacadeAdminResource.GameShowOneResponse{
		One: oProtoGame,
	}, nil
}

func (oSelf *GameHandler) ShowOnes(oContext context.Context, oRequest *pbFacadeAdminResource.GameShowOnesRequest) (*pbFacadeAdminResource.GameShowOnesResponse, error) {

	aFiltersInput := oRequest.GetFilters()
	aFilters := make([]*pkgInput.Filter, 0, len(aFiltersInput))
	for _, oOne := range aFiltersInput {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOperator := oOne.GetOperator()
		oValue := oOne.GetValue().AsInterface()
		oFilter := &pkgInput.Filter{
			Field:    &sField,
			Operator: &sOperator,
			Value:    oValue,
		}
		aFilters = append(aFilters, oFilter)
	}

	aSortersInput := oRequest.GetSorters()
	aSorters := make([]*pkgInput.Sorter, 0, len(aSortersInput))
	for _, oOne := range aSortersInput {
		if oOne == nil {
			continue
		}

		sField := oOne.GetField()
		sOrder := oOne.GetOrder()
		oSorter := &pkgInput.Sorter{
			Field: &sField,
			Order: &sOrder,
		}
		aSorters = append(aSorters, oSorter)
	}

	oPaginationInput := oRequest.GetPagination()
	iSize := uint(oPaginationInput.GetSize())
	iPage := uint(oPaginationInput.GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aGames, iTotal, oErr := oSelf.AdminResourceGameUsecase.ShowOnes(aFilters, aSorters, oPagination)
	if oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	aProtoGames := make([]*pb.Game, 0, len(aGames))
	for _, oGame := range aGames {
		oProtoGame := domainGameToProtoGame(oGame)
		aProtoGames = append(aProtoGames, oProtoGame)
	}

	return &pbFacadeAdminResource.GameShowOnesResponse{
		Ones:  aProtoGames,
		Total: uint64(iTotal),
	}, nil
}
