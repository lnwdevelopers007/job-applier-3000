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
		// basic info
		"title":           "Brr Brr Engineer",
		"companyID":       primitive.NewObjectID(),
		"location":        "Millenium Science School, Kivotos.",
		"workType":        "onsite",
		"workArrangement": "full-time",
		"currency":        "THB",
		"minSalary":       2000.34,
		"maxSalary":       300000.213213,

	// description
	"jobDescription": "This is a detailed job description that is intentionally long to meet the minimum length requirement for tests.",
	"jobSummary":     "Concise summary with sufficient length",

	// requirements
	"requiredSkills":  "skill1, skill2, skill3",
		"experienceLevel": "a lot",
		"education":       "maybe",
		"niceToHave":      "noting",

		// post settings
		"postOpenDate": now,
		"applicationDeadline": now,
		"numberOfPositions":   1,
		"visibility":          "public",
		"emailNotifications":  true,
		"autoReject":          false,
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

func TestJobNegativeMinSalary(t *testing.T) {
	payload := jobValidPayload()
	payload["minSalary"] = -1 // gte=0
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobMissingMinSalary(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "minSalary") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobNegativeMaxSalary(t *testing.T) {
	payload := jobValidPayload()
	payload["minSalary"] = 0  // gte=0
	payload["maxSalary"] = -1 // gte=0
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}

func TestJobMissingMaxSalary(t *testing.T) {
	payload := jobValidPayload()
	delete(payload, "maxSalary") // required
	_, err := bindMockJob(t, payload)
	assert.Error(t, err)
}
