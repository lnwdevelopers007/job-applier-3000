package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteOne[T schema.CollectionEntity](ctx context.Context, objID primitive.ObjectID) (*mongo.DeleteResult, error) {
	db := database.GetDatabase()
	var collEn T
	collection := db.Collection(collEn.GetCollectionName())

	// Delete the document from MongoDB
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

	return result, err

}
