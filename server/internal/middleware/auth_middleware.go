package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware validates JWT token and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		enableAuth, _ := strconv.ParseBool(os.Getenv("ENABLE_AUTH"))
		if !enableAuth {
			// Optional: allow test headers for user context
			if userID := c.GetHeader("X-User-Id"); userID != "" {
				c.Set("userID", userID)
			}
			if role := c.GetHeader("X-User-Role"); role != "" {
				c.Set("role", role)
			}
			c.Next()
			return
		}

		// 🔒 Enforce real authentication when enabled
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Extract claims
		userID, _ := claims["userID"].(string)
		role, _ := claims["role"].(string)
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)

		c.Set("userID", userID)
		c.Set("role", role)
		c.Set("email", email)
		c.Set("name", name)

		c.Next()
	}
}

