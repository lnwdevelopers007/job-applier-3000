package auth

import (
	"sync"
	"time"
)

var (
	blacklist = make(map[string]time.Time)
	mu        sync.RWMutex
)

func init() {
	// Start cleanup routine to remove expired tokens every 5 minutes
	go cleanupExpiredTokens()
}

// Add a token to the blacklist with its expiration time
func AddToBlacklist(token string, expiresAt time.Time) {
	mu.Lock()
	defer mu.Unlock()
	blacklist[token] = expiresAt
}

// Check if a token is in the blacklist
func IsBlacklisted(token string) bool {
	mu.RLock()
	defer mu.RUnlock()
	
	expiresAt, exists := blacklist[token]
	if !exists {
		return false
	}
	
	// Check if token has expired (can be removed from blacklist)
	if time.Now().After(expiresAt) {
		return false // Token expired anyway, no need to block
	}
	
	return true
}

// Remove expired tokens from blacklist every 5 minutes
func cleanupExpiredTokens() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		mu.Lock()
		now := time.Now()
		for token, expiresAt := range blacklist {
			if now.After(expiresAt) {
				delete(blacklist, token)
			}
		}
		mu.Unlock()
	}
}