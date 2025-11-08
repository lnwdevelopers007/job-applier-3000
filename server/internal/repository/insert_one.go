package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne[Schema any](ctx context.Context, collectionName string, raw Schema) (*mongo.InsertOneResult, error) {

	db := database.GetDatabase()
	collection := db.Collection(collectionName)

	res, err := collection.InsertOne(ctx, raw)
	return res, err

}
