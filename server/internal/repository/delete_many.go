package repository

import (
	"context"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
)

func DeleteMany[Schema any](
	ctx context.Context,
	collectionName string,
	filter any,
) error {
	db := database.GetDatabase()
	collection := db.Collection(collectionName)

	_, err := collection.DeleteMany(ctx, filter)
	return err

}
