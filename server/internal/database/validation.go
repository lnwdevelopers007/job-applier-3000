// Package database handles all connection between web server and mongodb.

package database

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// validateContent validates contents from HTTP request methods.
func ValidateRequestContent[Schema any](c *gin.Context) (bson.M, error) {
	var raw bson.M
	// decode context to bson map (dict)
	if err := c.ShouldBindJSON(&raw); err != nil {
		return nil, err
	}

	// validate whether the structure is correct
	if err := validateSchema[Schema](raw); err != nil {
		return nil, err
	}

	// Add an _id if missing
	if _, ok := raw["_id"]; !ok {
		raw["_id"] = primitive.NewObjectID()
	}
	return raw, nil
}

// validateSchema returns error if the raw data doesn't conform to the schema.
// Returns error when there're fields failed the validation criteria in the schema
// ex. validate:"gte:0" (that specific field must have value >= 0).
// Returns error when there're fields have different data type from the schema.
// Returns nil otherwies (aka. validation is successful).
// Validation will be successful even if there exist field(s) that are not specified
// in the schema. (we're using NoSQL duh).
func validateSchema[Schema any](raw bson.M) error {
	// Extract required fields into struct
	var s Schema
	b, _ := bson.Marshal(raw) // convert map -> bson
	bson.Unmarshal(b, &s)     // map -> struct

	err := validate.Struct(s)
	return err
}
