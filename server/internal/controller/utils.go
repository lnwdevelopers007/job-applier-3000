package controller

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
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
