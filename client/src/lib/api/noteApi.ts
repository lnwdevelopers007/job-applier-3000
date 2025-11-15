/**
 * Note API layer - handles all note-related HTTP requests
 */

import { ApiClient } from './client';
import type { Note, CreateNoteRequest, UpdateNoteRequest, NoteFilters } from '$lib/types';

export class NoteApi {
  constructor(private client: ApiClient) {}

  /**
   * Query notes with filters - GET /notes/ with query params
   * Backend endpoint returns all notes for jobs owned by the authenticated company user
   */
  async query(filters?: NoteFilters): Promise<Note[]> {
    const params = filters?.jobApplicationID ? { jobApplicationID: filters.jobApplicationID } : undefined;
    const response = await this.client.get<Note[]>('/notes/', params);
    
    // Sort notes by timestamp (newest first)
    const notes = response.data || [];
    return notes.sort((a, b) => {
      const timeA = new Date(a.timestamp).getTime();
      const timeB = new Date(b.timestamp).getTime();
      return timeB - timeA; // Descending order (newest first)
    });
  }

  /**
   * Get all notes - GET /notes/
   * Returns all notes for jobs owned by the authenticated company user
   */
  async getAll(): Promise<Note[]> {
    const response = await this.client.get<Note[]>('/notes/');
    
    // Sort notes by timestamp (newest first)
    const notes = response.data || [];
    return notes.sort((a, b) => {
      const timeA = new Date(a.timestamp).getTime();
      const timeB = new Date(b.timestamp).getTime();
      return timeB - timeA; // Descending order (newest first)
    });
  }

  /**
   * Get a specific note by ID - GET /notes/:id
   */
  async getById(id: string): Promise<Note> {
    const response = await this.client.get<Note>(`/notes/${id}`);
    return response.data;
  }

  /**
   * Create a new note - POST /notes/
   */
  async create(noteData: CreateNoteRequest): Promise<Note> {
    // Ensure timestamp is provided (backend requires it)
    const dataWithTimestamp = {
      ...noteData,
      timestamp: noteData.timestamp || new Date().toISOString()
    };
    const response = await this.client.post<any>('/notes/', dataWithTimestamp);
    
    // Construct the note object for optimistic update
    const createdNote: Note = {
      id: response.data.InsertedID || response.data.insertedID || response.data.id,
      jobApplicationID: noteData.jobApplicationID,
      content: noteData.content,
      timestamp: dataWithTimestamp.timestamp
    };
    
    return createdNote;
  }

  /**
   * Update a note - PUT /notes/:id
   * Backend expects PUT /notes/:id with the updated note data
   */
  async update(id: string, noteData: UpdateNoteRequest): Promise<Note> {
    // Ensure all required fields are present for update
    const updateData = {
      jobApplicationID: noteData.jobApplicationID!,
      content: noteData.content!,
      timestamp: noteData.timestamp || new Date().toISOString()
    };
    
    await this.client.put<any>(`/notes/${id}`, updateData);
    
    const updatedNote: Note = {
      id: id, // Keep the same ID
      jobApplicationID: updateData.jobApplicationID,
      content: updateData.content,
      timestamp: updateData.timestamp
    };
    
    return updatedNote;
  }

  /**
   * Delete a note - DELETE /notes/:id
   * Backend requires full note object in body for validation
   */
  async delete(id: string, noteData?: Note): Promise<void> {
    if (noteData) {
      // Send the complete note object for validation
      const body = {
        jobApplicationID: noteData.jobApplicationID,
        content: noteData.content,
        timestamp: noteData.timestamp
      };
      await this.client.delete(`/notes/${id}`, body);
    } else {
      await this.client.delete(`/notes/${id}`);
    }
  }
}