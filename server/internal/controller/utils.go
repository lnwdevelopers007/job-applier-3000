package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// extractUnique extracts unique values from a field of a given slice (list) of structs
func extractUnique[Schema any, ValueType comparable](items []Schema, getValue func(Schema) ValueType) []ValueType {
	visited := make(map[ValueType]bool)
	uniqueSlice := make([]ValueType, 0, len(items))

	for _, item := range items {
		value := getValue(item)
		if !visited[value] {
			visited[value] = true
			uniqueSlice = append(uniqueSlice, value)
		}
	}
	return uniqueSlice
}

// getUsersFromIDs returns a map of users from a slice of IDs.
func getUsersFromIDs(
	ctx context.Context,
	userIDs []primitive.ObjectID,
) (
	map[primitive.ObjectID]schema.User,
	error,
) {
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	result, err := repository.FindAll[schema.User](ctx, filter)

	// Create a map of job seekers for easy lookup
	resultMap := make(map[primitive.ObjectID]schema.User)
	for _, js := range result {
		resultMap[js.ID] = js
	}

	return resultMap, err
}

// getUserFromMiddleware returns user injected from the middleware, which gets the user from jwt.
func getUserFromMiddleware(c *gin.Context) (userID primitive.ObjectID, role string, err error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		err = http.ErrNotSupported
		return
	}

	userID, err = primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		return
	}

	roleVal, _ := c.Get("role")
	role, _ = roleVal.(string)
	return userID, role, nil
}

// getFakeUser generates a new userID from void.
func getFakeUser(c *gin.Context) (userID primitive.ObjectID, role string, err error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		userID = primitive.NewObjectID()
		role = "jobSeeker"
		err = nil
		return
	}

	userID, err = primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		userID = primitive.NewObjectID()
	}

	roleVal, _ := c.Get("role")
	role, _ = roleVal.(string)
	if role == "" {
		role = "jobSeeker"
	}
	err = nil

	return userID, role, nil

}
