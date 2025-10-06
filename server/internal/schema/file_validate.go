package schema

import "fmt"

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

// ValidateWithUserRole validates file based on user role
func (f *File) ValidateWithUserRole(userRole string) error {
	// Basic file validation
	if f.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}

	if f.ContentType != AllowedMimeType {
		return fmt.Errorf("only PDF files are allowed")
	}

	if f.FileExtension != "pdf" {
		return fmt.Errorf("only PDF file extension is allowed")
	}

	// Role-based category validation
	jobSeekerCategories := map[FileCategory]bool{
		CategoryResume:        true,
		CategoryCoverLetter:   true,
		CategoryCertification: true,
	}

	companyCategories := map[FileCategory]bool{
		CategoryVerification:  true,
		CategoryCertification: true, // Companies can also have certifications
	}

	switch userRole {
	case RoleJobSeeker:
		if !jobSeekerCategories[f.Category] {
			return fmt.Errorf("job seekers can only upload resume, cover_letter, or certification files")
		}
	case RoleCompany:
		if !companyCategories[f.Category] {
			return fmt.Errorf("companies can only upload verification or certification files")
		}
	default:
		return fmt.Errorf("invalid user role")
	}

	return nil
}

// Validate performs basic validation without role checking
func (f *File) Validate() error {
	if f.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 10MB")
	}

	if f.ContentType != AllowedMimeType {
		return fmt.Errorf("only PDF files are allowed")
	}

	if f.FileExtension != "pdf" {
		return fmt.Errorf("only PDF file extension is allowed")
	}

	validCategories := map[FileCategory]bool{
		CategoryResume:        true,
		CategoryCoverLetter:   true,
		CategoryCertification: true,
		CategoryVerification:  true,
	}

	if !validCategories[f.Category] {
		return fmt.Errorf("invalid file category")
	}

	return nil
}