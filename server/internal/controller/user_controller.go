package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

type CompanyController struct {
	baseController BaseController[schema.User]
}

func NewCompanyController() CompanyController {
	return CompanyController{
		baseController: BaseController[schema.User]{
			collectionName: "users",
			displayName:    "User",
		},
	}
}

// Update updates an existing company by ID.
func (jc CompanyController) Update(c *gin.Context) {
	jc.baseController.Update(c)
}

// Delete removes a company by ID.
func (jc CompanyController) Delete(c *gin.Context) {
	jc.baseController.Delete(c)
}

// RetrieveAll fetches all companies from the database.
func (jc CompanyController) RetrieveAll(c *gin.Context) {
	jc.baseController.RetrieveAll(c)
}

// RetrieveOne fetches a single company by ID.
func (jc CompanyController) RetrieveOne(c *gin.Context) {
	jc.baseController.RetrieveOne(c)
}
