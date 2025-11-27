package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/dto"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshRefreshToken refreshes the refresh token
// SECURITY FIX: Now checks fresh ban status from database
func RefreshRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
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

	// SECURITY: Check fresh ban status from database
	db := database.GetDatabase()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var dbUser schema.User
	err = db.Collection("users").FindOne(ctx, bson.M{"_id": oid}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Reject token refresh if user is banned
	if dbUser.Banned {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "account_banned",
			"message": "Your account has been banned. Please contact support.",
		})
		return
	}

	// Use fresh user data from database (includes updated ban/verified status)
	refreshTokenUser := dto.RefreshTokenUser{
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		AvatarURL: dbUser.AvatarURL,
		ID:        dbUser.ID,
		Role:      dbUser.Role,
		Verified:  dbUser.Verified,
		Banned:    dbUser.Banned,
	}

	accessToken, _, err := generateTokens(refreshTokenUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
