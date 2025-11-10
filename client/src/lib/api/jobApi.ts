/**
 * Job API layer - handles all job-related HTTP requests
 */

import { ApiClient } from './client';
import type { Job, JobFilters, DeleteJobRequest, JobFormData } from '$lib/types';

export class JobApi {
  constructor(private client: ApiClient) {}

  /**
   * Query jobs with filters - GET /jobs/query
   */
  async query(filters?: JobFilters): Promise<Job[]> {
    const response = await this.client.get<Job[]>('/jobs/query', filters);
    return response.data;
  }

  /**
   * Get all jobs - GET /jobs/
   */
  async getAll(): Promise<Job[]> {
    const response = await this.client.get<Job[]>('/jobs/');
    return response.data;
  }

  /**
   * Get a specific job by ID - GET /jobs/:id
   */
  async getById(id: string): Promise<Job> {
    const response = await this.client.get<Job>(`/jobs/${id}`);
    return response.data;
  }

  /**
   * Create a new job - POST /jobs/
   */
  async create(jobData: JobFormData): Promise<Job> {
    const response = await this.client.post<Job>('/jobs/', jobData);
    return response.data;
  }

  /**
   * Update a job - PUT /jobs/:id
   */
  async update(id: string, jobData: Partial<JobFormData>): Promise<Job> {
    const response = await this.client.put<Job>(`/jobs/${id}`, jobData);
    return response.data;
  }

  /**
   * Delete a job with reason - DELETE /jobs/:id
   */
  async deleteWithReason(id: string, reason: string): Promise<void> {
    const deleteRequest: DeleteJobRequest = { reason };
    await this.client.delete(`/jobs/${id}`, deleteRequest);
  }

  /**
   * Get jobs by company ID - GET /jobs/company/:companyId
   */
  async getByCompanyId(companyId: string): Promise<Job[]> {
    const response = await this.client.get<Job[]>(`/jobs/company/${companyId}`);
    return response.data;
  }
}