package outputApplicationMongodbModel

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	outputApplicationMongodb "example/internal/output/application/mongodb"
	pkgInput "example/pkg/input"

	domain "example/internal/domain"
	outputPortAnyModel "example/internal/output/port/any/model"
)

type AdminUserModel struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：name 唯一、deleted_at 供軟刪除過濾用。
func NewAdminUserModel(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) (outputPortAnyModel.AdminUserModel, error) {
	oCollection := oAbstractMongodb.Database.Collection("admin_users")
	oIndexView := oCollection.Indexes()

	oNameIndexOptions := options.Index()
	oNameIndexOptions.SetUnique(true)
	oNameIndexOptions.SetName("admin_users-name")

	oDeletedAtIndexOptions := options.Index()
	oDeletedAtIndexOptions.SetName("admin_users-da")

	aIndexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: oNameIndexOptions,
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: oDeletedAtIndexOptions,
		},
	}

	if _, oErr := oIndexView.CreateMany(oAbstractMongodb.Context, aIndexModels); oErr != nil {
		return nil, oErr
	}

	oCounters := oAbstractMongodb.Database.Collection("counters")

	return &AdminUserModel{
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oCounters,
	}, nil
}

/*
1. 單一 Document 本身有原子操作
MongoDB 不會讓它們拿到同一個 seq。
所以很多場景根本不需要自己加鎖。
*/

func (oSelf *AdminUserModel) nextId() (uint, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_user"},
		bson.M{"$inc": bson.M{"seq": 1}},
		oUpdateOptions,
	)

	var oCounter struct {
		Seq uint `bson:"seq"`
	}
	if oErr := oResult.Decode(&oCounter); oErr != nil {
		return 0, oErr
	}

	return oCounter.Seq, nil
}

func (oSelf *AdminUserModel) ShowOneByName(sName string) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	oResult := oSelf.Collection.FindOne(oSelf.Context, bson.M{"name": sName})
	oErr := oResult.Decode(&oAdminUser)
	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, errors.New("record not found")
		}
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint64) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	oResult := oSelf.Collection.FindOne(oSelf.Context, bson.M{"_id": iId})
	oErr := oResult.Decode(&oAdminUser)
	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, errors.New("record not found")
		}
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) RemoveOneById(iId uint64) error {
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

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oAdminUser.Name != nil {
		oSet["name"] = *oAdminUser.Name
	}
	if oAdminUser.Password != nil {
		oSet["password"] = *oAdminUser.Password
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

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aAdminUsers []*domain.AdminUser
	if oErr := oCursor.All(oSelf.Context, &aAdminUsers); oErr != nil {
		return nil, oErr
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkgInput.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.AdminUser{
		Id:        uint64(iId),
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oAdminUser.Name != nil {
		oNew.Name = *oAdminUser.Name
	}
	if oAdminUser.Password != nil {
		oNew.Password = *oAdminUser.Password
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}
