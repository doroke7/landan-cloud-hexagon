package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "example/internal/domain"
)

const TABLE_COLLECTION = "tables"

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
	var oDocument domain.TableDocument

	oFilter := bson.M{"_id": iId, "deleted_at": oNotDeletedAt}

	if oErr := oSelf.collection().FindOne(oSelf.Context, oFilter).Decode(&oDocument); oErr != nil {
		if errors.Is(oErr, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, oErr
	}

	return domain.TableDocumentToTable(&oDocument), nil
}
