package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/schema"
)

// CompanyController wraps the BaseController with schema.Company
// and provides REST handlers for the "companies" collection.
type CompanyController struct {
	baseController BaseController[schema.Company]
}

// NewCompanyController creates a new CompanyController
// configured to operate on the "companies" collection.
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
