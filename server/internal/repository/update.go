package repository

import (
	"context"
	"reflect"
	"strings"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update[collectionEntity schema.CollectionEntity, dto any](
	ctx context.Context, objID primitive.ObjectID, newData dto,
) (*mongo.UpdateResult, error) {

	updateFields := buildUpdateMap(newData)

	update := bson.M{"$set": updateFields}

	db := database.GetDatabase()
	var collEn collectionEntity
	collection := db.Collection(collEn.GetCollectionName())

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return res, err
}

func buildUpdateMap(input any) bson.M {
	result := bson.M{}

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Extract bson tag
		bsonTag := fieldType.Tag.Get("bson")
		bsonKey := strings.Split(bsonTag, ",")[0]

		if bsonKey == "" || bsonKey == "-" {
			continue
		}

		// Skip UpdatedAt — always set manually
		if bsonKey == "updatedAt" {
			continue
		}

		// Only include non-nil pointer fields
		if fieldVal.Kind() == reflect.Pointer && !fieldVal.IsNil() {
			result[bsonKey] = fieldVal.Elem().Interface()
		}
	}

	return result
}
