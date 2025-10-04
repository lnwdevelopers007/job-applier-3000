package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// addProvider adds the OAuth provider to the request from the query params.
func addProvider(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.String(http.StatusBadRequest, "Provider not specified")
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
}

// upsertUser update or insert user into the database.
// "role" can only have 3 values: company, jobSeeker, login.
// Due to me being too lazy to catch every edge cases,
func upsertUser(user goth.User, role string) (any, error) {
	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userID": user.UserID}

	// Check if user already exists
	var existingUser schema.User
	err := usersCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to query existing user: %w", err)
	}

	// If the user does not exist and is trying to login, throw error.
	if role == "login" {
		if existingUser.ID.IsZero() {
			return nil, fmt.Errorf("please register first before using our service")
		}
	}

	update := bson.M{
		"$set": bson.M{
			"provider":  user.Provider,
			"email":     user.Email,
			"name":      user.Name,
			"avatarURL": user.AvatarURL,
			"updatedAt": time.Now(),
		},
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
			"role":      role,
			"verified":  false,
		},
	}

	opts := options.Update().SetUpsert(true)

	res, err := usersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %w", err)
	}
	if res.UpsertedID != nil {
		return res.UpsertedID, nil
	}

	return existingUser.ID, nil
}

// find user
func findUser(userID any) (schema.User, error) {
	db := database.GetDatabase()
	users := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert to ObjectID if it’s a string
	var oid primitive.ObjectID
	switch v := userID.(type) {
	case string:
		var err error
		oid, err = primitive.ObjectIDFromHex(v)
		if err != nil {
			return schema.User{}, fmt.Errorf("invalid userID string: %w", err)
		}
	case primitive.ObjectID:
		oid = v
	default:
		return schema.User{}, fmt.Errorf("unsupported userID type: %T", v)
	}

	// Now query with ObjectID
	fmt.Println(oid)
	var registeredUser schema.User
	filter := bson.M{"_id": oid}
	if err := users.FindOne(ctx, filter).Decode(&registeredUser); err != nil {
		return registeredUser, fmt.Errorf("can't find user: %w", err)
	}
	return registeredUser, nil

}
