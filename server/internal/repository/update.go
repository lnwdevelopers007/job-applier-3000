package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Update[T schema.CollectionEntity](
	ctx context.Context, objID primitive.ObjectID, newData T,
) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": newData}

	db := database.GetDatabase()
	var collEn T
	collection := db.Collection(collEn.GetCollectionName())

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return res, err
}
