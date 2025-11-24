package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RejectUnknownOriginsMiddleware(allowed []string) gin.HandlerFunc {
	allowedMap := make(map[string]bool)
	for _, origin := range allowed {
		allowedMap[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Pass non browser requests
		if origin == "" {
			c.Next()
			return
		}

		// Reject non-allowed origins
		if !allowedMap[origin] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden origin",
			})
			return
		}

		c.Next()
	}
}