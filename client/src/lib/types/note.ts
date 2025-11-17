/**
 * Note-related types matching backend schema
 */

export interface Note {
  id?: string;
  jobApplicationID: string;
  content: string;
  timestamp: string;
}

export interface CreateNoteRequest {
  jobApplicationID: string;
  content: string;
  timestamp?: string;
}

export interface UpdateNoteRequest {
  id?: string;
  jobApplicationID?: string;
  content?: string;
  timestamp?: string;
}

export interface NoteFilters {
  jobApplicationID?: string;
}