package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lnwdevelopers007/job-applier-3000/server/config"
)

var jwtSecret = []byte(config.LoadEnv("JWT_SECRET"))

// AuthMiddleware validates JWT token and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Check for Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Set user information in context
		userID, _ := claims["userID"].(string)
		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)

		c.Set("userID", userID)
		c.Set("role", role)
		c.Set("email", email)
		c.Set("name", name)
		
		// Note: 'verified' is managed by admin, not passed in middleware
		// If you need to check verification status, query the database

		c.Next()
	}
}

// OptionalAuthMiddleware makes auth optional based on environment variable
func OptionalAuthMiddleware() gin.HandlerFunc {
	enableAuth := config.LoadBoolean("ENABLE_AUTH")
	
	if enableAuth {
		return AuthMiddleware()
	}
	
	// If auth is disabled, just pass through
	return func(c *gin.Context) {
		c.Next()
	}
}