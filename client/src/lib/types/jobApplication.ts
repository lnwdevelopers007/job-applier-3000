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

export interface JobApplicationFilters {
  applicantID?: string;
  jobID?: string;
  companyID?: string;
  status?: string;
}