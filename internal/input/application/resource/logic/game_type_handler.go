package inputApplicationResourceLogic

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "example/internal/domain"
	inputApplicationResource "example/internal/input/application/resource"
	usecasePortAnyLogic "example/internal/usecase/port/any/logic"
	pbResourceLogic "example/pb/resource/logic"
)

type GameTypeHandler struct {
	*inputApplicationResource.AbstractHandler
	pbResourceLogic.UnimplementedGameTypeLogicServer
	usecasePortAnyLogic.GameTypeUsecase
}

func NewGameTypeHandler(oAbstractHandler *inputApplicationResource.AbstractHandler, oGameTypeUsecase usecasePortAnyLogic.GameTypeUsecase) *GameTypeHandler {
	return &GameTypeHandler{
		AbstractHandler: oAbstractHandler,
		GameTypeUsecase: oGameTypeUsecase,
	}
}

// ShowTree 從 usecase 拿到組好的 tree，攤平成平的 GameTypeNode 回傳，
// client 端再自己組回 tree（組 tree 每個 adapter 各寫一份）。
func (oSelf *GameTypeHandler) ShowTree(oContext context.Context, oReq *pbResourceLogic.GameTypeShowTreeInput) (*pbResourceLogic.GameTypeShowTreeOutput, error) {

	aRoots, oErr := oSelf.GameTypeUsecase.ShowTree()
	if oErr != nil {
		return nil, oErr
	}

	aNodes := make([]*pbResourceLogic.GameTypeNode, 0, len(aRoots))

	var fnFlatten func(oNode domain.GameType)
	fnFlatten = func(oNode domain.GameType) {
		aNodes = append(aNodes, &pbResourceLogic.GameTypeNode{
			Id:        uint32(oNode.Id),
			ParentId:  uint32(oNode.ParentId),
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
