package controller

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	result, err := repository.FindAll[schema.User](ctx, "users", filter)

	// Create a map of job seekers for easy lookup
	resultMap := make(map[primitive.ObjectID]schema.User)
	for _, js := range result {
		resultMap[js.ID] = js
	}

	return resultMap, err
}
