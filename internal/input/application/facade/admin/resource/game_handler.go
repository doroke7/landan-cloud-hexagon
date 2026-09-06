package inputApplicationFacadeAdminResource

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pbFacade "example/pb/facade"
	pbFacadeAdminResource "example/pb/facade/admin/resource"
	pkgInput "example/pkg/input"
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
		return status.Error(codes.Aborted, oDefaultError.Error())
	}

	return status.Error(codes.Internal, oErr.Error())
}

// domainGameTypeToProtoGameType 遞迴帶出 Parent / Children。
func domainGameTypeToProtoGameType(oGameType *domain.GameType) *pbFacade.GameType {
	if oGameType == nil {
		return nil
	}

	oProto := &pbFacade.GameType{
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
		oProto.Children = append(oProto.Children, domainGameTypeToProtoGameType(&oGameType.Children[i]))
	}

	return oProto
}

func domainGameToProtoGame(oGame *domain.Game) *pbFacade.Game {
	if oGame == nil {
		return nil
	}

	return &pbFacade.Game{
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
func idFromFilters(aFilters []*pbFacade.Filter) (uint, error) {
	if len(aFilters) == 0 || aFilters[0] == nil || aFilters[0].GetField() != "id" {
		return 0, status.Error(codes.InvalidArgument, "filter position error")
	}

	fId, bOk := aFilters[0].GetValue().AsInterface().(float64)
	if !bOk {
		return 0, status.Error(codes.InvalidArgument, "filter.id format error")
	}

	return uint(fId), nil
}

func protoGameValueToDomainGameValue(oValue *pbFacadeAdminResource.GameValue) *domain.GameValue {
	oDomain := &domain.GameValue{}
	if oValue == nil {
		return oDomain
	}

	if oValue.GameTypeId != nil {
		iGameTypeId := uint64(oValue.GetGameTypeId())
		oDomain.GameTypeId = &iGameTypeId
	}

	oDomain.Key = oValue.Key
	oDomain.Name = oValue.Name
	oDomain.Description = oValue.Description

	return oDomain
}

func (oSelf *GameHandler) AddOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameAddOneRequest) (*pbFacadeAdminResource.GameAddOneResponse, error) {

	oValue := protoGameValueToDomainGameValue(oRequest.GetValue())

	if oErr := oSelf.AdminResourceGameUsecase.AddOne(oValue); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameAddOneResponse{}, nil
}

func (oSelf *GameHandler) EditOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameEditOneRequest) (*pbFacadeAdminResource.GameEditOneResponse, error) {

	iId, oErr := idFromFilters(oRequest.GetFilters())
	if oErr != nil {
		return nil, oErr
	}

	oValue := protoGameValueToDomainGameValue(oRequest.GetValue())

	if oErr := oSelf.AdminResourceGameUsecase.EditOne(oValue, uint64(iId)); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameEditOneResponse{}, nil
}

func (oSelf *GameHandler) RemoveOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameRemoveOneRequest) (*pbFacadeAdminResource.GameRemoveOneResponse, error) {

	iId, oErr := idFromFilters(oRequest.GetFilters())
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

	iId, oErr := idFromFilters(oRequest.GetFilters())
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

	aFilters := make([]*pkgInput.Filter, 0, len(oRequest.GetFilters()))
	for _, oOne := range oRequest.GetFilters() {
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

	aSorters := make([]*pkgInput.Sorter, 0, len(oRequest.GetSorters()))
	for _, oOne := range oRequest.GetSorters() {
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

	iSize := uint(oRequest.GetPagination().GetSize())
	iPage := uint(oRequest.GetPagination().GetPage())
	oPagination := &pkgInput.Pagination{
		Size: &iSize,
		Page: &iPage,
	}

	aGames, iTotal, oErr := oSelf.AdminResourceGameUsecase.ShowOnes(aFilters, aSorters, oPagination)
	if oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	aProtoGames := make([]*pbFacade.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aProtoGames = append(aProtoGames, domainGameToProtoGame(oGame))
	}

	return &pbFacadeAdminResource.GameShowOnesResponse{
		Ones:  aProtoGames,
		Total: uint64(iTotal),
	}, nil
}
