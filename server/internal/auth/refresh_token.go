package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshRefreshToken refreshes the refresh token
func RefreshRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	// Check if refresh token is blacklisted
	if IsBlacklisted(refreshToken) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token has been revoked"})
		return
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}

	// Issue new access token
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}
	refreshTokenUser := dto.RefreshTokenUser{
		Email:     claims["email"].(string),
		Name:      claims["name"].(string),
		AvatarURL: claims["avatarURL"].(string),
		ID:        oid,
		Role:      claims["role"].(string),
		Verified:  claims["verified"].(bool),
		Banned:    claims["banned"].(bool),
	}

	accessToken, _, err := generateTokens(refreshTokenUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
