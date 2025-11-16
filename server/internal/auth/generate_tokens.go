package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
)

var jwtSecret = []byte(config.LoadEnv("JWT_SECRET"))

func generateTokens(email, name, avatarURL string, userID any) (accessToken, refreshToken string, err error) {
	// Access token (15m)
	user, err := findUser(userID)
	if err != nil {
		return
	}
	accessClaims := jwt.MapClaims{
		"email":     email,
		"name":      name,
		"avatarURL": avatarURL,
		"userID":    userID,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"role":      user.Role,
		"verified":  user.Verified,
		"banned":    user.Banned,
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)

	if !config.LoadBoolean("IS_PROD") {
		fmt.Println(accessToken)
	}

	if err != nil {
		return
	}

	expDays := time.Duration(config.LoadInt("REFRESH_TOKEN_AGE_DAYS"))
	// Refresh token (7d)
	refreshClaims := jwt.MapClaims{
		"email":     email,
		"name":      name,
		"avatarURL": avatarURL,
		"userID":    userID,
		"type":      "refresh",
		"exp":       time.Now().Add(expDays * 24 * time.Hour).Unix(),
		"role":      user.Role,
		"verified":  user.Verified,
		"banned":    user.Banned,
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)

	return
}
