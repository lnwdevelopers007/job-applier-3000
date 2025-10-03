package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func upsertUser(user goth.User) (any, error) {
	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userID": user.UserID}

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
		},
	}

	opts := options.Update().SetUpsert(true)

	res, err := usersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", err
	}
	if res.UpsertedID != nil {
		return res.UpsertedID, nil
	}

	var existingUser bson.M
	if err := usersCollection.FindOne(ctx, filter).Decode(&existingUser); err != nil {
		return "", err
	}
	return existingUser["_id"], nil
}

func findUserRole(userID any) (string, error) {
	db := database.GetDatabase()
	jobSeekerCollection := db.Collection("jobSeeker")
	companyCollection := db.Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert to ObjectID if it’s a string
	var oid primitive.ObjectID
	switch v := userID.(type) {
	case string:
		var err error
		oid, err = primitive.ObjectIDFromHex(v)
		if err != nil {
			return "", fmt.Errorf("invalid userID string: %w", err)
		}
	case primitive.ObjectID:
		oid = v
	default:
		return "", fmt.Errorf("unsupported userID type: %T", v)
	}

	// Now query with ObjectID
	var registeredUser bson.M
	filter := bson.M{"userID": oid}
	if err := jobSeekerCollection.FindOne(ctx, filter).Decode(&registeredUser); err == nil {
		return "jobSeeker", nil
	}
	fmt.Println(registeredUser)
	if err := companyCollection.FindOne(ctx, filter).Decode(&registeredUser); err == nil {
		return "company", nil
	}
	fmt.Println(registeredUser)

	return "unverified", nil
}
