package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- ValidateFileExtension tests ---

func TestValidateFileExtensionValid(t *testing.T) {
	err := ValidateFileExtension("pdf")
	assert.NoError(t, err)
}

func TestValidateFileExtensionInvalid(t *testing.T) {
	err := ValidateFileExtension("docx")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PDF")
}

// --- ValidateFileCategory tests ---

func TestValidateFileCategoryJobSeekerResume(t *testing.T) {
	err := ValidateFileCategory(CategoryResume, RoleJobSeeker)
	assert.NoError(t, err)
}

func TestValidateFileCategoryJobSeekerTranscript(t *testing.T) {
	err := ValidateFileCategory(CategoryTranscript, RoleJobSeeker)
	assert.NoError(t, err)
}

func TestValidateFileCategoryJobSeekerCertification(t *testing.T) {
	err := ValidateFileCategory(CategoryCertification, RoleJobSeeker)
	assert.NoError(t, err)
}

func TestValidateFileCategoryJobSeekerVerification(t *testing.T) {
	err := ValidateFileCategory(CategoryVerification, RoleJobSeeker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job seekers can only upload")
}

func TestValidateFileCategoryCompanyVerification(t *testing.T) {
	err := ValidateFileCategory(CategoryVerification, RoleCompany)
	assert.NoError(t, err)
}

func TestValidateFileCategoryCompanyCertification(t *testing.T) {
	err := ValidateFileCategory(CategoryCertification, RoleCompany)
	assert.NoError(t, err)
}

func TestValidateFileCategoryCompanyResume(t *testing.T) {
	err := ValidateFileCategory(CategoryResume, RoleCompany)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "companies can only upload")
}

func TestValidateFileCategoryCompanyTranscript(t *testing.T) {
	err := ValidateFileCategory(CategoryTranscript, RoleCompany)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "companies can only upload")
}

func TestValidateFileCategoryInvalidRole(t *testing.T) {
	err := ValidateFileCategory(CategoryResume, "invalid_role")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user role")
}

// --- ValidateFile tests ---

func TestValidateFileSuccess(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "application/pdf",
		FileExtension: "pdf",
		Category:      CategoryResume,
	}

	err := ValidateFile(file, RoleJobSeeker)
	assert.NoError(t, err)
}

func TestValidateFileFileTooLarge(t *testing.T) {
	file := &File{
		Size:          11 * 1024 * 1024, // 11MB
		ContentType:   "application/pdf",
		FileExtension: "pdf",
		Category:      CategoryResume,
	}

	err := ValidateFile(file, RoleJobSeeker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum")
}

func TestValidateFileInvalidContentType(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "image/png",
		FileExtension: "pdf",
		Category:      CategoryResume,
	}

	err := ValidateFile(file, RoleJobSeeker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PDF")
}

func TestValidateFileInvalidExtension(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "application/pdf",
		FileExtension: "docx",
		Category:      CategoryResume,
	}

	err := ValidateFile(file, RoleJobSeeker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PDF")
}

func TestValidateFileInvalidCategoryForRole(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "application/pdf",
		FileExtension: "pdf",
		Category:      CategoryVerification,
	}

	err := ValidateFile(file, RoleJobSeeker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "job seekers can only upload")
}

func TestValidateFileCompanySuccess(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "application/pdf",
		FileExtension: "pdf",
		Category:      CategoryVerification,
	}

	err := ValidateFile(file, RoleCompany)
	assert.NoError(t, err)
}

func TestValidateFileCompanyInvalidCategory(t *testing.T) {
	file := &File{
		Size:          1024,
		ContentType:   "application/pdf",
		FileExtension: "pdf",
		Category:      CategoryResume,
	}

	err := ValidateFile(file, RoleCompany)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "companies can only upload")
}