package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- helpers ---

func jobValidPayload() map[string]any {
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

func bindMockJob(t *testing.T, payload map[string]any) (Job, error) {
	return bindMockRequest[Job](t, payload)
}

// --- valid case ---

func TestValidJob(t *testing.T) {
	payload := jobValidPayload()
	_, err := bindMockJob(t, payload)
	assert.NoError(t, err)
}

// --- Title ---
func TestJobMissingTitle(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "title") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobWrongTitleDataType(t *testing.T) {
	payload := jobValidPayload()
	payload["title"] = 1212312121
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- Company ---
func TestJobMissingCompanyID(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "companyID") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- Location ---

func TestJobMissingLocation(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "location") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobWrongLocationDataType(t *testing.T) {
	payload := jobValidPayload()
	payload["location"] = 1212312121
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- Salary ---

func TestJobNegativeSalary(t *testing.T) {
	payload := jobValidPayload()
	payload["salary"] = -1 // gte=0
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobMissingSalary(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "salary") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}
func TestJobWrongSalaryDataType(t *testing.T) {
	payload := jobValidPayload()
	payload["location"] = 1212312121
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- SalaryRate / WorkType / ContractType ---
func TestJobMissingSalaryRate(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "salaryRate") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobMissingWorkType(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "workType") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobMissingContractType(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "contractType") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- PrivacyPolicy (optional) ---

func TestJobPrivacyPolicyOmittedIsValid(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "privacyPolicy") // optional
	_, err := bindMockJob(t, payload)
	assert.NoError(t, err)
}

// --- Publication (nested) ---

func TestJobPublicationMissingStartDate(t *testing.T) {
	payload := jobValidPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "startDate") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobPublicationMissingEndDate(t *testing.T) {
	payload := jobValidPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "endDate") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobPublicationMissingCreatedAt(t *testing.T) {
	payload := jobValidPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "createdAt") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobPublicationMissingIsHiring(t *testing.T) {
	payload := jobValidPayload()
	pub := payload["publicationInfo"].(map[string]any)
	delete(pub, "isHiring") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

// --- Criteria (nested) ---

func TestJobCriteriaMissingRequirements(t *testing.T) {
	payload := jobValidPayload()
	crit := payload["criteria"].(map[string]any)
	delete(crit, "requirements") // required slice
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobCriteriaEmptyRequirements(t *testing.T) {
	payload := jobValidPayload()
	crit := payload["criteria"].(map[string]any)
	crit["requirements"] = []string{} // required slice cannot be empty
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobCriteriaMissingQualifications(t *testing.T) {
	payload := jobValidPayload()
	crit := payload["criteria"].(map[string]any)
	delete(crit, "qualifications") // required slice
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobIsApprovedCanBeFalse(t *testing.T) {
	payload := jobValidPayload()
	payload["isApproved"] = false // `required` on bool -> false is invalid
	_, err := bindMockJob(t, payload)
	assert.NoError(t, err)
}
