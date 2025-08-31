package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// --- Google Provider Setup ---
	clientID, _ := config.LoadEnv("CLIENT_ID")
	clientSecret, _ := config.LoadEnv("CLIENT_SECRET")

	googleAuth := google.New(
		clientID,
		clientSecret,
		"http://localhost:8080/auth/google/callback",
		"email",
		"profile",
	)
	goth.UseProviders(googleAuth)

	hashKey, _ := config.LoadEnv("SESSION_HASH_KEY")

	blockKey, _ := config.LoadEnv("SESSION_BLOCK_KEY")
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(hashKey), []byte(blockKey))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd
	store.Options.SameSite = http.SameSiteLaxMode
	gothic.Store = store
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
func upsertUser(user goth.User) (*mongo.UpdateResult, error) {
	db := database.GetDatabase()
	usersCollection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.UserID}

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
		return nil, err
	}
	return res, nil
}

// OAuthCallback is a callback function after user login.
func OAuthCallback(c *gin.Context) {
	addProvider(c)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := json.Marshal(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := upsertUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, user)
}

// Login log user in via OAuth.
func Login(c *gin.Context) {
	addProvider(c)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Logout(c *gin.Context) {
	addProvider(c)
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, "User logged out successfully")
}
