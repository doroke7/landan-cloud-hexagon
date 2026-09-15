package outputApplicationMongodbLogic

import (
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
	pkgInput "example/pkg/input"
)

type AdminUserLogic struct {
	*AbstractLogic
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

func NewAdminUserLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminUserLogic {
	oCollection := oAbstractLogic.Database.Collection("admin_users")
	oCounters := oAbstractLogic.Database.Collection("counters")

	return &AdminUserLogic{
		AbstractLogic: oAbstractLogic,
		Collection:    oCollection,
		Counters:      oCounters,
	}
}

func (oSelf *AdminUserLogic) nextId() (uint, error) {
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

func (oSelf *AdminUserLogic) ShowAdminUsersTotalByFiltersWithSortersPagination(aFilters []*pkgInput.Filter, aSorters []*pkgInput.Sorter, oPagination *pkgInput.Pagination) ([]*domain.AdminUser, uint64, error) {
	oFilter := oSelf.FiltersToFilter(aFilters)
	oFilter["deleted_at"] = oDeletedAtZero

	oFindOptions := oSelf.SortersPaginationToFindOptions(aSorters, oPagination)

	var aAdminUsers []*domain.AdminUser
	var iTotal int64
	var oFindErr error
	var oCountErr error

	var oWaitGroup sync.WaitGroup
	oWaitGroup.Add(2)

	go func() {
		defer oWaitGroup.Done()

		oCursor, oErr := oSelf.Collection.Find(oSelf.Context, oFilter, oFindOptions)
		if oErr != nil {
			oFindErr = oErr
			return
		}
		defer oCursor.Close(oSelf.Context)

		oFindErr = oCursor.All(oSelf.Context, &aAdminUsers)
	}()

	go func() {
		defer oWaitGroup.Done()

		iTotal, oCountErr = oSelf.Collection.CountDocuments(oSelf.Context, oFilter)
	}()

	oWaitGroup.Wait()

	if oFindErr != nil {
		return aAdminUsers, 0, oFindErr
	}

	if oCountErr != nil {
		return aAdminUsers, 0, oCountErr
	}

	return aAdminUsers, uint64(iTotal), nil
}

func (oSelf *AdminUserLogic) AddAdminUser(oValue *domain.AdminUserVariable) error {
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

	if oValue.Name != nil {
		oNew.Name = *oValue.Name
	}
	if oValue.Password != nil {
		oNew.Password = *oValue.Password
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *AdminUserLogic) EditAdminUserById(oValue *domain.AdminUserVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oValue.Name != nil {
		oSet["name"] = *oValue.Name
	}
	if oValue.Password != nil {
		oSet["password"] = *oValue.Password
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

func (oSelf *AdminUserLogic) ShowAdminUserById(iId uint64) (*domain.AdminUser, error) {
	var oAdminUser domain.AdminUser

	oResult := oSelf.Collection.FindOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
	)
	oErr := oResult.Decode(&oAdminUser)

	if errors.Is(oErr, mongo.ErrNoDocuments) {
		return nil, errors.New("record not found")
	}

	if oErr != nil {
		return nil, oErr
	}

	return &oAdminUser, nil
}

func (oSelf *AdminUserLogic) RemoveAdminUserById(iId uint64) error {
	oNow := time.Now()
	oResult, oErr := oSelf.Collection.UpdateOne(
		oSelf.Context,
		bson.M{"_id": iId, "deleted_at": oDeletedAtZero},
		bson.M{"$set": bson.M{"deleted_at": oNow}},
	)
	if oErr != nil {
		return oErr
	}

	if oResult.MatchedCount == 0 {
		return errors.New("0 rows deleted")
	}

	return nil
}
