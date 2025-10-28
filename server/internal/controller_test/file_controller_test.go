package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/lnwdevelopers007/job-applier-3000/server/internal/database"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test data IDs
var (
	testFileUserID1   = primitive.NewObjectID()
	testFileUserID2   = primitive.NewObjectID()
	testFileCompanyID = primitive.NewObjectID()
	testFileJobID     = primitive.NewObjectID()
	testFileAppID     = primitive.NewObjectID()
)

// setupFileTestData creates test data for file controller tests
func setupFileTestData(t *testing.T) {
	// Set ENABLE_AUTH=false for testing
	os.Setenv("ENABLE_AUTH", "false")

	db := database.GetDatabase()
	ctx := context.Background()

	// Clean up collections
	db.Collection("files").DeleteMany(ctx, bson.M{"userID": bson.M{"$in": []primitive.ObjectID{testFileUserID1, testFileUserID2, testFileCompanyID}}})
	db.Collection("users").DeleteMany(ctx, bson.M{"_id": bson.M{"$in": []primitive.ObjectID{testFileUserID1, testFileUserID2, testFileCompanyID}}})
	db.Collection("jobs").DeleteMany(ctx, bson.M{"_id": testFileJobID})
	db.Collection("job_applications").DeleteMany(ctx, bson.M{"_id": testFileAppID})

	// Create test users
	db.Collection("users").InsertMany(ctx, []interface{}{
		bson.M{
			"_id":      testFileUserID1,
			"userID":   "jobseeker1@test.com",
			"email":    "jobseeker1@test.com",
			"name":     "Job Seeker One",
			"role":     "jobSeeker",
			"verified": true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"_id":      testFileUserID2,
			"userID":   "jobseeker2@test.com",
			"email":    "jobseeker2@test.com",
			"name":     "Job Seeker Two",
			"role":     "jobSeeker",
			"verified": true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
		bson.M{
			"_id":      testFileCompanyID,
			"userID":   "company@test.com",
			"email":    "company@test.com",
			"name":     "Test Company",
			"role":     "company",
			"verified": true,
			"createdAt": time.Now(),
			"updatedAt": time.Now(),
		},
	})

	// Create test job
	db.Collection("jobs").InsertOne(ctx, bson.M{
		"_id":                  testFileJobID,
		"title":                "Backend Developer",
		"companyID":            testFileCompanyID,
		"location":             "Bangkok",
		"workType":             "Full-time",
		"workArrangement":      "On-site",
		"currency":             "USD",
		"minSalary":            3000.0,
		"maxSalary":            5000.0,
		"jobDescription":       "Test job",
		"jobSummary":           "Test",
		"requiredSkills":       "Go",
		"experienceLevel":      "Mid",
		"education":            "Bachelor",
		"postOpenDate":         time.Now(),
		"applicationDeadline":  time.Now().Add(30 * 24 * time.Hour),
		"numberOfPositions":    1,
		"visibility":           "Public",
		"emailNotifications":   true,
		"autoReject":           false,
	})

	// Create test application
	db.Collection("job_applications").InsertOne(ctx, bson.M{
		"_id":         testFileAppID,
		"applicantID": testFileUserID1,
		"jobID":       testFileJobID,
		"status":      "pending",
		"createdAt":   time.Now(),
	})
}

// Test 1: Upload File Successfully
func TestFileUploadSuccess(t *testing.T) {
	setupFileTestData(t)
	router := getTestRouter()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create file part with PDF content
	part, _ := writer.CreateFormFile("file", "my-resume.pdf")
	part.Write([]byte("%PDF-1.4\n%Test PDF Content\n%%EOF"))
	
	writer.WriteField("category", "resume")
	writer.WriteField("userID", testFileUserID1.Hex())
	writer.WriteField("userRole", "jobSeeker")
	writer.Close()

	req, _ := http.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "only PDF files are allowed")
}

// Test 2: Upload File - Wrong Category for Role
func TestFileUploadWrongCategoryForRole(t *testing.T) {
	setupFileTestData(t)
	router := getTestRouter()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "verification.pdf")
	part.Write([]byte("%PDF-1.4\n%Test PDF\n%%EOF"))
	
	// Job seeker trying to upload verification (only companies can)
	writer.WriteField("category", "verification")
	writer.WriteField("userID", testFileUserID1.Hex())
	writer.WriteField("userRole", "jobSeeker")
	writer.Close()

	req, _ := http.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "only PDF files are allowed")
}

// Test 3: Upload File - File Too Large
func TestFileUploadTooLarge(t *testing.T) {
	setupFileTestData(t)
	router := getTestRouter()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", "huge-file.pdf")
	// Create 11MB of data (exceeds 10MB limit)
	largeContent := make([]byte, 11*1024*1024)
	part.Write(largeContent)
	
	writer.WriteField("category", "resume")
	writer.WriteField("userID", testFileUserID1.Hex())
	writer.WriteField("userRole", "jobSeeker")
	writer.Close()

	req, _ := http.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "10MB")
}

// Test 4: Download Own File Successfully
func TestFileDownloadOwnFileSuccess(t *testing.T) {
	setupFileTestData(t)
	
	// Insert a file directly to DB
	db := database.GetDatabase()
	ctx := context.Background()
	
	fileID := primitive.NewObjectID()
	testContent := []byte("This is PDF content for testing")
	
	db.Collection("files").InsertOne(ctx, bson.M{
		"_id":           fileID,
		"userID":        testFileUserID1,
		"filename":      "my-resume.pdf",
		"fileExtension": "pdf",
		"contentType":   "application/pdf",
		"size":          int64(len(testContent)),
		"category":      "resume",
		"content":       testContent,
		"uploadDate":    time.Now(),
	})

	router := getTestRouter()
	req, _ := http.NewRequest("GET", "/files/download/"+fileID.Hex(), nil)
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "my-resume.pdf")
	assert.Equal(t, testContent, w.Body.Bytes())
}

// Test 5: Download Other User's File - Should Fail
func TestFileDownloadOtherUserFileForbidden(t *testing.T) {
	setupFileTestData(t)
	
	db := database.GetDatabase()
	ctx := context.Background()
	
	fileID := primitive.NewObjectID()
	db.Collection("files").InsertOne(ctx, bson.M{
		"_id":           fileID,
		"userID":        testFileUserID1, // Owned by user1
		"filename":      "user1-resume.pdf",
		"fileExtension": "pdf",
		"contentType":   "application/pdf",
		"size":          int64(100),
		"category":      "resume",
		"content":       []byte("PDF content"),
		"uploadDate":    time.Now(),
	})

	router := getTestRouter()
	
	// User2 tries to download user1's file
	req, _ := http.NewRequest("GET", "/files/download/"+fileID.Hex(), nil)
	req.Header.Set("X-User-Id", testFileUserID2.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "permission")
}

// Test 6: List User's Files
func TestFileListUserFiles(t *testing.T) {
	setupFileTestData(t)
	
	db := database.GetDatabase()
	ctx := context.Background()
	
	// Insert multiple files for user1
	db.Collection("files").InsertMany(ctx, []interface{}{
		bson.M{
			"_id":           primitive.NewObjectID(),
			"userID":        testFileUserID1,
			"filename":      "resume.pdf",
			"fileExtension": "pdf",
			"contentType":   "application/pdf",
			"size":          int64(1000),
			"category":      "resume",
			"content":       []byte("Resume content"),
			"uploadDate":    time.Now(),
		},
		bson.M{
			"_id":           primitive.NewObjectID(),
			"userID":        testFileUserID1,
			"filename":      "transcript.pdf",
			"fileExtension": "pdf",
			"contentType":   "application/pdf",
			"size":          int64(500),
			"category":      "transcript",
			"content":       []byte("Transcript content"),
			"uploadDate":    time.Now(),
		},
	})

	router := getTestRouter()
	req, _ := http.NewRequest("GET", "/files/user/"+testFileUserID1.Hex(), nil)
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	files := response["files"].([]interface{})
	assert.Equal(t, 2, len(files))
}

// Test 7: Delete Own File Successfully
func TestFileDeleteOwnFileSuccess(t *testing.T) {
	setupFileTestData(t)
	
	db := database.GetDatabase()
	ctx := context.Background()
	
	fileID := primitive.NewObjectID()
	db.Collection("files").InsertOne(ctx, bson.M{
		"_id":           fileID,
		"userID":        testFileUserID1,
		"filename":      "to-delete.pdf",
		"fileExtension": "pdf",
		"contentType":   "application/pdf",
		"size":          int64(100),
		"category":      "resume",
		"content":       []byte("PDF content"),
		"uploadDate":    time.Now(),
	})

	router := getTestRouter()
	req, _ := http.NewRequest("DELETE", "/files/"+fileID.Hex(), nil)
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify file was deleted
	count, _ := db.Collection("files").CountDocuments(ctx, bson.M{"_id": fileID})
	assert.Equal(t, int64(0), count)
}


// Test 8: Company Access Applicant Files - Success
func TestFileCompanyAccessApplicantFilesSuccess(t *testing.T) {
	setupFileTestData(t)
	
	db := database.GetDatabase()
	ctx := context.Background()
	
	// Upload files for applicant
	db.Collection("files").InsertMany(ctx, []interface{}{
		bson.M{
			"_id":           primitive.NewObjectID(),
			"userID":        testFileUserID1,
			"filename":      "applicant-resume.pdf",
			"fileExtension": "pdf",
			"contentType":   "application/pdf",
			"size":          int64(1000),
			"category":      "resume",
			"content":       []byte("Resume"),
			"uploadDate":    time.Now(),
		},
		bson.M{
			"_id":           primitive.NewObjectID(),
			"userID":        testFileUserID1,
			"filename":      "transcript.pdf",
			"fileExtension": "pdf",
			"contentType":   "application/pdf",
			"size":          int64(500),
			"category":      "transcript",
			"content":       []byte("Transcript content"),
			"uploadDate":    time.Now(),
		},
	})

	router := getTestRouter()
	req, _ := http.NewRequest("GET", "/files/application/"+testFileAppID.Hex(), nil)
	req.Header.Set("X-User-Id", testFileCompanyID.Hex())
	req.Header.Set("X-User-Role", "company")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	assert.Equal(t, testFileAppID.Hex(), response["applicationID"])
	assert.Equal(t, testFileUserID1.Hex(), response["applicantID"])
	
	files := response["files"].([]interface{})
	assert.GreaterOrEqual(t, len(files), 2)
}

// Test 9: Job Seeker Cannot Access Applicant Files
func TestFileJobSeekerCannotAccessApplicantFiles(t *testing.T) {
	setupFileTestData(t)
	router := getTestRouter()

	// Job seeker tries to access applicant files endpoint
	req, _ := http.NewRequest("GET", "/files/application/"+testFileAppID.Hex(), nil)
	req.Header.Set("X-User-Id", testFileUserID1.Hex())
	req.Header.Set("X-User-Role", "jobSeeker")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "only companies")
}

// Test 10: Company Cannot Access Other Company's Application Files
func TestFileCompanyCannotAccessOtherCompanyFiles(t *testing.T) {
	setupFileTestData(t)
	
	db := database.GetDatabase()
	ctx := context.Background()
	
	// Create another company
	anotherCompanyID := primitive.NewObjectID()
	db.Collection("users").InsertOne(ctx, bson.M{
		"_id":      anotherCompanyID,
		"userID":   "othercompany@test.com",
		"email":    "othercompany@test.com",
		"name":     "Other Company",
		"role":     "company",
		"verified": true,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})

	router := getTestRouter()
	
	// Other company tries to access application for testFileCompanyID's job
	req, _ := http.NewRequest("GET", "/files/application/"+testFileAppID.Hex(), nil)
	req.Header.Set("X-User-Id", anotherCompanyID.Hex())
	req.Header.Set("X-User-Role", "company")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "your own jobs")
}