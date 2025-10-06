package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
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
