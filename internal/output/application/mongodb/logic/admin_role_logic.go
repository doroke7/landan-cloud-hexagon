package outputApplicationMongodbLogic

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "example/internal/domain"
	outputApplicationMongodb "example/internal/output/application/mongodb"
	outputPortAnyLogic "example/internal/output/port/any/logic"
)

type AdminRoleLogic struct {
	*outputApplicationMongodb.AbstractMongodb
	Collection                   *mongo.Collection
	AdminRolesToAdminPermissions *mongo.Collection
	Counters                     *mongo.Collection
}

func NewAdminRoleLogic(oAbstractMongodb *outputApplicationMongodb.AbstractMongodb) outputPortAnyLogic.AdminRoleLogic {
	oCollection := oAbstractMongodb.Database.Collection("admin_roles")
	oAdminRolesToAdminPermissions := oAbstractMongodb.Database.Collection("admin_roles_to_admin_permissions")
	oCounters := oAbstractMongodb.Database.Collection("counters")

	return &AdminRoleLogic{
		AbstractMongodb:              oAbstractMongodb,
		Collection:                   oCollection,
		AdminRolesToAdminPermissions: oAdminRolesToAdminPermissions,
		Counters:                     oCounters,
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

	// 避免髒數據（例如 id 被重用）導致異常，插入前先把這個 admin_role_id 底下的關聯清乾淨
	if _, oErr := oSelf.AdminRolesToAdminPermissions.DeleteMany(oSelf.Context, bson.M{"admin_role_id": iId}); oErr != nil {
		return oErr
	}

	if oVariable.AdminPermissionIds == nil || len(*oVariable.AdminPermissionIds) == 0 {
		return nil
	}

	aAdminRolesToAdminPermissions := make([]interface{}, 0, len(*oVariable.AdminPermissionIds))
	for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
		oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
			AdminRoleId:       iId,
			AdminPermissionId: iAdminPermissionId,
		}
		aAdminRolesToAdminPermissions = append(aAdminRolesToAdminPermissions, oAdminRolesToAdminPermission)
	}

	if _, oErr := oSelf.AdminRolesToAdminPermissions.InsertMany(oSelf.Context, aAdminRolesToAdminPermissions); oErr != nil {
		return oErr
	}

	return nil
}

func (oSelf *AdminRoleLogic) EditAdminRoleById(oVariable *domain.AdminRoleVariable, iId uint64) error {
	oSet := bson.M{"updated_at": time.Now()}

	if oVariable.Key != nil {
		oSet["key"] = *oVariable.Key
	}
	if oVariable.Name != nil {
		oSet["name"] = *oVariable.Name
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
		oZeroRowsError := errors.New("0 rows updated")

		return oZeroRowsError
	}

	if oVariable.AdminPermissionIds == nil {
		return nil
	}

	oCursor, oErr := oSelf.AdminRolesToAdminPermissions.Find(oSelf.Context, bson.M{"admin_role_id": iId})
	if oErr != nil {
		return oErr
	}

	var aExistingAdminRolesToAdminPermissions []domain.AdminRolesToAdminPermission
	if oErr := oCursor.All(oSelf.Context, &aExistingAdminRolesToAdminPermissions); oErr != nil {
		return oErr
	}

	oMapExistingAdminPermissionIds := make(map[uint64]bool, len(aExistingAdminRolesToAdminPermissions))
	for _, oExisting := range aExistingAdminRolesToAdminPermissions {
		oMapExistingAdminPermissionIds[oExisting.AdminPermissionId] = true
	}

	oInputtingAdminPermissionIds := make(map[uint64]bool, len(*oVariable.AdminPermissionIds))
	for _, iAdminPermissionId := range *oVariable.AdminPermissionIds {
		oInputtingAdminPermissionIds[iAdminPermissionId] = true
	}

	// a. 傳進來的 id 在 DB 不存在 -> 插入
	aAdminPermissionIdsToInsert := make([]interface{}, 0, len(oInputtingAdminPermissionIds))
	for iAdminPermissionId := range oInputtingAdminPermissionIds {
		if oMapExistingAdminPermissionIds[iAdminPermissionId] {
			continue
		}

		oAdminRolesToAdminPermission := domain.AdminRolesToAdminPermission{
			AdminRoleId:       iId,
			AdminPermissionId: iAdminPermissionId,
		}
		aAdminPermissionIdsToInsert = append(aAdminPermissionIdsToInsert, oAdminRolesToAdminPermission)
	}

	// b. DB 存在但不在傳進來的 id 裡面 -> 刪除
	aAdminPermissionIdsToDelete := make([]uint64, 0, len(oMapExistingAdminPermissionIds))
	for iAdminPermissionId := range oMapExistingAdminPermissionIds {
		if oInputtingAdminPermissionIds[iAdminPermissionId] {
			continue
		}

		aAdminPermissionIdsToDelete = append(aAdminPermissionIdsToDelete, iAdminPermissionId)
	}

	// 交集的部分不動，維持原樣

	if len(aAdminPermissionIdsToDelete) > 0 {
		oFilter := bson.M{"admin_role_id": iId, "admin_permission_id": bson.M{"$in": aAdminPermissionIdsToDelete}}
		if _, oErr := oSelf.AdminRolesToAdminPermissions.DeleteMany(oSelf.Context, oFilter); oErr != nil {
			return oErr
		}
	}

	if len(aAdminPermissionIdsToInsert) > 0 {
		if _, oErr := oSelf.AdminRolesToAdminPermissions.InsertMany(oSelf.Context, aAdminPermissionIdsToInsert); oErr != nil {
			return oErr
		}
	}

	return nil
}
