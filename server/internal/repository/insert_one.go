package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne[T schema.CollectionEntity](ctx context.Context, raw T) (*mongo.InsertOneResult, error) {

	db := database.GetDatabase()
	collection := db.Collection((*new(T)).GetCollectionName())

	res, err := collection.InsertOne(ctx, raw)
	return res, err

}
