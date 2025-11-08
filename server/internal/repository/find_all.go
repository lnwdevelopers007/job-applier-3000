package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// findAll finds all document which matched the filter from a collection.
// note: opts is an optional parameter.
func FindAll[T schema.CollectionEntity](
	ctx context.Context,
	filter bson.M,
	opts ...*options.FindOptions,
) ([]T, error) {
	var collEn T
	collection := database.GetDatabase().Collection(collEn.GetCollectionName())
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var result []T
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
