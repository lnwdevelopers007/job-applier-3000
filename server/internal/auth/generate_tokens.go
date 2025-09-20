package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
)

var jwtSecret = []byte(config.LoadEnv("JWT_SECRET"))

func generateTokens(email, name, avatarURL string) (accessToken, refreshToken string, err error) {
	// Access token (15m)
	accessClaims := jwt.MapClaims{
		"email":     email,
		"name":      name,
		"avatarURL": avatarURL,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)
	if err != nil {
		return
	}

	expDays := time.Duration(config.LoadInt("REFRESH_TOKEN_AGE_DAYS"))
	// Refresh token (7d)
	refreshClaims := jwt.MapClaims{
		"email":     email,
		"name":      name,
		"avatarURL": avatarURL,
		"type":      "refresh",
		"exp":       time.Now().Add(expDays * 24 * time.Hour).Unix(),
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)

	return
}
