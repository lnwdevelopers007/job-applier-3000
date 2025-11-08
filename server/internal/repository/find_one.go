package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindOne[Schema any](
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
