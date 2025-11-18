package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindOne[T schema.CollectionEntity](
	ctx context.Context,
	objID any,
) (T, error) {
	var result T
	primitiveObjectID, err := getPrimitiveObjID(objID)
	if err != nil {
		return result, err
	}
	collection := database.GetDatabase().Collection(result.GetCollectionName())
	filter := bson.M{"_id": primitiveObjectID}

	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func getPrimitiveObjID(objID any) (primitive.ObjectID, error) {
	var oid primitive.ObjectID

	switch v := objID.(type) {
	case string:
		v = strings.TrimSpace(v)
		if len(v) != 24 {
			return oid, fmt.Errorf("invalid ObjectID string: wrong length")
		}
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return oid, fmt.Errorf("invalid ObjectID string: %w", err)
		}
		oid = id
	case primitive.ObjectID:
		oid = v
	default:
		return oid, fmt.Errorf("unsupported ObjectID type: %T", objID)
	}

	return oid, nil
}
