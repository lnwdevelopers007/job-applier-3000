/**
 * File API layer - handles all file-related HTTP requests
 */

import { ApiClient } from './client';

export interface FileMetadata {
  id: string;
  userID: string;
  filename: string;
  fileExtension: string;
  contentType: string;
  size: number;
  category: string;
  uploadDate: string;
}

export interface ApplicantFilesResponse {
  applicationID: string;
  applicantID: string;
  jobID: string;
  files: FileMetadata[];
}

export class FileApi {
  constructor(private client: ApiClient) {}

  /**
   * Upload a file with category - POST /files/upload
   * Note: This uses FormData and doesn't go through the standard JSON client
   */
  async upload(file: File, category: string): Promise<FileMetadata> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('category', category);

    const response = await fetch('/files/upload', {
      method: 'POST',
      credentials: 'include',
      body: formData
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Upload failed' }));
      throw new Error(error.error || 'Failed to upload file');
    }

    return await response.json();
  }

  /**
   * List all files for a user - GET /files/user/:userId
   */
  async listUserFiles(userId: string): Promise<FileMetadata[]> {
    const response = await this.client.get<{ files: FileMetadata[] }>(`/files/user/${userId}`);
    return response.data.files || [];
  }

  /**
   * Download a file as Blob - GET /files/download/:id
   */
  async download(fileId: string): Promise<Blob> {
    const response = await fetch(`/files/download/${fileId}`, {
      credentials: 'include'
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Download failed' }));
      throw new Error(error.error || 'Failed to download file');
    }

    return await response.blob();
  }

  /**
   * Delete a file - DELETE /files/:id
   */
  async delete(fileId: string): Promise<void> {
    await this.client.delete(`/files/${fileId}`);
  }

  /**
   * Get applicant files for a job application - GET /files/application/:applicationId
   */
  async getApplicantFiles(applicationId: string): Promise<ApplicantFilesResponse> {
    const response = await this.client.get<ApplicantFilesResponse>(`/files/application/${applicationId}`);
    return response.data;
  }

  /**
   * Download an applicant's file - GET /files/application/:applicationId/download/:fileId
   */
  async downloadApplicantFile(applicationId: string, fileId: string): Promise<Blob> {
    const response = await fetch(`/files/application/${applicationId}/download/${fileId}`, {
      credentials: 'include'
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Download failed' }));
      throw new Error(error.error || 'Failed to download file');
    }

    return await response.blob();
  }

  /**
   * Get file URL for preview (creates object URL from blob)
   */
  async getFilePreviewUrl(fileId: string): Promise<string> {
    const blob = await this.download(fileId);
    return URL.createObjectURL(blob);
  }

  /**
   * Get applicant file preview URL
   */
  async getApplicantFilePreviewUrl(applicationId: string, fileId: string): Promise<string> {
    const blob = await this.downloadApplicantFile(applicationId, fileId);
    return URL.createObjectURL(blob);
  }
}