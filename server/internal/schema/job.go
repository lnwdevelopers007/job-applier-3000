// Package database handles all connection between web server and mongodb.
package schema

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Schema for various datatypes.
// more info here: https://gin-gonic.com/en/docs/examples/binding-and-validation/

type Job struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title           string             `bson:"title" json:"title" binding:"required"`
	CompanyID       primitive.ObjectID `bson:"companyID" json:"companyID" binding:"required"`
	Location        string             `bson:"location" json:"location" binding:"required"`
	Salary          int                `bson:"salary" json:"salary" binding:"required,gte=0"`
	SalaryRate      string             `bson:"salaryRate" json:"salaryRate" binding:"required"`
	WorkType        string             `bson:"workType" json:"workType" binding:"required"`
	ContractType    string             `bson:"contractType" json:"contractType" binding:"required"`
	PrivacyPolicy   string             `bson:"privacyPolicy" json:"privacyPolicy,omitempty"`
	PublicationInfo Publication        `bson:"publicationInfo" json:"publicationInfo" binding:"required"`
	Criteria        JobCriteria        `bson:"criteria" json:"criteria" binding:"required" `
	IsApproved      bool               `bson:"isApproved" json:"isApproved" binding:"required"`
}

type Publication struct {
	IsHiring  bool      `bson:"isHiring" json:"isHiring" binding:"required"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" binding:"required"`
	StartDate time.Time `bson:"startDate" json:"startDate" binding:"required"`
	EndDate   time.Time `bson:"endDate" json:"endDate" binding:"required"`
}

type JobCriteria struct {
	Requirements    []string `bson:"requirements" json:"requirements" binding:"required,min=1"`
	Qualifications  []string `bson:"qualifications" json:"qualifications" binding:"required"`
	CommonQuestions []string `bson:"commonQuestions" json:"commonQuestions,omitempty"`
}
