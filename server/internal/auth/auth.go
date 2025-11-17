package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	googleAuth := google.New(
		config.LoadEnv("CLIENT_ID"),
		config.LoadEnv("CLIENT_SECRET"),
		config.LoadCallbackURI("http", "google"),
		"email",
		"profile",
	)

	goth.UseProviders(googleAuth)
	maxAgeSeconds := 86400 * config.LoadInt("MAX_COOKIE_AGE_DAYS")
	store := sessions.NewCookieStore(
		[]byte(config.LoadEnv("SESSION_HASH_KEY")),
		[]byte(config.LoadEnv("SESSION_BLOCK_KEY")),
	)

	store.MaxAge(maxAgeSeconds)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = config.LoadBoolean("IS_PROD")
	store.Options.SameSite = http.SameSiteLaxMode
	gothic.Store = store
}

// OAuthCallback is a callback function after user login.
func OAuthCallback(c *gin.Context) {
	addProvider(c)
	role := c.Query("state")
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		msg := "cannot complete user authentication"
		slog.Error(msg + ": " + err.Error())
		c.AbortWithError(http.StatusInternalServerError, errors.New(msg))
		return
	}

	dbUser, isNewUser, err := upsertUser(user, role)
	if err != nil {
		msg := "cannot upsert user"
		slog.Error(msg + ": " + err.Error())
		if isNewUser {
			c.AbortWithError(http.StatusBadRequest, errors.New(msg))
		} else {
			c.AbortWithError(http.StatusInternalServerError, errors.New(msg))
		}
		return
	}

	accessToken, refreshToken, err := generateTokens(user.Email, user.Name, user.AvatarURL, dbUser.ID.Hex())
	if err != nil {
		msg := "cannot generate token"
		slog.Error(msg + err.Error())
		c.AbortWithError(http.StatusInternalServerError, errors.New(msg))
		return
	}

	refreshTokenAge := config.LoadInt("REFRESH_TOKEN_AGE_DAYS") * 24 * 3600
	// Set refresh token as HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		refreshTokenAge,
		"/", "localhost",
		false,
		true,
	)

	// Set access token as HttpOnly cookie (more secure than localStorage)
	// Access tokens typically have shorter lifespan (e.g., 1 hour = 3600 seconds)
	c.SetCookie(
		"access_token",
		accessToken,
		3600, // 1 hour
		"/", "localhost",
		false,
		true,
	)

	slog.Info("User: " + dbUser.ID.Hex() + ", Role: " + dbUser.Role + " logged in")

	if dbUser.Banned {
		// Now redirect to banned page WITH cookies
		redirectURL := config.LoadEnv("FRONTEND") + "/banned"
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Redirect to frontend callback without token in URL
	redirectURL := config.LoadEnv("FRONTEND") + "/callback"
	if isNewUser {
		redirectURL += "?step=signup&token=" + accessToken
	} else {
		redirectURL += "?step=login&token=" + accessToken
	}
	c.Redirect(http.StatusFound, redirectURL)
}

// Login log user in via OAuth.
func Login(c *gin.Context) {
	addProvider(c)
	role := c.Query("role")
	q := c.Request.URL.Query()
	q.Set("state", role)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Logout(c *gin.Context) {
	addProvider(c)
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.SetCookie("refresh_token", "", -1, "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusPermanentRedirect, config.LoadEnv("FRONTEND"))
}

// Me returns the current authenticated user's information from JWT
func Me(c *gin.Context) {
	enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
	if !enableAuth {
		c.JSON(http.StatusOK, gin.H{
			"message": "authentication is disabled in this environment",
		})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	role, _ := c.Get("role")
	email, _ := c.Get("email")
	name, _ := c.Get("name")
	banned, _ := c.Get("banned")

	c.JSON(http.StatusOK, gin.H{
		"userID": userID,
		"role":   role,
		"email":  email,
		"name":   name,
		"banned": banned, // Include ban status
		// Note: 'verified' status is managed by admin and stored in database
		// Query user document if you need verification status
	})
}
