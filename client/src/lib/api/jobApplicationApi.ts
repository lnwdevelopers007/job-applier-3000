/**
 * Job Application API layer - handles all job application-related HTTP requests
 */

import { ApiClient } from './client';
import type { JobApplication, JobApplicationFilters, JobApplicationWithApplicant } from '$lib/types';

export class JobApplicationApi {
  constructor(private client: ApiClient) {}

  /**
   * Query job applications with filters - GET /apply/query
   */
  async query(filters?: JobApplicationFilters): Promise<JobApplication[]> {
    const response = await this.client.get<JobApplication[]>('/apply/query', filters);
    return response.data;
  }

  /**
   * Get all job applications - GET /apply/
   */
  async getAll(): Promise<JobApplication[]> {
    const response = await this.client.get<JobApplication[]>('/apply/');
    return response.data;
  }

  /**
   * Get a specific application by ID - GET /apply/:id
   */
  async getById(id: string): Promise<JobApplication> {
    const response = await this.client.get<JobApplication>(`/apply/${id}`);
    return response.data;
  }

  /**
   * Create a new application - POST /apply/
   */
  async create(applicationData: Partial<JobApplication>): Promise<JobApplication> {
    const response = await this.client.post<JobApplication>('/apply/', applicationData);
    return response.data;
  }

  /**
   * Update an application - PUT /apply/:id
   */
  async update(id: string, applicationData: Partial<JobApplication>): Promise<JobApplication> {
    const response = await this.client.put<JobApplication>(`/apply/${id}`, applicationData);
    return response.data;
  }

  /**
   * Delete an application - DELETE /apply/:id
   */
  async delete(id: string): Promise<void> {
    await this.client.delete(`/apply/${id}`);
  }

  /**
   * Get applications by applicant ID
   */
  async getByApplicantId(applicantId: string): Promise<JobApplication[]> {
    return this.query({ applicantID: applicantId });
  }

  /**
   * Get applications by job ID
   */
  async getByJobId(jobId: string): Promise<JobApplicationWithApplicant[]> {
    return this.query({ jobID: jobId }) as Promise<JobApplicationWithApplicant[]>;
  }

  /**
   * Get applications by company ID
   */
  async getByCompanyId(companyId: string): Promise<JobApplication[]> {
    return this.query({ companyID: companyId });
  }

  /**
   * Get applications by status
   */
  async getByStatus(status: string): Promise<JobApplication[]> {
    return this.query({ status });
  }
}