/**
 * Job-related type definitions
 */

// Backend Job schema (matches Go struct)
export interface Job {
  id: string;
  title: string;
  companyID: string;
  location: string;
  workType: string;
  workArrangement: string;
  currency: string;
  minSalary: number;
  maxSalary: number;
  jobDescription: string;
  jobSummary: string;
  requiredSkills: string;
  experienceLevel: string;
  education: string;
  niceToHave?: string;
  questions?: string;
  postOpenDate: string;
  applicationDeadline: string;
  numberOfPositions: number;
  visibility: string;
  emailNotifications: boolean;
  autoReject: boolean;
}

// Frontend-enhanced job (for display purposes)
export interface JobDisplay extends Job {
  company?: string;
  companyLogo?: string;
  status?: string;
  applicants?: number;
  posted?: string;
  expires?: string;
}

// Query filters (matches backend query params)
export interface JobFilters {
  id?: string;
  title?: string;
  companyID?: string;
  location?: string;
  minSalary?: number;
  maxSalary?: number;
  workType?: string;
  workArrangement?: string;
  postOpenDate?: string; // "1d" | "6w"
  latest?: boolean;
  sort?: "dateAsc" | "dateDesc" | "title";
}

export interface DeleteJobRequest {
  reason: string;
}


// Form-related types
export interface JobFormData {
  // Basic Info
  jobTitle?: string;
  companyID?: string;
  location?: string;
  workType?: string;
  workArrangement?: string;
  currency?: string;
  minSalary?: number | string;
  maxSalary?: number | string;
  
  // Description
  jobDescription?: string;
  jobSummary?: string;
  
  // Requirements
  requiredSkills?: string[];
  yearsOfExperience?: string;
  educationLevel?: string;
  
  // Post Settings
  postingOpenDate?: string;
  postingCloseDate?: string;
  screeningQuestions?: string;
  emailNotifications?: boolean;
  
  // Application Requirements
  applicationRequirements?: {
    resume?: boolean;
    coverLetter?: boolean;
    portfolio?: boolean;
    linkedin?: boolean;
  };
}

export interface ValidationState {
  errors: Record<string, string>;
  showErrors: boolean;
}