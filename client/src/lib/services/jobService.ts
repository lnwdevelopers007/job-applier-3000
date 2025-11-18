import { goto } from '$app/navigation';
import { withDeadlineTime, formatDateCompact } from '$lib/utils/datetime';
import { fetchCompanyNameLogo, DEFAULT_COMPANY_NAME, DEFAULT_COMPANY_LOGO } from '$lib/utils/fetcher';
import { JobApplicationService } from './jobApplicationService';
import { jobApi } from '$lib/api';
import toast from 'svelte-french-toast';
import type { 
  Job, 
  JobDisplay,
  JobFilters, 
  JobFormData, 
  ValidationState 
} from '$lib/types';


export class JobService {
  /**
   * Query jobs with filters - GET /jobs/query
   */
  static async queryJobs(filters?: JobFilters): Promise<Job[]> {
    try {
      return await jobApi.query(filters);
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
      return await jobApi.getAll();
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
      return await jobApi.getById(id);
    } catch (error) {
      console.error('Error fetching job:', error);
      throw error;
    }
  }

  /**
   * Create a new job - POST /jobs/
   */
  static async createJobAPI(jobData: JobFormData): Promise<Job> {
    try {
      return await jobApi.create(jobData);
    } catch (error) {
      console.error('Error creating job:', error);
      throw error;
    }
  }

  /**
   * Update a job - PUT /jobs/:id
   */
  static async updateJobAPI(id: string, jobData: Partial<JobFormData>): Promise<Job> {
    try {
      return await jobApi.update(id, jobData);
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
      await jobApi.deleteWithReason(id, reason);
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
      return await jobApi.query({ companyID: companyId });
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

  static buildPayload(formData: JobFormData, isUpdate: boolean) {

    const payload: any = {
      // Basic Info
      title: formData.jobTitle,
      location: formData.location,
      workType: formData.workType,
      workArrangement: formData.workArrangement,
      currency: formData.currency,
      minSalary: Number(formData.minSalary),
      maxSalary: Number(formData.maxSalary),

      // Description
      jobDescription: formData.jobDescription,
      jobSummary: formData.jobSummary,

      // Requirements
      requiredSkills: Array.isArray(formData.requiredSkills) && formData.requiredSkills.length
          ? formData.requiredSkills.join(", ")
          : "",
      experienceLevel: formData.yearsOfExperience,
      education: formData.educationLevel,
      niceToHave: "",
      questions: formData.screeningQuestions,

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

    // Only include companyID when creating, not when updating
    if (!isUpdate) {
      payload.companyID = formData.companyID;
    }

    return payload;
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
      
      // Job Summary validation
      if (!formData.jobSummary?.trim()) {
        errors.jobSummary = 'Job summary is required';
      }
    }
    
    if (step === 3) {
      // Requirements validation
      if (!formData.requiredSkills || formData.requiredSkills.length === 0) {
        errors.requiredSkills = 'At least one required skill is required';
      }
      
      if (!formData.yearsOfExperience) {
        errors.yearsOfExperience = 'Years of experience is required';
      }
      
      if (!formData.educationLevel) {
        errors.educationLevel = 'Education level is required';
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
      const payload = this.buildPayload(formData, false);
      console.log("Creating job with payload:", payload);

      await jobApi.create(payload as unknown as JobFormData);
      console.log('Job created successfully');
      return { success: true };
    } catch (err) {
      console.error('Error creating job:', err);
      return { success: false, error: err instanceof Error ? err.message : 'Unknown error' };
    }
  }

  static async updateJob(jobId: string, formData: JobFormData): Promise<{ success: boolean; error?: string }> {
    try {
      const payload = this.buildPayload(formData, true);
      console.log("Updating job with payload:", payload);

      await jobApi.update(jobId, payload as unknown as Partial<JobFormData>);
      console.log('Job updated successfully');
      return { success: true };
    } catch (err) {
      console.error('Error updating job:', err);
      return { success: false, error: err instanceof Error ? err.message : 'Unknown error' };
    }
  }

  static async loadJob(jobId: string): Promise<{ success: boolean; data?: JobFormData; error?: string }> {
    try {
      const jobs = await jobApi.query({ id: jobId });
      const job = Array.isArray(jobs) ? jobs[0] : jobs;

      const formattedDeadline = job.applicationDeadline
        ? new Date(job.applicationDeadline).toISOString().slice(0, 10)
        : '';

      const formattedOpenDate = job.postOpenDate
        ? new Date(job.postOpenDate).toISOString().slice(0, 10)
        : '';

      // Destructure to exclude companyID from form data for editing
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { companyID, ...jobWithoutCompanyID } = job;
      
      const formattedData: JobFormData = {
        ...jobWithoutCompanyID,
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