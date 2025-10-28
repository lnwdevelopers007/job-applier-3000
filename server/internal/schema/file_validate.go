package schema

import (
	"fmt"
	"mime/multipart"
)

// Validation constants
const (
	MaxFileSize     = 10 * 1024 * 1024 // 10MB in bytes
	AllowedMimeType = "application/pdf"
)

// User roles
const (
	RoleJobSeeker = "jobSeeker"
	RoleCompany   = "company"
)

// ValidateUploadedFile validates a file from multipart form before processing
func ValidateUploadedFile(header *multipart.FileHeader) error {
	// Validate file size
	if header.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if contentType != AllowedMimeType {
		return fmt.Errorf("only PDF files are allowed")
	}

	return nil
}

// ValidateFileCategory validates if the category is valid for the given user role
func ValidateFileCategory(category FileCategory, userRole string) error {
	jobSeekerCategories := map[FileCategory]bool{
		CategoryResume:        true,
		CategoryTranscript:   true,
		CategoryCertification: true,
	}

	companyCategories := map[FileCategory]bool{
		CategoryVerification:  true,
		CategoryCertification: true,
	}

	switch userRole {
	case RoleJobSeeker:
		if !jobSeekerCategories[category] {
			return fmt.Errorf("job seekers can only upload resume, transcript, or certification files")
		}
	case RoleCompany:
		if !companyCategories[category] {
			return fmt.Errorf("companies can upload verification files only")
		}
	default:
		return fmt.Errorf("invalid user role")
	}

	return nil
}

// ValidateFileExtension validates if the file extension is allowed
func ValidateFileExtension(extension string) error {
	if extension != "pdf" {
		return fmt.Errorf("only PDF file extension is allowed")
	}
	return nil
}

// ValidateFile performs comprehensive validation on a File struct
func ValidateFile(file *File, userRole string) error {
	// Validate file size
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}

	// Validate content type
	if file.ContentType != AllowedMimeType {
		return fmt.Errorf("only PDF files are allowed")
	}

	// Validate file extension
	if err := ValidateFileExtension(file.FileExtension); err != nil {
		return err
	}

	// Validate category for user role
	if err := ValidateFileCategory(file.Category, userRole); err != nil {
		return err
	}

	return nil
}