/**
 * Note API layer - handles all note-related HTTP requests
 */

import { ApiClient } from './client';
import type { Note, CreateNoteRequest, UpdateNoteRequest, NoteFilters } from '$lib/types';

export class NoteApi {
  private mockNotes: Map<string, Note> = new Map();
  private mockMode = true; // Toggle for mock mode while backend is in PR
  
  constructor(private client: ApiClient) {
    // Initialize with some mock data
    this.initMockData();
  }

  /**
   * Query notes with filters - GET /notes/query
   */
  async query(filters?: NoteFilters): Promise<Note[]> {
    if (this.mockMode) {
      return this.mockQuery(filters);
    }
    const response = await this.client.get<Note[]>('/notes/query', filters);
    return response.data;
  }

  /**
   * Get all notes - GET /notes/
   */
  async getAll(): Promise<Note[]> {
    if (this.mockMode) {
      return Array.from(this.mockNotes.values());
    }
    const response = await this.client.get<Note[]>('/notes/');
    return response.data;
  }

  /**
   * Get a specific note by ID - GET /notes/:id
   */
  async getById(id: string): Promise<Note> {
    if (this.mockMode) {
      const note = this.mockNotes.get(id);
      if (!note) {
        throw new Error(`Note with id ${id} not found`);
      }
      return note;
    }
    const response = await this.client.get<Note>(`/notes/${id}`);
    return response.data;
  }

  /**
   * Create a new note - POST /notes/
   */
  async create(noteData: CreateNoteRequest): Promise<Note> {
    if (this.mockMode) {
      return this.mockCreate(noteData);
    }
    const response = await this.client.post<Note>('/notes/', noteData);
    return response.data;
  }

  /**
   * Update a note - PUT /notes/
   */
  async update(id: string, noteData: UpdateNoteRequest): Promise<Note> {
    if (this.mockMode) {
      return this.mockUpdate(id, noteData);
    }
    // Backend uses PUT /notes/ with full note data
    const response = await this.client.put<Note>('/notes/', { ...noteData, id });
    return response.data;
  }

  /**
   * Delete a note - DELETE /notes/:id
   */
  async delete(id: string): Promise<void> {
    if (this.mockMode) {
      this.mockDelete(id);
      return;
    }
    await this.client.delete(`/notes/${id}`);
  }

  // ============ Mock Methods ============

  private initMockData(): void {
    const now = new Date();
    const mockData: Note[] = [
      {
        id: '507f1f77bcf86cd799439011',
        jobApplicationID: '507f1f77bcf86cd799439001',
        content: 'Strong technical background with React and Node.js experience. Good communication skills during the initial screening. Candidate showed enthusiasm for the role and company culture.',
        timestamp: new Date(now.getTime() - 86400000 * 3).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439012',
        jobApplicationID: '507f1f77bcf86cd799439001',
        content: 'Follow up needed - candidate mentioned availability for next week. Prefers Tuesday or Thursday afternoons for the technical interview.',
        timestamp: new Date(now.getTime() - 43200000).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439013',
        jobApplicationID: '507f1f77bcf86cd799439002',
        content: 'Excellent portfolio showcasing several production applications. Has experience with our tech stack including MongoDB, Express, and React. Previously worked on similar e-commerce projects.',
        timestamp: new Date(now.getTime() - 86400000 * 2).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439014',
        jobApplicationID: '507f1f77bcf86cd799439003',
        content: 'Salary expectation is within our budget range. Candidate is looking for remote work opportunities and values work-life balance.',
        timestamp: new Date(now.getTime() - 86400000).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439015',
        jobApplicationID: '507f1f77bcf86cd799439002',
        content: 'Technical assessment completed - scored 85/100. Strong in algorithms and system design, some gaps in DevOps knowledge but eager to learn.',
        timestamp: new Date(now.getTime() - 3600000 * 6).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439016',
        jobApplicationID: '507f1f77bcf86cd799439004',
        content: 'First round interview scheduled for next Monday at 2 PM. Sent calendar invite and preparation materials.',
        timestamp: new Date(now.getTime() - 3600000 * 2).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439017',
        jobApplicationID: '507f1f77bcf86cd799439005',
        content: 'Reference check completed. Previous employer gave positive feedback about work ethic and team collaboration. No red flags identified.',
        timestamp: new Date(now.getTime() - 86400000 * 4).toISOString()
      },
      {
        id: '507f1f77bcf86cd799439018',
        jobApplicationID: '507f1f77bcf86cd799439003',
        content: 'Candidate has competing offers. Need to expedite our decision-making process if we want to proceed.',
        timestamp: new Date(now.getTime() - 3600000).toISOString()
      }
    ];

    mockData.forEach(note => this.mockNotes.set(note.id!, note));
  }

  private async mockQuery(filters?: NoteFilters): Promise<Note[]> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 300));

    let notes = Array.from(this.mockNotes.values());

    if (filters?.jobApplicationID) {
      notes = notes.filter(n => n.jobApplicationID === filters.jobApplicationID);
    }

    // Sort by timestamp (newest first)
    return notes.sort((a, b) => 
      new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );
  }

  private async mockCreate(noteData: CreateNoteRequest): Promise<Note> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 500));

    const newNote: Note = {
      id: new Date().getTime().toString(16), // Generate hex-like ID
      jobApplicationID: noteData.jobApplicationID,
      content: noteData.content,
      timestamp: noteData.timestamp || new Date().toISOString()
    };

    this.mockNotes.set(newNote.id!, newNote);
    return newNote;
  }

  private async mockUpdate(id: string, noteData: UpdateNoteRequest): Promise<Note> {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 400));

    const existingNote = this.mockNotes.get(id);
    if (!existingNote) {
      throw new Error(`Note with id ${id} not found`);
    }

    const updatedNote: Note = {
      ...existingNote,
      jobApplicationID: noteData.jobApplicationID || existingNote.jobApplicationID,
      content: noteData.content || existingNote.content,
      timestamp: noteData.timestamp || new Date().toISOString()
    };

    this.mockNotes.set(id, updatedNote);
    return updatedNote;
  }

  private mockDelete(id: string): void {
    if (!this.mockNotes.has(id)) {
      throw new Error(`Note with id ${id} not found`);
    }
    this.mockNotes.delete(id);
  }

  /**
   * Toggle mock mode (for development/testing)
   */
  setMockMode(enabled: boolean): void {
    this.mockMode = enabled;
  }
}