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

export interface UploadResponse {
  id: string;
  userID: string;
  filename: string;
  fileExtension: string;
  contentType: string;
  size: number;
  category: string;
  uploadDate: string;
}

class FileService {
  /**
   * Upload a file with category
   */
  async uploadFile(file: File, category: string): Promise<FileMetadata> {
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
   * List all files for a user
   */
  async listUserFiles(userId: string): Promise<FileMetadata[]> {
    const response = await fetch(`/files/user/${userId}`, {
      credentials: 'include'
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Failed to fetch files' }));
      throw new Error(error.error || 'Failed to fetch files');
    }

    const data = await response.json();
    return data.files || [];
  }

  /**
   * Download a file as Blob
   */
  async downloadFile(fileId: string): Promise<Blob> {
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
   * Get file URL for preview (creates object URL from blob)
   */
  async getFilePreviewUrl(fileId: string): Promise<string> {
    const blob = await this.downloadFile(fileId);
    return URL.createObjectURL(blob);
  }

  /**
   * Delete a file
   */
  async deleteFile(fileId: string): Promise<void> {
    const response = await fetch(`/files/${fileId}`, {
      method: 'DELETE',
      credentials: 'include'
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete file');
    }
  }

  /**
   * Get applicant files for a job application (company use)
   */
  async getApplicantFiles(applicationId: string): Promise<{
    applicationID: string;
    applicantID: string;
    jobID: string;
    files: FileMetadata[];
  }> {
    const response = await fetch(`/files/application/${applicationId}`, {
      credentials: 'include'
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to fetch applicant files');
    }

    return await response.json();
  }

  /**
   * Download an applicant's file (company use)
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
   * Get applicant file preview URL (company use)
   */
  async getApplicantFilePreviewUrl(applicationId: string, fileId: string): Promise<string> {
    const blob = await this.downloadApplicantFile(applicationId, fileId);
    return URL.createObjectURL(blob);
  }

  /**
   * Format file size for display
   */
  formatFileSize(bytes: number): string {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }

  /**
   * Get category display label
   */
  getCategoryLabel(category: string): string {
    const labels: Record<string, string> = {
      resume: 'Resume',
      transcript: 'Transcript',
      certification: 'Certification',
      verification: 'Verification'
    };
    return labels[category] || category;
  }

  /**
   * Get category color for badges
   */
  getCategoryColor(category: string): string {
    const colors: Record<string, string> = {
      resume: 'blue',
      transcript: 'purple',
      certification: 'green',
      verification: 'yellow'
    };
    return colors[category] || 'gray';
  }

  /**
   * Get valid categories for user role
   */
  getValidCategories(role: string): { value: string; label: string }[] {
    if (role === 'jobSeeker') {
      return [
        { value: 'resume', label: 'Resume' },
        { value: 'transcript', label: 'Transcript' },
        { value: 'certification', label: 'Certification' }
      ];
    }
    if (role === 'company') {
      return [
        { value: 'verification', label: 'Verification' }
      ];
    }
    return [];
  }
}

export const fileService = new FileService();