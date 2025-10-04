package auth

import (
	"fmt"
	"net/http"

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
	fmt.Println(role)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res, err := upsertUser(user, role)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(res)
	accessToken, refreshToken, err := generateTokens(user.Email, user.Name, user.AvatarURL, res)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
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

	// Build redirect URL for frontend
	redirectURL := config.LoadEnv("FRONTEND") + "/callback?token=" + accessToken

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
