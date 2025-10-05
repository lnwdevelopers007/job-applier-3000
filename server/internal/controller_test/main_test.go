// package controller_test contains test files for controllers.
package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var globalApplicantID string

// TestMain is the main test function of every packageName_test package.
// more info: https://pkg.go.dev/testing#hdr-Main
func TestMain(m *testing.M) {
	db := database.GetDatabase()

	collections := []string{
		"job_applications",
		"jobs",
		"users",
	}

	createMockCollections(db, collections)
	initJobApplication()
	// m.Run() run all other *_test.go files in controller_test.
	// Lines before this is equivalent to python's setUp
	m.Run()
	// and below this is equivalent to tearDown.
}

// initJobApplication initialize job application for retrieval tests.
func initJobApplication() {
	router := getTestRouter()
	w1 := httptest.NewRecorder()

	body, _ := json.Marshal(rawUser("tralalero tralala"))

	req, _ := http.NewRequest("POST", "/users/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req)
	r := regexp.MustCompile(`"InsertedID":"(.+)"`)
	matches := r.FindStringSubmatch(w1.Body.String())
	id := matches[1]

	w2 := httptest.NewRecorder()
	objID, _ := primitive.ObjectIDFromHex(id)
	body, _ = json.Marshal(rawJobApplication(objID))
	req, _ = http.NewRequest("POST", "/apply/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req)
	globalApplicantID = id

}

// createMockCollections ensures that each collection in `names` exists
// and wipes all data from them.
//
// IMPORTANT: PLEASE USE DIFFERENT DB_NAME IN .env FILE. Otherwise, ALL YOUR
// DATA IN REAL PRODUCTION DATABASE WILL BE LOST WHEN RUNNING CONTROLLER TESTS.
func createMockCollections(db *mongo.Database, names []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get existing collections
	existing, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	existingMap := make(map[string]bool)
	for _, name := range existing {
		existingMap[name] = true
	}

	for _, name := range names {
		// Create collection if it doesn't exist
		if !existingMap[name] {
			if err := db.CreateCollection(ctx, name); err != nil {
				log.Printf("[WARN] Failed to create collection %q: %v\n", name, err)
				continue
			}
			log.Printf("[INFO] Created collection: %s\n", name)
		}

		// Wipe data from collection
		// TODO: not working properly, data not deleted at all! someone pls fix this.
		// I suspected that the issue is database singleton not singletonning.
		coll := db.Collection(name)
		if _, err := coll.DeleteMany(ctx, bson.D{}); err != nil {
			log.Printf("[WARN] Failed to clear collection %q: %v\n", name, err)
			continue
		}
		log.Printf("[INFO] Cleared collection: %s\n", name)
	}

	return nil
}

func rawJobApplication(applicantID primitive.ObjectID) map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		"applicantID": applicantID,
		"jobID":       primitive.NewObjectID(),
		"status":      "waiting for approval",
		"createdAt":   now,
	}
}

// TODO: add more fields in here if we decided to enforce required fields in user schema.
func rawUser(name string) map[string]any {
	return map[string]any{
		"name": name,
	}
}

func rawJob(title string) map[string]any {
	now := time.Now().UTC().Truncate(time.Second)
	return map[string]any{
		// basic info
		"title":           title,
		"companyID":       primitive.NewObjectID().Hex(),
		"location":        "Kivotos",
		"workType":        "onsite",
		"workArrangement": "full-time",
		"currency":        "THB",
		"minSalary":       2000.34,
		"maxSalary":       300000.213213,

		// description
		"jobDescription": "long",
		"jobSummary":     "longer",

		// requirements
		"requiredSkills":  "none",
		"experienceLevel": "a lot",
		"education":       "maybe",
		"niceToHave":      "noting",

		// post settings
		"postOpenDate":        now,
		"applicationDeadline": now,
		"numberOfPositions":   1,
		"visibility":          "public",
		"emailNotifications":  true,
		"autoReject":          false,
	}
}

// getTestRouter gets the router configured for running controller tests.
func getTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := controller.NewRouter()
	return router
}
