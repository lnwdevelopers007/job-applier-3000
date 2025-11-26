package middleware

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

var methodLimits = map[string]uint{
	http.MethodGet:    2500,
	http.MethodPost:   20,
	http.MethodPut:		25,
	http.MethodPatch:  25,
	http.MethodDelete: 25,
}

var rateLimiters map[string]gin.HandlerFunc

func init() {
	rateLimiters = make(map[string]gin.HandlerFunc, len(methodLimits))
	for method, limit := range methodLimits {
		rateLimiters[method] = newRateLimiter(limit)
	}
}

func RateLimiterMiddleware(method string) gin.HandlerFunc {
	if mw, ok := rateLimiters[method]; ok {
		return mw
	}

	return func(c *gin.Context) { c.Next() }
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func rateLimitErrorHandler(c *gin.Context, info ratelimit.Info) {
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"error": "too many requests, please try again in " + time.Until(info.ResetTime).String(),
	})
}

func newRateLimiter(limit uint) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: limit,
	})

	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: rateLimitErrorHandler,
		KeyFunc:      keyFunc,
	})

	return mw
}