package auth

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/repository"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getAllowedRoles() []string {
	return []string{"jobSeeker", "admin", "faculty", "company"}
}

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
// "role" can be valid role and login.
// Due to me being too lazy to catch every edge cases,
func upsertUser(user goth.User, role string) (dbUser schema.User, isNewUser bool, err error) {
	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"userID": user.UserID}

	var existingUser schema.User
	err = usersCollection.FindOne(ctx, filter).Decode(&existingUser)
	// If there's error when querying user
	if err != nil && err != mongo.ErrNoDocuments {
		return existingUser, false, fmt.Errorf("failed to query user: %w", err)
	}

	// If user tries to log in BUT the user does not exist in the database.
	if role == "login" && err == mongo.ErrNoDocuments {
		return existingUser, true, fmt.Errorf("please register first before using our service")
	}

	// If user tries to register (role != login) but the role is invalid
	if role != "login" && !slices.Contains(getAllowedRoles(), role) {
		return existingUser, false, fmt.Errorf("role is not valid")
	}

	// notice that role is setOnInsert meaning that there's no way a user will have role "login".
	update := bson.M{
		"$set": bson.M{
			"provider":  user.Provider,
			"email":     user.Email,
			"avatarURL": user.AvatarURL,
			"updatedAt": time.Now(),
		},
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
			"name":      user.Name,
			"role":      role, // unchanged behavior
			"verified":  false,
		},
	}

	opts := options.Update().SetUpsert(true)

	res, err := usersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return existingUser, false, fmt.Errorf("failed to upsert user: %w", err)
	}

	// check whether upsert operation went successfully
	if res.UpsertedID != nil {
		isNewUser = true
	}

	return existingUser, isNewUser, nil
}

// find user
func findUser(userID any) (schema.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return repository.FindOne[schema.User](ctx, userID)
}
