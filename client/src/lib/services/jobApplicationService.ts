/**
 * Job Application Service - handles all job application-related business logic
 */

import { jobApplicationApi } from '$lib/api';
import type { JobApplication, JobApplicationFilters } from '$lib/types';

export class JobApplicationService {
  /**
   * Query job applications with filters - GET /apply/query
   */
  static async queryApplications(filters?: JobApplicationFilters): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.query(filters);
    } catch (error) {
      console.error('Error querying job applications:', error);
      throw error;
    }
  }

  /**
   * Get all job applications - GET /apply/
   */
  static async getAllApplications(): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.getAll();
    } catch (error) {
      console.error('Error fetching all applications:', error);
      throw error;
    }
  }

  /**
   * Get a specific application by ID - GET /apply/:id
   */
  static async getApplicationById(id: string): Promise<JobApplication> {
    try {
      return await jobApplicationApi.getById(id);
    } catch (error) {
      console.error('Error fetching application:', error);
      throw error;
    }
  }

  /**
   * Create a new application - POST /apply/
   */
  static async createApplication(applicationData: Partial<JobApplication>): Promise<JobApplication> {
    try {
      return await jobApplicationApi.create(applicationData);
    } catch (error) {
      console.error('Error creating application:', error);
      throw error;
    }
  }

  /**
   * Update an application - PUT /apply/:id
   */
  static async updateApplication(id: string, applicationData: Partial<JobApplication>): Promise<JobApplication> {
    try {
      return await jobApplicationApi.update(id, applicationData);
    } catch (error) {
      console.error('Error updating application:', error);
      throw error;
    }
  }

  /**
   * Delete an application - DELETE /apply/:id
   */
  static async deleteApplication(id: string): Promise<void> {
    try {
      await jobApplicationApi.delete(id);
    } catch (error) {
      console.error('Error deleting application:', error);
      throw error;
    }
  }

  /**
   * Get applications by applicant ID
   */
  static async getApplicationsByApplicant(applicantId: string): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.getByApplicantId(applicantId);
    } catch (error) {
      console.error('Error fetching applications by applicant:', error);
      throw error;
    }
  }

  /**
   * Get applications by job ID
   */
  static async getApplicationsByJob(jobId: string): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.getByJobId(jobId);
    } catch (error) {
      console.error('Error fetching applications by job:', error);
      throw error;
    }
  }

  /**
   * Get applications by company ID
   */
  static async getApplicationsByCompany(companyId: string): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.getByCompanyId(companyId);
    } catch (error) {
      console.error('Error fetching applications by company:', error);
      throw error;
    }
  }

  /**
   * Get applications by status
   */
  static async getApplicationsByStatus(status: string): Promise<JobApplication[]> {
    try {
      return await jobApplicationApi.getByStatus(status);
    } catch (error) {
      console.error('Error fetching applications by status:', error);
      throw error;
    }
  }
  
  /**
   * Get application statistics for a company
   */
  static async getApplicationStats(companyId: string): Promise<{
    total: number;
    pending: number;
    approved: number;
    rejected: number;
  }> {
    try {
      const applications = await this.getApplicationsByCompany(companyId);
      
      return {
        total: applications.length,
        pending: applications.filter(app => app.status === 'pending').length,
        approved: applications.filter(app => app.status === 'approved').length,
        rejected: applications.filter(app => app.status === 'rejected').length
      };
    } catch (error) {
      console.error('Error calculating application stats:', error);
      throw error;
    }
  }

  /**
   * Check if user has already applied to a job
   */
  static async hasUserAppliedToJob(applicantId: string, jobId: string): Promise<boolean> {
    try {
      const applications = await this.queryApplications({ 
        applicantID: applicantId, 
        jobID: jobId 
      });
      return applications.length > 0;
    } catch (error) {
      console.error('Error checking application status:', error);
      return false;
    }
  }

  /**
   * Get recent applications for dashboard
   */
  static async getRecentApplications(companyId: string, limit: number = 5): Promise<JobApplication[]> {
    try {
      const applications = await this.getApplicationsByCompany(companyId);
      // Sort by application date (most recent first) and limit
      return applications
        .sort((a, b) => new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime())
        .slice(0, limit);
    } catch (error) {
      console.error('Error fetching recent applications:', error);
      throw error;
    }
  }
}