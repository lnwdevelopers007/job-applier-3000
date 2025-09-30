package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

type CompanyController struct {
	baseController BaseController[schema.Company]
}

func NewCompanyController() CompanyController {
	return CompanyController{
		baseController: BaseController[schema.Company]{
			collectionName: "companies",
			displayName:    "Company",
		},
	}
}

// Create inserts a new company into the database.
func (jc CompanyController) Create(c *gin.Context) {
	jc.baseController.Create(c)
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
