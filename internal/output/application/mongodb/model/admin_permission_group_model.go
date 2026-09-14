package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	pkgInput "example/pkg/input"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminPermissionGroupModel struct {
	*AbstractModel
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：key 唯一、name+deleted_at、deleted_at 供查詢用。
func NewAdminPermissionGroupModel(oAbstractModel *AbstractModel) (outputPortAnyModel.AdminPermissionGroupModel, error) {
	oCollection := oAbstractModel.Database.Collection("admin_permission_groups")
	oIndexView := oCollection.Indexes()

	oKeyIndexOptions := options.Index()
	oKeyIndexOptions.SetUnique(true)
	oKeyIndexOptions.SetName("admin_permission_groups-k")

	oNameDeletedAtIndexOptions := options.Index()
	oNameDeletedAtIndexOptions.SetName("admin_permission_groups-n-da")

	oDeletedAtIndexOptions := options.Index()
	oDeletedAtIndexOptions.SetName("admin_permission_groups-da")

	aIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "key", Value: 1}},
			Options: oKeyIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "deleted_at", Value: 1}},
			Options: oNameDeletedAtIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: oDeletedAtIndexOptions,
		},
	}

	if _, oErr := oIndexView.CreateMany(oAbstractModel.Context, aIndexModels); oErr != nil {
		return nil, oErr
	}

	oCounters := oAbstractModel.Database.Collection("counters")

	return &AdminPermissionGroupModel{
		AbstractModel: oAbstractModel,
		Collection:    oCollection,
		Counters:      oCounters,
	}, nil
}

func (oSelf *AdminPermissionGroupModel) nextId() (uint64, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_permission_group"},
		bson.M{"$inc": bson.M{"seq": 1}},
		oUpdateOptions,
	)

	var oCounter struct {
		Seq uint64 `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByParentId(iParentId uint64) ([]*domain.AdminPermissionGroup, error) {
	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, bson.M{
		"parent_id":  iParentId,
		"deleted_at": oDeletedAtZero,
	})
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup
	if oErr := oCursor.All(oSelf.Context, &aAdminPermissionGroups); oErr != nil {
		return nil, oErr
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminPermissionGroup, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aAdminPermissionGroups []*domain.AdminPermissionGroup
	if oErr := oCursor.All(oSelf.Context, &aAdminPermissionGroups); oErr != nil {
		return nil, oErr
	}

	return aAdminPermissionGroups, nil
}

func (oSelf *AdminPermissionGroupModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *AdminPermissionGroupModel) AddOne(oValue *domain.AdminPermissionGroupVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.AdminPermissionGroup{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oValue.Key != nil {
		oNew.Key = *oValue.Key
	}
	if oValue.Name != nil {
		oNew.Name = *oValue.Name
	}
	if oValue.ParentId != nil {
		oNew.ParentId = *oValue.ParentId
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *AdminPermissionGroupModel) EditOneById(oValue *domain.AdminPermissionGroupVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oValue.Key != nil {
		oSet["key"] = *oValue.Key
	}
	if oValue.Name != nil {
		oSet["name"] = *oValue.Name
	}

	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId},
		bson.M{"$set": oSet},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.MatchedCount == 0 {
		return errors.New("0 rows updated")
	}

	return nil
}

func (oSelf *AdminPermissionGroupModel) RemoveOneById(iId uint64) error {
	oNow := time.Now()
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": oNow}},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.ModifiedCount == 0 {
		return errors.New("0 rows deleted")
	}

	return nil
}
