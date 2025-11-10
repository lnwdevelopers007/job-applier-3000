/**
 * Note Service - Business logic layer for note operations
 */

import { noteApi } from '$lib/api';
import type { Note, CreateNoteRequest, UpdateNoteRequest, NoteFilters } from '$lib/types';

export class NoteService {
  /**
   * Query notes with filters
   */
  static async queryNotes(filters?: NoteFilters): Promise<Note[]> {
    try {
      return await noteApi.query(filters);
    } catch (error) {
      console.error('Error querying notes:', error);
      throw error;
    }
  }

  /**
   * Get all notes
   */
  static async getAllNotes(): Promise<Note[]> {
    try {
      return await noteApi.getAll();
    } catch (error) {
      console.error('Error fetching all notes:', error);
      throw error;
    }
  }

  /**
   * Get all notes for a specific job application
   */
  static async getNotesByJobApplication(jobApplicationID: string): Promise<Note[]> {
    try {
      return await noteApi.query({ jobApplicationID });
    } catch (error) {
      console.error('Error fetching job application notes:', error);
      throw error;
    }
  }

  /**
   * Get a specific note by ID
   */
  static async getNoteById(id: string): Promise<Note> {
    try {
      return await noteApi.getById(id);
    } catch (error) {
      console.error('Error fetching note:', error);
      throw error;
    }
  }

  /**
   * Create a new note for a job application
   */
  static async createNote(noteData: CreateNoteRequest): Promise<Note> {
    try {
      // Validate content before creating
      const validation = this.validateNoteContent(noteData.content);
      if (!validation.isValid) {
        throw new Error(validation.error);
      }

      // Add timestamp if not provided
      const enrichedData: CreateNoteRequest = {
        ...noteData,
        timestamp: noteData.timestamp || new Date().toISOString()
      };

      const newNote = await noteApi.create(enrichedData);
      console.log('Note created successfully:', newNote.id);
      return newNote;
    } catch (error) {
      console.error('Error creating note:', error);
      throw error;
    }
  }

  /**
   * Update an existing note
   */
  static async updateNote(noteId: string, noteData: UpdateNoteRequest): Promise<Note> {
    try {
      // Validate content if it's being updated
      if (noteData.content) {
        const validation = this.validateNoteContent(noteData.content);
        if (!validation.isValid) {
          throw new Error(validation.error);
        }
      }

      // Update timestamp to current time
      const updatedData: UpdateNoteRequest = {
        ...noteData,
        timestamp: new Date().toISOString()
      };

      const updatedNote = await noteApi.update(noteId, updatedData);
      console.log('Note updated successfully:', noteId);
      return updatedNote;
    } catch (error) {
      console.error('Error updating note:', error);
      throw error;
    }
  }

  /**
   * Delete a note
   */
  static async deleteNote(noteId: string): Promise<void> {
    try {
      await noteApi.delete(noteId);
      console.log('Note deleted successfully:', noteId);
    } catch (error) {
      console.error('Error deleting note:', error);
      throw error;
    }
  }

  /**
   * Format note timestamp for display
   */
  static formatNoteTimestamp(timestamp: string): string {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
    
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric', 
      year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined 
    });
  }

  /**
   * Validate note content
   */
  static validateNoteContent(content: string): { isValid: boolean; error?: string } {
    if (!content || content.trim().length === 0) {
      return { isValid: false, error: 'Note content cannot be empty' };
    }
    if (content.length > 5000) {
      return { isValid: false, error: 'Note content cannot exceed 5000 characters' };
    }
    return { isValid: true };
  }

  /**
   * Search notes by keyword
   */
  static async searchNotes(keyword: string, filters?: NoteFilters): Promise<Note[]> {
    try {
      const notes = await noteApi.query(filters);
      const lowercaseKeyword = keyword.toLowerCase();
      return notes.filter(note => 
        note.content.toLowerCase().includes(lowercaseKeyword)
      );
    } catch (error) {
      console.error('Error searching notes:', error);
      throw error;
    }
  }

  /**
   * Get notes grouped by date
   */
  static groupNotesByDate(notes: Note[]): Map<string, Note[]> {
    const grouped = new Map<string, Note[]>();
    
    notes.forEach(note => {
      const date = new Date(note.timestamp);
      const today = new Date();
      const yesterday = new Date(today);
      yesterday.setDate(yesterday.getDate() - 1);
      
      let key: string;
      if (date.toDateString() === today.toDateString()) {
        key = 'Today';
      } else if (date.toDateString() === yesterday.toDateString()) {
        key = 'Yesterday';
      } else {
        key = date.toLocaleDateString('en-US', { 
          weekday: 'long',
          month: 'long', 
          day: 'numeric',
          year: date.getFullYear() !== today.getFullYear() ? 'numeric' : undefined
        });
      }
      
      if (!grouped.has(key)) {
        grouped.set(key, []);
      }
      grouped.get(key)!.push(note);
    });
    
    return grouped;
  }

  /**
   * Sort notes by timestamp
   */
  static sortNotesByTimestamp(notes: Note[], order: 'asc' | 'desc' = 'desc'): Note[] {
    return [...notes].sort((a, b) => {
      const timeA = new Date(a.timestamp).getTime();
      const timeB = new Date(b.timestamp).getTime();
      return order === 'desc' ? timeB - timeA : timeA - timeB;
    });
  }

  /**
   * Get the most recent note for a job application
   */
  static async getMostRecentNote(jobApplicationID: string): Promise<Note | null> {
    try {
      const notes = await this.getNotesByJobApplication(jobApplicationID);
      if (notes.length === 0) return null;
      
      const sorted = this.sortNotesByTimestamp(notes, 'desc');
      return sorted[0];
    } catch (error) {
      console.error('Error fetching most recent note:', error);
      return null;
    }
  }

  /**
   * Count notes for a job application
   */
  static async getNotesCount(jobApplicationID: string): Promise<number> {
    try {
      const notes = await this.getNotesByJobApplication(jobApplicationID);
      return notes.length;
    } catch (error) {
      console.error('Error counting notes:', error);
      return 0;
    }
  }

  /**
   * Export notes as formatted text
   */
  static exportNotesAsText(notes: Note[]): string {
    const sorted = this.sortNotesByTimestamp(notes, 'asc');
    let output = 'Job Application Notes Export\n';
    output += '=' .repeat(50) + '\n\n';

    sorted.forEach((note, index) => {
      const date = new Date(note.timestamp);
      output += `Note ${index + 1}\n`;
      output += `-`.repeat(30) + '\n';
      output += `Date: ${date.toLocaleString()}\n`;
      output += `Content:\n${note.content}\n\n`;
    });

    return output;
  }

  /**
   * Create a summary of all notes for a job application
   */
  static async getApplicationNotesSummary(jobApplicationID: string): Promise<{
    totalNotes: number;
    mostRecent: Note | null;
    oldestNote: Note | null;
  }> {
    try {
      const notes = await this.getNotesByJobApplication(jobApplicationID);
      const sorted = this.sortNotesByTimestamp(notes, 'desc');

      return {
        totalNotes: notes.length,
        mostRecent: sorted[0] || null,
        oldestNote: sorted[sorted.length - 1] || null
      };
    } catch (error) {
      console.error('Error getting notes summary:', error);
      return {
        totalNotes: 0,
        mostRecent: null,
        oldestNote: null
      };
    }
  }
}