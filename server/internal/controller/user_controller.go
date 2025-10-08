package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	baseController BaseController[schema.User]
}

func NewUserController() UserController {
	return UserController{
		baseController: BaseController[schema.User]{
			collectionName: "users",
			displayName:    "User",
		},
	}
}

func (jc UserController) Query(c *gin.Context) {
	userFilter, shouldReturn := userFilter(c)
	if shouldReturn {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := findAll[schema.User](ctx, jc.baseController.collectionName, userFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func userFilter(c *gin.Context) (bson.M, bool) {
	allowedParams := map[string]func(string) (any, error){
		"id": func(v string) (any, error) {
			if v == "" {
				return nil, fmt.Errorf("id parameter is empty")
			}
			return primitive.ObjectIDFromHex(v)
		},
		"role": func(v string) (any, error) {
			if v == "" {
				return nil, nil
			}
			return bson.M{"$eq": v}, nil
		},
	}

	filter := bson.M{}

	// Loop through query params
	for key, value := range c.Request.URL.Query() {
		if fn, ok := allowedParams[key]; ok {
			val, err := fn(value[0])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return nil, true
			}
			if val != nil {
				if key == "id" {
					filter["_id"] = val
				} else {
					filter[key] = val
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported query parameter: " + key})
			return nil, true
		}
	}
	return filter, false
}

// Update updates an existing user by ID.
func (jc UserController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Create creates a user.
func (jc UserController) Create(c *gin.Context) {
	jc.baseController.Create(c)
}

// Delete removes a user by ID.
func (jc UserController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}

// RetrieveAll fetches all companies from the database.
func (jc UserController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// RetrieveOne fetches a single user by ID.
func (jc UserController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}
