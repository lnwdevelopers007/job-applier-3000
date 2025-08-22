// Package database handles all connection between web server and mongodb.

package database

import (
	"github.com/gin-gonic/gin"
)

// validateContent validates contents from HTTP request methods.
func ValidateRequestContent[Schema any](c *gin.Context) (any, error) {
	var json Schema
	// decode context to bson map (dict)
	if err := c.ShouldBindJSON(&json); err != nil {
		return nil, err
	}

	return json, nil
}
