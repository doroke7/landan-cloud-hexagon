package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	pkg "example/pkg"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
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

	if _, oErr := oCollection.Indexes().CreateMany(oAbstractMongodb.Context, []mongo.IndexModel{
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
		AbstractMongodb: oAbstractMongodb,
		Collection:      oCollection,
		Counters:        oAbstractMongodb.Database.Collection("counters"),
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
	var oDoc domain.AdminUserDocument

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{"name": sName}).Decode(&oDoc)
	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return domain.AdminUserDocumentToAdminUser(&oDoc), nil
}

func (oSelf *AdminUserModel) ShowOneById(iId uint) (*domain.AdminUser, error) {
	var oDoc domain.AdminUserDocument

	oErr := oSelf.Collection.FindOne(oSelf.Context, bson.M{"_id": iId}).Decode(&oDoc)
	if oErr != nil {
		if oErr == mongo.ErrNoDocuments {
			return nil, errors.New("資料不存在")
		}
		return nil, oErr
	}

	return domain.AdminUserDocumentToAdminUser(&oDoc), nil
}

func (oSelf *AdminUserModel) ShowOnesByFiltersWithSortersPagination(aFilters []*pkg.Filter, aSorters []*pkg.Sorter, oPagination *pkg.Pagination) ([]*domain.AdminUser, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
	if oErr != nil {
		return nil, oErr
	}
	defer oCursor.Close(oSelf.Context)

	var aDocs []*domain.AdminUserDocument
	if oErr := oCursor.All(oSelf.Context, &aDocs); oErr != nil {
		return nil, oErr
	}

	aAdminUsers := make([]*domain.AdminUser, len(aDocs))
	for i, oDoc := range aDocs {
		aAdminUsers[i] = domain.AdminUserDocumentToAdminUser(oDoc)
	}

	return aAdminUsers, nil
}

func (oSelf *AdminUserModel) TotalByFilters(aFilters []*pkg.Filter) (uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	iCount, oErr := oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	if oErr != nil {
		return 0, oErr
	}

	return uint64(iCount), nil
}

func (oSelf *AdminUserModel) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return false, oErr
	}

	oNow := time.Now()
	oDoc := &domain.AdminUserDocument{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oAdminUser.Name != nil {
		oDoc.Name = *oAdminUser.Name
	}
	if oAdminUser.Password != nil {
		oDoc.Password = *oAdminUser.Password
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oDoc); oErr != nil {
		return false, oErr
	}

	return true, nil
}
