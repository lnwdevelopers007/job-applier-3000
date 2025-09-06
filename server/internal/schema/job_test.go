package schema

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func validPayload() map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"title":        "Brr Brr Engineer",
		"companyID":    primitive.NewObjectID(),
		"location":     "Millenium Science School, Kivotos.",
		"salary":       120000,
		"salaryRate":   "yearly",
		"workType":     "onsite",
		"contractType": "full-time",
		// privacyPolicy is optional
		"publicationInfo": map[string]any{
			"isHiring":  true,
			"createdAt": now.Format(time.RFC3339),
			"startDate": now.Format(time.RFC3339),
			"endDate":   now.AddDate(0, 1, 0).Format(time.RFC3339),
		},
		"criteria": map[string]any{
			"requirements":   []string{"Go", "MongoDB"},
			"qualifications": []string{"Bachelor's Degree"},
			// commonQuestions is optional
		},
		"isApproved": true, // NOTE: with `binding:"required"` on a bool, false will fail validation
	}
}

func bindJobMock(t *testing.T, payload map[string]any) (Job, error) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("POST", "/jobs", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	var j Job
	err = c.ShouldBindJSON(&j)
	return j, err
}

// --- valid case ---

func TestValidJob(t *testing.T) {
	payload := validPayload()
	_, err := bindJobMock(t, payload)
	assert.NoError(t, err)
}

// --- Title ---

func TestValidTitle(t *testing.T) {
	payload := validPayload()
	payload["title"] = "Software Engineer"
	_, err := bindJobMock(t, payload)
	assert.NoError(t, err)
}

func TestInvalidTitle(t *testing.T) {
	payload := validPayload()
	delete(payload, "title") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- Company ---

func TestInvalidCompany(t *testing.T) {
	payload := validPayload()
	delete(payload, "company") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- Location ---

func TestInvalidLocation(t *testing.T) {
	payload := validPayload()
	delete(payload, "location") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- Salary ---

func TestValidSalary(t *testing.T) {
	payload := validPayload()
	payload["salary"] = 50000
	_, err := bindJobMock(t, payload)
	assert.NoError(t, err)
}

func TestInvalidSalaryNegative(t *testing.T) {
	payload := validPayload()
	payload["salary"] = -1 // gte=0
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidSalaryMissing(t *testing.T) {
	payload := validPayload()
	delete(payload, "salary") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- SalaryRate / WorkType / ContractType ---

func TestInvalidSalaryRateMissing(t *testing.T) {
	payload := validPayload()
	delete(payload, "salaryRate") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidWorkTypeMissing(t *testing.T) {
	payload := validPayload()
	delete(payload, "workType") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidContractTypeMissing(t *testing.T) {
	payload := validPayload()
	delete(payload, "contractType") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- PrivacyPolicy (optional) ---

func TestPrivacyPolicyOmittedIsValid(t *testing.T) {
	payload := validPayload()
	delete(payload, "privacyPolicy") // optional
	_, err := bindJobMock(t, payload)
	assert.NoError(t, err)
}

// --- Publication (nested) ---

func TestInvalidPublicationMissingStartDate(t *testing.T) {
	payload := validPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "startDate") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidPublicationMissingEndDate(t *testing.T) {
	payload := validPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "endDate") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidPublicationMissingCreatedAt(t *testing.T) {
	payload := validPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "createdAt") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidPublicationMissingIsHiring(t *testing.T) {
	payload := validPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "isHiring") // required
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- Criteria (nested) ---

func TestInvalidCriteriaMissingRequirements(t *testing.T) {
	payload := validPayload()
	crit := payload["criteria"].(map[string]any)
	delete(crit, "requirements") // required slice
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidCriteriaEmptyRequirements(t *testing.T) {
	payload := validPayload()
	crit := payload["criteria"].(map[string]any)
	crit["requirements"] = []string{} // required slice cannot be empty
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestInvalidCriteriaMissingQualifications(t *testing.T) {
	payload := validPayload()
	crit := payload["criteria"].(map[string]any)
	delete(crit, "qualifications") // required slice
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

// --- IsApproved (bool with `required`) ---

func TestInvalidIsApprovedFalse(t *testing.T) {
	payload := validPayload()
	payload["isApproved"] = false // `required` on bool -> false is invalid
	_, err := bindJobMock(t, payload)
	assert.Error(t, err)
}

func TestValidIsApprovedTrue(t *testing.T) {
	payload := validPayload()
	payload["isApproved"] = true
	_, err := bindJobMock(t, payload)
	assert.NoError(t, err)
}
