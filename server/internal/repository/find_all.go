package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findAll finds all document which matched the filter from a collection.
// note: opts is an optional parameter.
func FindAll[Schema any](
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
