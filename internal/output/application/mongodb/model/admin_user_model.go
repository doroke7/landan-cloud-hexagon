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

type AdminUserModel struct {
	*AbstractModel
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

// 索引跟 script/mongodb/resource.js 建的一致：name 唯一、deleted_at 供軟刪除過濾用。
func NewAdminUserModel(oAbstractModel *AbstractModel) (outputPortAnyModel.AdminUserModel, error) {
	oCollection := oAbstractModel.Database.Collection("admin_users")

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractModel.Context, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("admin_users-name"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("admin_users-da"),
		},
	}); oErr != nil {
		return nil, oErr
	}

	return &AdminUserModel{
		AbstractModel: oAbstractModel,
		Collection:    oCollection,
		Counters:      oAbstractModel.Database.Collection("counters"),
	}, nil
}

/*
1. 單一 Document 本身有原子操作
MongoDB 不會讓它們拿到同一個 seq。
所以很多場景根本不需要自己加鎖。
*/

func (oSelf *AdminUserModel) nextId() (uint, error) {
	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_user"},
		bson.M{"$inc": bson.M{"seq": 1}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
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

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{"name": sName}).Decode(&oAdminUser)
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

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{"_id": iId}).Decode(&oAdminUser)
	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, errors.New("record not found")
		}
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserModel) RemoveOneById(iId uint64) error {
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": time.Now()}},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.ModifiedCount == 0 {
		return errors.New("0 rows deleted")
	}

	return nil
}

func (oSelf *AdminUserModel) EditOneById(oAdminUser *domain.AdminUserValue, iId uint64) error {
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

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) error {
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
