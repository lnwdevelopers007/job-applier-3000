/**
 * Job Application-related type definitions
 */

export interface JobApplication {
  id: string;
  applicantID: string;
  jobID: string;
  companyID?: string;
  status: string;
  createdAt: string;
  // Add other application fields as needed
}

// Type for the nested response structure that includes user data
export interface JobApplicationWithApplicant {
  jobApplication: JobApplication;
  applicant: {
    id: string;
    userID: string;
    provider: string;
    email: string;
    name: string;
    avatarURL?: string;
    role: string;
    verified: boolean;
    updatedAt: string;
    createdAt: string;
    userInfo: Array<{
      Key: string;
      Value: string;
    }>;
  };
}

export interface JobApplicationFilters {
  applicantID?: string;
  jobID?: string;
  companyID?: string;
  status?: string;
}