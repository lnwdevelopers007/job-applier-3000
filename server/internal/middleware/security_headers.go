package middleware

import (
	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevents browsers from trying to guess ("sniff") the MIME type
		// Stops malicious HTML from being executed when expected to be plain text
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevents your site from being embedded in iframes (clickjacking protection)
		// DENY = no framing allowed, SAMEORIGIN = only same domain can frame
		c.Header("X-Frame-Options", "DENY")
		
		// Legacy XSS protection for older browsers
		// Modern browsers have this built-in, but doesn't hurt to include
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Forces HTTPS connections for your domain
		// Browser will remember for 1 year (31536000 seconds) to always use HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy - controls what resources can be loaded
		// Less restrictive for images and media to allow external sources
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net https://unpkg.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src * data: blob:; media-src *; font-src 'self' data: https://fonts.gstatic.com; connect-src 'self' https://api.github.com https://api.google.com")
		
		// Controls how much referrer info is sent with requests
		// strict-origin-when-cross-origin = full URL for same-origin, only origin for cross-origin
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Disables browser features that could be privacy/security risks
		// Empty () means feature is disabled for all origins
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")
		
		c.Next()
	}
}