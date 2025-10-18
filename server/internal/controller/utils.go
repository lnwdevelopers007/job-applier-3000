package controller

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// extractUnique extracts unique values from a field of a given slice (list) of structs
func extractUnique[Schema any, ValueType comparable](items []Schema, getValue func(Schema) ValueType) []ValueType {
	visited := make(map[ValueType]bool)
	uniqueSlice := make([]ValueType, 0, len(items))

	for _, item := range items {
		value := getValue(item)
		if !visited[value] {
			visited[value] = true
			uniqueSlice = append(uniqueSlice, value)
		}
	}
	return uniqueSlice
}

// getUsersFromID returns a map of users from a slice of IDs.
func getUsersFromID(
	ctx context.Context,
	userIDs []primitive.ObjectID,
) (
	map[primitive.ObjectID]schema.User,
	error,
) {
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	result, err := findAll[schema.User](ctx, "users", filter)

	// Create a map of job seekers for easy lookup
	resultMap := make(map[primitive.ObjectID]schema.User)
	for _, js := range result {
		resultMap[js.ID] = js
	}

	return resultMap, err
}

// findAll finds all document which matched the filter from a collection.
// note: opts is an optional parameter.
func findAll[Schema any](
	ctx context.Context,
	collectionName string,
	filter bson.M,
	opts ...*options.FindOptions,
) ([]Schema, error) {
	collection := database.GetDatabase().Collection(collectionName)
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var result []Schema
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func findOne[Schema any](
	ctx context.Context,
	collectionName string,
	objID primitive.ObjectID,
) (Schema, error) {
	collection := database.GetDatabase().Collection(collectionName)
	var result Schema
	err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		var thing Schema
		return thing, err
	}
	return result, nil
}

func deleteMany[Schema any](
	ctx context.Context,
	collectionName string,
	filter any,
) error {
	db := database.GetDatabase()
	collection := db.Collection(collectionName)

	_, err := collection.DeleteMany(ctx, filter)
	return err

}
