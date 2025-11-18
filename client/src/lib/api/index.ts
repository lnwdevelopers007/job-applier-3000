/**
 * API layer exports and instances
 */

import { ApiClient } from './client';
import { JobApi } from './jobApi';
import { UserApi } from './userApi';
import { JobApplicationApi } from './jobApplicationApi';
import { FileApi } from './fileApi';
import { NoteApi } from './noteApi';

// Create shared API client instance
export const apiClient = new ApiClient();

// Create domain-specific API instances
export const jobApi = new JobApi(apiClient);
export const userApi = new UserApi(apiClient);
export const jobApplicationApi = new JobApplicationApi(apiClient);
export const fileApi = new FileApi(apiClient);
export const noteApi = new NoteApi(apiClient);

// Export classes for testing or custom instances
export { ApiClient, JobApi, UserApi, JobApplicationApi, FileApi, NoteApi };

// Re-export types
export type { ApiResponse, ApiError } from '$lib/types';
export type { FileMetadata, ApplicantFilesResponse } from './fileApi';