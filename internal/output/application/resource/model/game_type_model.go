package outputApplicationResourceModel

import (
	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
	pb "example/pb"
	pbResourceModel "example/pb/resource/model"
	pkgInput "example/pkg/input"
)

type GameTypeModel struct {
	*AbstractModel
}

func NewGameTypeModel(oAbstractModel *AbstractModel) outputPortAnyModel.GameTypeModel {
	return &GameTypeModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *GameTypeModel) AddOne(oGameTypeParm *domain.GameTypeVariable) error {

	oRequest := &pbResourceModel.GameTypeAddOneInput{
		Variable: &pb.GameTypeVariable{
			Key:  oGameTypeParm.Key,
			Name: oGameTypeParm.Name,
		},
	}

	_, oErr := oSelf.ResourceModelClient.GameType.AddOne(oSelf.Context, oRequest)

	return oErr
}

func (oSelf *GameTypeModel) TotalByParentId(iParentId uint64) (uint64, error) {
	sField := "parent_id"
	sOperator := "eq"
	aFilters := []*pkgInput.Filter{
		{Field: &sField, Operator: &sOperator, Value: iParentId},
	}

	iTotal, oErr := oSelf.TotalByFilters(aFilters)
	return iTotal, oErr
}

func (oSelf *GameTypeModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {

	oRequest := &pbResourceModel.GameTypeTotalByFiltersInput{
		Filters: oSelf.ToFilters(aFilters),
	}

	oResponse, oErr := oSelf.ResourceModelClient.GameType.TotalByFilters(oSelf.Context, oRequest)

	iTotal := uint64(oResponse.GetTotal())

	return iTotal, oErr
}

func (oSelf *GameTypeModel) EditOneById(oGameType *domain.GameTypeVariable, iId uint64) error {

	oRequest := &pbResourceModel.GameTypeEditOneByIdInput{
		Id: uint64(iId),
		Variable: &pb.GameTypeVariable{
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
