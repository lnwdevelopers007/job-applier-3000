/**
 * Job Application Service - handles all job application-related API operations
 */

import { apiClient } from '$lib/api/client';
import type { JobApplication, JobApplicationFilters } from '$lib/types/jobApplication';

export class JobApplicationService {
  /**
   * Query job applications with filters - GET /apply/query
   */
  static async queryApplications(filters?: JobApplicationFilters): Promise<JobApplication[]> {
    try {
      const response = await apiClient.get<JobApplication[]>('/apply/query', filters);
      return response.data;
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
      const response = await apiClient.get<JobApplication[]>('/apply/');
      return response.data;
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
      const response = await apiClient.get<JobApplication>(`/apply/${id}`);
      return response.data;
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
      const response = await apiClient.post<JobApplication>('/apply/', applicationData);
      return response.data;
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
      const response = await apiClient.put<JobApplication>(`/apply/${id}`, applicationData);
      return response.data;
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
      await apiClient.delete(`/apply/${id}`);
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
      return this.queryApplications({ applicantID: applicantId });
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
      return this.queryApplications({ jobID: jobId });
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
      return this.queryApplications({ companyID: companyId });
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
      return this.queryApplications({ status });
    } catch (error) {
      console.error('Error fetching applications by status:', error);
      throw error;
    }
  }
}