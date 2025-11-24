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
	primitiveObjectID primitive.ObjectID,
) (T, error) {
	var result T

	collection := database.GetDatabase().Collection(result.GetCollectionName())
	filter := bson.M{"_id": primitiveObjectID}

	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}
