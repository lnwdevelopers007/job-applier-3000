package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// TODO: add OAuth code here
func init() {
	// --- Google Provider Setup ---
	clientID, err := config.LoadEnv("CLIENT_ID")
	if err != nil {
		log.Fatalf("Failed to load CLIENT_ID: %v", err)
	}
	clientSecret, err := config.LoadEnv("CLIENT_SECRET")
	if err != nil {
		log.Fatalf("Failed to load CLIENT_SECRET: %v", err)
	}

	googleAuth := google.New(
		clientID,
		clientSecret,
		"http://localhost:8080/auth/google/callback", // Ensure this matches your Google Cloud console configuration
		"email",
		"profile",
	)
	goth.UseProviders(googleAuth)

	// --- Session Store Setup ---
	// IMPORTANT: These keys should be long, random strings.
	// 32 or 64 bytes are standard lengths.
	hashKey, err := config.LoadEnv("SESSION_HASH_KEY")
	if err != nil {
		log.Fatalf("Failed to load SESSION_HASH_KEY: %v", err)
	}

	blockKey, err := config.LoadEnv("SESSION_BLOCK_KEY")
	if err != nil {
		log.Fatalf("Failed to load SESSION_BLOCK_KEY: %v", err)
	}

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

// callback function
func OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.String(http.StatusBadRequest, "Provider not specified")
	}

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, err := json.Marshal(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func BeginAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.String(http.StatusBadRequest, "Provider not specified")
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}
