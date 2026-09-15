package outputApplicationMongodbLogic

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AdminRoleLogic struct {
	*AbstractLogic
	Collection *mongo.Collection
	Counters   *mongo.Collection
}

func NewAdminRoleLogic(oAbstractLogic *AbstractLogic) outputPortAnyLogic.AdminRoleLogic {
	oCollection := oAbstractLogic.Database.Collection("admin_roles")
	oCounters := oAbstractLogic.Database.Collection("counters")

	return &AdminRoleLogic{
		AbstractLogic: oAbstractLogic,
		Collection:    oCollection,
		Counters:      oCounters,
	}
}

func (oSelf *AdminRoleLogic) nextId() (uint64, error) {
	oUpdateOptions := options.FindOneAndUpdate()
	oUpdateOptions.SetUpsert(true)
	oUpdateOptions.SetReturnDocument(options.After)

	oResult := oSelf.Counters.FindOneAndUpdate(
		oSelf.Context,
		bson.M{"_id": "admin_role"},
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

func (oSelf *AdminRoleLogic) AddAdminRole(oVariable *domain.AdminRoleVariable) error {
	iId, oErr := oSelf.nextId()
	if oErr != nil {
		return oErr
	}

	oNow := time.Now()
	oNew := &domain.AdminRole{
		Id:        iId,
		CreatedAt: oNow,
		UpdatedAt: oNow,
		DeletedAt: oDeletedAtZero,
	}

	if oVariable.Key != nil {
		oNew.Key = *oVariable.Key
	}
	if oVariable.Name != nil {
		oNew.Name = *oVariable.Name
	}

	if _, oErr := oSelf.Collection.InsertOne(oSelf.Context, oNew); oErr != nil {
		return oErr
	}

	return nil
}
