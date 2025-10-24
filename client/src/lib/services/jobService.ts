import { goto } from '$app/navigation';
import { withDeadlineTime } from '$lib/utils/datetime';

interface JobFormData {
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
  requiredSkills?: string[] | string;
  experienceLevel?: string;
  education?: string;
  niceToHave?: string;
  screeningQuestions?: string;
  yearsOfExperience?: string;
  educationLevel?: string;
  
  // Post Settings
  postingOpenDate?: string;
  postingCloseDate?: string;
  applicationDeadline?: string;
  numberOfPositions?: number;
  visibility?: string;
  emailNotifications?: boolean;
  autoReject?: boolean;
}

interface ValidationState {
  errors: Record<string, string>;
  showErrors: boolean;
}

export class JobService {
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
      requiredSkills: '',
      experienceLevel: '',
      education: '',
      niceToHave: '',
      
      // Post Settings
      postingOpenDate: '',
      postingCloseDate: '',
      applicationDeadline: '',
      numberOfPositions: 1,
      visibility: 'public',
      emailNotifications: true,
      autoReject: false
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
      niceToHave: formData.niceToHave || "",
      questions: formData.screeningQuestions || "What is your expected salary?",

      // Post Settings
      postOpenDate: formData.postingOpenDate
        ? new Date(formData.postingOpenDate).toISOString()
        : new Date().toISOString(),
      applicationDeadline: formData.postingCloseDate
        ? withDeadlineTime(formData.postingCloseDate)
        : withDeadlineTime(new Date().toISOString()),
      numberOfPositions: Number(formData.numberOfPositions || 1),
      visibility: formData.visibility || "public",
      emailNotifications: Boolean(formData.emailNotifications),
      autoReject: Boolean(formData.autoReject)
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
      goto('/company/dashboard');
      return true;
    } else {
      console.error('Operation failed:', result.error);
      return false;
    }
  }
}