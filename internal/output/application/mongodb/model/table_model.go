package model

import (
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
)

const TABLE_COLLECTION = "tables"

type tableDocument struct {
	Id          uint      `bson:"_id"`
	No          string    `bson:"no"`
	GameId      uint      `bson:"game_id"`
	Key         string    `bson:"key"`
	State       uint8     `bson:"state"`
	Description string    `bson:"description"`
	Result      string    `bson:"result"`
	StartedAt   time.Time `bson:"started_at"`
	EndedAt     time.Time `bson:"ended_at"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at"`
}

func (oDoc *tableDocument) toDomain() *domain.Table {
	return &domain.Table{
		Id:          oDoc.Id,
		No:          oDoc.No,
		GameId:      oDoc.GameId,
		Key:         oDoc.Key,
		State:       oDoc.State,
		Description: oDoc.Description,
		Result:      json.RawMessage(oDoc.Result),
		StartedAt:   oDoc.StartedAt,
		EndedAt:     oDoc.EndedAt,
		CreatedAt:   oDoc.CreatedAt,
		UpdatedAt:   oDoc.UpdatedAt,
		DeletedAt:   oDoc.DeletedAt,
	}
}

type TableModel struct {
	*AbstractModel
}

func NewTableModel(oAbstractModel *AbstractModel) *TableModel {
	return &TableModel{
		AbstractModel: oAbstractModel,
	}
}

func (oSelf *TableModel) collection() *mongo.Collection {
	return oSelf.Database.Collection(TABLE_COLLECTION)
}

func (oSelf *TableModel) ShowOneById(iId uint) (*domain.Table, error) {
	var oDoc tableDocument

	oFilter := bson.M{"_id": iId, "deleted_at": oNotDeletedAt}

	if oErr := oSelf.collection().FindOne(oSelf.Context, oFilter).Decode(&oDoc); oErr != nil {
		if errors.Is(oErr, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, oErr
	}

	return oDoc.toDomain(), nil
}
