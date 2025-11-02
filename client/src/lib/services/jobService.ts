import { goto } from '$app/navigation';
import { withDeadlineTime, formatDateCompact } from '$lib/utils/datetime';
import { fetchCompanyNameLogo, DEFAULT_COMPANY_NAME, DEFAULT_COMPANY_LOGO } from '$lib/utils/fetcher';
import { JobApplicationService } from './jobApplicationService';
import { apiClient } from '$lib/api/client';
import toast from 'svelte-french-toast';
import type { 
  Job, 
  JobDisplay,
  JobFilters, 
  DeleteJobRequest, 
  JobFormData, 
  ValidationState 
} from '$lib/types';


export class JobService {
  // === API METHODS (matching backend endpoints) ===

  /**
   * Query jobs with filters - GET /jobs/query
   */
  static async queryJobs(filters?: JobFilters): Promise<Job[]> {
    try {
      const response = await apiClient.get<Job[]>('/jobs/query', filters);
      return response.data;
    } catch (error) {
      console.error('Error querying jobs:', error);
      throw error;
    }
  }

  /**
   * Get all jobs - GET /jobs/
   */
  static async getAllJobs(): Promise<Job[]> {
    try {
      const response = await apiClient.get<Job[]>('/jobs/');
      return response.data;
    } catch (error) {
      console.error('Error fetching all jobs:', error);
      throw error;
    }
  }

  /**
   * Get a specific job by ID - GET /jobs/:id
   */
  static async getJobById(id: string): Promise<Job> {
    try {
      const response = await apiClient.get<Job>(`/jobs/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching job:', error);
      throw error;
    }
  }

  /**
   * Create a new job - POST /jobs/
   */
  static async createJobAPI(jobData: Partial<Job>): Promise<Job> {
    try {
      const response = await apiClient.post<Job>('/jobs/', jobData);
      return response.data;
    } catch (error) {
      console.error('Error creating job:', error);
      throw error;
    }
  }

  /**
   * Update a job - PUT /jobs/:id
   */
  static async updateJobAPI(id: string, jobData: Partial<Job>): Promise<Job> {
    try {
      const response = await apiClient.put<Job>(`/jobs/${id}`, jobData);
      return response.data;
    } catch (error) {
      console.error('Error updating job:', error);
      throw error;
    }
  }

  /**
   * Delete a job with reason - DELETE /jobs/:id
   */
  static async deleteJobWithReason(id: string, reason: string): Promise<void> {
    try {
      const deleteData: DeleteJobRequest = { reason };
      await apiClient.delete(`/jobs/${id}`, deleteData);
    } catch (error) {
      console.error('Error deleting job:', error);
      throw error;
    }
  }

  /**
   * Get jobs by company ID using query filters
   */
  static async getJobsByCompany(companyId: string): Promise<Job[]> {
    try {
      return this.queryJobs({ companyID: companyId });
    } catch (error) {
      console.error('Error fetching company jobs:', error);
      throw error;
    }
  }

  /**
   * Get latest jobs (last 3)
   */
  static async getLatestJobs(): Promise<Job[]> {
    try {
      return this.queryJobs({ latest: true });
    } catch (error) {
      console.error('Error fetching latest jobs:', error);
      throw error;
    }
  }
  static createEmptyFormData(userID?: string): JobFormData {
    return {
      // Basic Info
      jobTitle: '',
      companyID: userID || '',
      location: '',
      workType: 'full-time',
      workArrangement: 'on-site',
      currency: 'THB',
      minSalary: 1,
      maxSalary: 1,
      
      // Description
      jobDescription: '',
      jobSummary: '',
      
      // Requirements
      requiredSkills: [],
      yearsOfExperience: '',
      educationLevel: '',
      
      // Post Settings
      postingOpenDate: '',
      postingCloseDate: '',
      screeningQuestions: '',
      emailNotifications: false,
      
      // Application Requirements
      applicationRequirements: {
        resume: false,
        coverLetter: false,
        portfolio: false,
        linkedin: false
      }
    };
  }

  static buildPayload(formData: JobFormData) {

    return {
      // Basic Info
      title: formData.jobTitle || "Test Job Title",
      companyID: String(formData.companyID || "64f0c44a27b1c27f4d92e9a2"),
      location: formData.location || "Bangkok, Thailand",
      workType: formData.workType,
      workArrangement: formData.workArrangement,
      currency: formData.currency,
      minSalary: Number(formData.minSalary || 0),
      maxSalary: Number(formData.maxSalary || 0),

      // Description
      jobDescription: formData.jobDescription || "Test description",
      jobSummary: formData.jobSummary || "Test summary",

      // Requirements
      requiredSkills: Array.isArray(formData.requiredSkills) && formData.requiredSkills.length
          ? formData.requiredSkills.join(", ")
          : "JS, Node",
      experienceLevel: formData.yearsOfExperience || "Mid-Level",
      education: formData.educationLevel || "Bachelor",
      niceToHave: "",
      questions: formData.screeningQuestions || "What is your expected salary?",

      // Post Settings
      postOpenDate: formData.postingOpenDate
        ? new Date(formData.postingOpenDate).toISOString()
        : new Date().toISOString(),
      applicationDeadline: formData.postingCloseDate
        ? withDeadlineTime(formData.postingCloseDate)
        : withDeadlineTime(new Date().toISOString()),
      numberOfPositions: 1,
      visibility: "public",
      emailNotifications: Boolean(formData.emailNotifications),
      autoReject: false
    };
  }

  static validateStep(step: number, formData: JobFormData): Record<string, string> {
    const errors: Record<string, string> = {};
    
    if (step === 1) {
      // Basic Info validation
      if (!formData.jobTitle?.trim()) {
        errors.jobTitle = 'Job title is required';
      }
      
      if (!formData.location?.trim()) {
        errors.location = 'Location is required';
      }
      
      if (!formData.workType) {
        errors.workType = 'Work type is required';
      }
      
      if (!formData.workArrangement) {
        errors.workArrangement = 'Work arrangement is required';
      }
      
      // Pay range validation
      if (!formData.minSalary || Number(formData.minSalary) <= 0) {
        errors.minSalary = 'Minimum salary must be greater than 0';
      }
      
      if (!formData.maxSalary || Number(formData.maxSalary) <= 0) {
        errors.maxSalary = 'Maximum salary must be greater than 0';
      }
      
      if (formData.minSalary && formData.maxSalary && 
          Number(formData.minSalary) >= Number(formData.maxSalary)) {
        errors.maxSalary = 'Maximum salary must be greater than minimum salary';
      }
    }
    
    if (step === 2) {
      // Description validation
      if (!formData.jobDescription?.trim()) {
        errors.jobDescription = 'Job description is required';
      }
    }
    
    if (step === 4) {
      // Post Settings validation
      if (!formData.postingOpenDate) {
        errors.postingOpenDate = 'Posting open date is required';
      }
      
      if (formData.postingOpenDate && formData.postingCloseDate) {
        const openDate = new Date(formData.postingOpenDate);
        const closeDate = new Date(formData.postingCloseDate);
        
        if (closeDate <= openDate) {
          errors.postingCloseDate = 'Close date must be after open date';
        }
      }
    }
    
    return errors;
  }

  static validateAllSteps(formData: JobFormData): Record<string, string> {
    let allErrors: Record<string, string> = {};
    for (let step = 1; step <= 4; step++) {
      const stepErrors = this.validateStep(step, formData);
      allErrors = { ...allErrors, ...stepErrors };
    }
    return allErrors;
  }

  static validateStepData(step: number, formData: JobFormData): ValidationState {
    const errors = this.validateStep(step, formData);
    return {
      errors,
      showErrors: true
    };
  }

  static validateAllData(formData: JobFormData): ValidationState {
    const errors = this.validateAllSteps(formData);
    return {
      errors,
      showErrors: true
    };
  }

  static async createJob(formData: JobFormData): Promise<{ success: boolean; error?: string }> {
    try {
      const payload = this.buildPayload(formData);
      console.log("Creating job with payload:", payload);

      const res = await fetch('/jobs/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include'
      });

      if (!res.ok) {
        const text = await res.text();
        console.error('Server error body:', text);
        throw new Error(`Failed: ${res.status} ${text}`);
      }

      const data = await res.json();
      console.log('Job created:', data);
      return { success: true };
    } catch (err) {
      console.error('Error creating job:', err);
      return { success: false, error: err instanceof Error ? err.message : 'Unknown error' };
    }
  }

  static async updateJob(jobId: string, formData: JobFormData): Promise<{ success: boolean; error?: string }> {
    try {
      const payload = this.buildPayload(formData);
      console.log("Updating job with payload:", payload);

      const res = await fetch(`/jobs/${jobId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'include'
      });

      if (!res.ok) {
        throw new Error(`Failed to update job: ${res.status}`);
      }

      const updatedJob = await res.json();
      console.log('Job updated successfully:', updatedJob);
      return { success: true };
    } catch (err) {
      console.error('Error updating job:', err);
      return { success: false, error: err instanceof Error ? err.message : 'Unknown error' };
    }
  }

  static async loadJob(jobId: string): Promise<{ success: boolean; data?: JobFormData; error?: string }> {
    try {
      const res = await fetch(`/jobs/query?id=${jobId}`, {
        credentials: 'include'
      });
      if (!res.ok) throw new Error(`Failed to load job: ${res.status}`);
      
      const data = await res.json();
      const job = Array.isArray(data) ? data[0] : data;

      const formattedDeadline = job.applicationDeadline
        ? new Date(job.applicationDeadline).toISOString().slice(0, 10)
        : '';

      const formattedOpenDate = job.postOpenDate
        ? new Date(job.postOpenDate).toISOString().slice(0, 10)
        : '';

      const formattedData: JobFormData = {
        ...job,
        requiredSkills: typeof job.requiredSkills === 'string'
          ? job.requiredSkills.split(',').map((s: string) => s.trim()).filter(Boolean)
          : job.requiredSkills || [],
        postingCloseDate: formattedDeadline,
        postingOpenDate: formattedOpenDate,
        jobTitle: job.title || '',
        screeningQuestions: job.questions || '',
        educationLevel: job.education || '',
        yearsOfExperience: job.experienceLevel || ''
      };

      return { success: true, data: formattedData };
    } catch (err) {
      console.error('Error loading job:', err);
      return { success: false, error: err instanceof Error ? err.message : 'Unknown error' };
    }
  }

  static async handleStepValidation(
    step: number,
    formData: JobFormData,
    onValidationChange: (validation: ValidationState) => void
  ): Promise<boolean> {
    const validation = this.validateStepData(step, formData);
    onValidationChange(validation);
    return Object.keys(validation.errors).length === 0;
  }

  static async handleFormSubmit(
    formData: JobFormData,
    onValidationChange: (validation: ValidationState) => void,
    isEdit = false,
    jobId?: string
  ): Promise<boolean> {
    const validation = this.validateAllData(formData);
    onValidationChange(validation);
    
    if (Object.keys(validation.errors).length > 0) {
      console.error('Form has validation errors:', validation.errors);
      return false;
    }

    const result = isEdit && jobId 
      ? await this.updateJob(jobId, formData)
      : await this.createJob(formData);

    if (result.success) {
      // Show success toast with appropriate message
      if (isEdit) {
        toast.success('Job updated successfully!');
      } else {
        toast.success('Job created successfully!');
      }
      goto('/company/dashboard');
      return true;
    } else {
      console.error('Operation failed:', result.error);
      // Show error toast
      toast.error(result.error || 'Failed to save job. Please try again.');
      return false;
    }
  }

  // === ENHANCED DISPLAY METHODS ===

  /**
   * Get applicant count for a specific job
   */
  static async getApplicantCount(jobId: string): Promise<number> {
    try {
      const applications = await JobApplicationService.queryApplications({ jobID: jobId });
      return applications.length;
    } catch (err) {
      console.warn(`Failed to fetch applicant count for job ${jobId}:`, err);
      return 0;
    }
  }

  /**
   * Transform a job from API to display format
   */
  static async transformJobForDisplay(job: Job): Promise<JobDisplay> {
    // Fetch company details
    let companyName = DEFAULT_COMPANY_NAME;
    let companyLogo = DEFAULT_COMPANY_LOGO;
    try {
      const [name, logo] = await fetchCompanyNameLogo(job.companyID || '');
      companyName = name;
      companyLogo = logo || DEFAULT_COMPANY_LOGO;
    } catch (err) {
      console.warn(`Failed to fetch company for job ${job.id}:`, err);
    }

    // Get real applicant count
    const applicantCount = await this.getApplicantCount(job.id || '');

    // Determine job status
    const now = new Date();
    const closeDate = job.applicationDeadline ? new Date(job.applicationDeadline) : null;
    const isExpired = closeDate && closeDate < now;
    
    let status: JobDisplay['status'] = 'Active';
    if (isExpired) status = 'Closed';
    if (job.visibility === 'draft') status = 'Draft';

    return {
      ...job,
      company: companyName,
      companyLogo,
      status,
      location: job.location || 'Remote',
      posted: job.postOpenDate ? formatDateCompact(job.postOpenDate) : 'Unknown',
      expires: job.applicationDeadline ? formatDateCompact(job.applicationDeadline) : 'No deadline',
      applicants: applicantCount
    };
  }

  /**
   * Load jobs with enhanced display data
   */
  static async loadJobsForDisplay(filters?: JobFilters): Promise<JobDisplay[]> {
    try {
      const jobs = await this.queryJobs(filters);
      return Promise.all(jobs.map(job => this.transformJobForDisplay(job)));
    } catch (err) {
      console.error('Error loading jobs for display:', err);
      throw err;
    }
  }
}