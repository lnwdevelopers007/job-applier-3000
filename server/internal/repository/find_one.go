package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindOne[T schema.CollectionEntity](
	ctx context.Context,
	objID primitive.ObjectID,
) (T, error) {
	var result T
	collection := database.GetDatabase().Collection(result.GetCollectionName())
	filter := bson.M{"_id": objID}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		var thing T
		return thing, err
	}
	return result, nil
}
