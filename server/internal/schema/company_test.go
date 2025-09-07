package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---
func companyValidPayload() map[string]any {
	return map[string]any{
		"id":          primitive.NewObjectID(),
		"name":        "Millenium Science School co. ltd.",
		"userID":      primitive.NewObjectID(),
		"aboutUs":     "Brr Brr patapim",
		"companyType": "cute and funny",
	}
}

func bindMockCompany(t *testing.T, payload map[string]any) (Company, error) {
	return bindMockRequest[Company](t, payload)
}

func TestValidCompany(t *testing.T) {
	_, err := bindMockCompany(t, companyValidPayload())
	assert.NoError(t, err)
}

func TestCompanyWrongNameDataType(t *testing.T) {
	payload := companyValidPayload()
	payload["name"] = 123213213
	_, err := bindMockCompany(t, payload)
	assert.Error(t, err)
}

func TestCompanyMissingName(t *testing.T) {
	payload := companyValidPayload()
	delete(payload, "name")
	_, err := bindMockCompany(t, payload)
	assert.Error(t, err)
}

func TestCompanyMissingAboutUs(t *testing.T) {
	p := companyValidPayload()
	delete(p, "aboutUs")
	_, err := bindMockCompany(t, p)
	assert.Error(t, err)
}

func TestCompanyWrongAboutUsDataType(t *testing.T) {
	payload := companyValidPayload()
	payload["aboutUs"] = 123
	_, err := bindMockCompany(t, payload)
	assert.Error(t, err)
}

func TestCompanyMissingCompanyType(t *testing.T) {
	p := companyValidPayload()
	delete(p, "companyType")
	_, err := bindMockCompany(t, p)
	assert.Error(t, err)
}

func TestCompanyWrongCompanyDataType(t *testing.T) {
	payload := companyValidPayload()
	payload["companyType"] = 123
	_, err := bindMockCompany(t, payload)
	assert.Error(t, err)
}

func TestCompanyMissingUserID(t *testing.T) {
	payload := companyValidPayload()
	delete(payload, "userID")
	_, err := bindMockCompany(t, payload)
	assert.Error(t, err)
}
