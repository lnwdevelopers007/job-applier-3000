package schema

import "fmt"

// Validation constants
const (
	MaxFileSize     = 10 * 1024 * 1024 // 10MB in bytes
	AllowedMimeType = "application/pdf"
)

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

	validParentColls := map[string]bool{
		"job_seekers": true,
		"companies":   true,
	}

	if !validParentColls[f.ParentColl] {
		return fmt.Errorf("invalid parent collection")
	}

	return nil
}
