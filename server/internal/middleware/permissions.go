package middleware

// Permission defines the access control rules for a route
type Permission struct {
	AllowedRoles     []string
	RequireOwnership bool // For resources like "company can only edit own jobs"
}

// RoutePermissions defines RBAC rules for each endpoint
// Format: "METHOD:PATH_PATTERN" -> Permission
var RoutePermissions = map[string]Permission{
	// ===== Job Routes =====
	"POST:/jobs/": {
		AllowedRoles: []string{"company", "admin"},
	},
	"GET:/jobs/": {
		AllowedRoles: []string{"jobSeeker", "company", "faculty", "admin"},
	},
	"GET:/jobs/query": {
		AllowedRoles: []string{"jobSeeker", "company", "faculty", "admin"},
	},
	"GET:/jobs/:id": {
		AllowedRoles: []string{"jobSeeker", "company", "faculty", "admin"},
	},
	"PUT:/jobs/:id": {
		AllowedRoles:     []string{"company", "admin"},
		RequireOwnership: true, // Company can only update their own jobs
	},
	"DELETE:/jobs/:id": {
		AllowedRoles:     []string{"company", "admin"},
		RequireOwnership: true, // Company can only delete their own jobs
	},

	// ===== Job Application Routes =====
	"POST:/apply/": {
		AllowedRoles: []string{"jobSeeker", "admin"},
	},
	"GET:/apply/": {
		AllowedRoles: []string{"jobSeeker", "company", "admin"},
	},
	"GET:/apply/query": {
		AllowedRoles: []string{"jobSeeker", "company", "admin"},
	},
	"GET:/apply/:id": {
		AllowedRoles:     []string{"jobSeeker", "company", "admin"},
		RequireOwnership: true, // Users can only view their own applications
	},
	"PUT:/apply/:id": {
		AllowedRoles: []string{"company", "admin"}, // Only company/admin can update status
	},
	"DELETE:/apply/:id": {
		AllowedRoles:     []string{"jobSeeker", "admin"},
		RequireOwnership: true, // Job seekers can only delete their own applications
	},

	// ===== File Routes =====
	"POST:/files/upload": {
		AllowedRoles: []string{"jobSeeker", "company", "admin"},
	},
	"GET:/files/download/:id": {
		AllowedRoles:     []string{"jobSeeker", "company", "admin"},
		RequireOwnership: true, // Users can only download their own files
	},
	"GET:/files/user/:userId": {
		AllowedRoles:     []string{"jobSeeker", "company", "admin"},
		RequireOwnership: true, // Users can only list their own files
	},
	"DELETE:/files/:id": {
		AllowedRoles:     []string{"jobSeeker", "company", "admin"},
		RequireOwnership: true, // Users can only delete their own files
	},
	"GET:/files/application/:applicationId": {
		AllowedRoles: []string{"company", "admin"}, // Company can view applicant files for their jobs
	},
	"GET:/files/application/:applicationId/download/:fileId": {
		AllowedRoles: []string{"company", "admin"}, // Company can download applicant files
	},

	// ===== User Routes =====
	"GET:/users/": {
		AllowedRoles: []string{"admin"},
	},
	"GET:/users/query": {
		AllowedRoles: []string{"admin"},
	},
	"GET:/users/:id": {
		AllowedRoles: []string{"admin"},
	},
	"POST:/users/": {
		AllowedRoles: []string{"admin"},
	},
	"PUT:/users/:id": {
		AllowedRoles: []string{"admin"},
	},
	"DELETE:/users/:id": {
		AllowedRoles: []string{"admin"},
	},
}

// IsRoleAllowed checks if a role is allowed for a given permission
func (p Permission) IsRoleAllowed(role string) bool {
	for _, allowedRole := range p.AllowedRoles {
		if allowedRole == role {
			return true
		}
	}
	return false
}