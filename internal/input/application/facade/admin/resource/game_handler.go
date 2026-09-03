package inputApplicationFacadeAdminResource

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationFacade "example/internal/input/application/facade"
	usecasePortAnyAdminResource "example/internal/usecase/port/any/admin/resource"
	pbFacadeAdminResource "example/pb/facade/admin/resource"
	pkgInput "example/pkg/input"
	pkgUtility "example/pkg/utility"
)

// GameHandler 是 http admin resource GameHandler 的 gRPC facade 版本，
// 走同一個 usecasePortAnyAdminResource.GameUsecase，只差在 in/out 是 proto。
type GameHandler struct {
	pbFacadeAdminResource.UnimplementedGameControllerServer
	*inputApplicationFacade.AbstractHandler
	GameUsecase usecasePortAnyAdminResource.GameUsecase
}

func NewGameHandler(oGameUsecase usecasePortAnyAdminResource.GameUsecase, oAbstractHandler *inputApplicationFacade.AbstractHandler) *GameHandler {
	return &GameHandler{
		AbstractHandler: oAbstractHandler,
		GameUsecase:     oGameUsecase,
	}
}

// toStatusError DefaultError（業務錯誤）用 Aborted，其餘用 Internal。
func toStatusError(oErr error) error {
	if oDefaultError, bOk := oErr.(*pkgUtility.DefaultError); bOk {
		return status.Error(codes.Aborted, oDefaultError.Error())
	}

	return status.Error(codes.Internal, oErr.Error())
}

func domainGameToProtoGame(oGame *domain.Game) *pbFacadeAdminResource.Game {
	if oGame == nil {
		return nil
	}

	return &pbFacadeAdminResource.Game{
		Id:          uint32(oGame.Id),
		GameTypeId:  uint32(oGame.GameTypeId),
		Key:         oGame.Key,
		Name:        oGame.Name,
		Description: oGame.Description,
		CreatedAt:   timestamppb.New(oGame.CreatedAt),
		UpdatedAt:   timestamppb.New(oGame.UpdatedAt),
		DeletedAt:   timestamppb.New(oGame.DeletedAt),
	}
}

func protoGameValueToDomainGameValue(oValue *pbFacadeAdminResource.GameValue) *domain.GameValue {
	oDomain := &domain.GameValue{}
	if oValue == nil {
		return oDomain
	}

	if oValue.GameTypeId != nil {
		iGameTypeId := uint(oValue.GetGameTypeId())
		oDomain.GameTypeId = &iGameTypeId
	}

	oDomain.Key = oValue.Key
	oDomain.Name = oValue.Name
	oDomain.Description = oValue.Description

	return oDomain
}

func (oSelf *GameHandler) AddOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameAddOneRequest) (*pbFacadeAdminResource.GameAddOneResponse, error) {

	oValue := protoGameValueToDomainGameValue(oRequest.GetValue())

	if _, oErr := oSelf.GameUsecase.AddOne(oValue); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameAddOneResponse{}, nil
}

func (oSelf *GameHandler) EditOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameEditOneRequest) (*pbFacadeAdminResource.GameEditOneResponse, error) {

	oValue := protoGameValueToDomainGameValue(oRequest.GetValue())

	if _, oErr := oSelf.GameUsecase.EditOne(oValue, uint(oRequest.GetId())); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameEditOneResponse{}, nil
}

func (oSelf *GameHandler) RemoveOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameRemoveOneRequest) (*pbFacadeAdminResource.GameRemoveOneResponse, error) {

	if _, oErr := oSelf.GameUsecase.RemoveOne(uint(oRequest.GetId())); oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	return &pbFacadeAdminResource.GameRemoveOneResponse{}, nil
}

func (oSelf *GameHandler) ShowOne(oContext context.Context, oRequest *pbFacadeAdminResource.GameShowOneRequest) (*pbFacadeAdminResource.GameShowOneResponse, error) {

	oGame, oErr := oSelf.GameUsecase.ShowOne(uint(oRequest.GetId()))
	if oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	oProtoGame := domainGameToProtoGame(oGame)

	return &pbFacadeAdminResource.GameShowOneResponse{
		Game: oProtoGame,
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

	aGames, iTotal, oErr := oSelf.GameUsecase.ShowOnes(aFilters, aSorters, oPagination)
	if oErr != nil {
		oStatusErr := toStatusError(oErr)
		return nil, oStatusErr
	}

	aProtoGames := make([]*pbFacadeAdminResource.Game, 0, len(aGames))
	for _, oGame := range aGames {
		aProtoGames = append(aProtoGames, domainGameToProtoGame(oGame))
	}

	return &pbFacadeAdminResource.GameShowOnesResponse{
		Games: aProtoGames,
		Total: iTotal,
	}, nil
}
