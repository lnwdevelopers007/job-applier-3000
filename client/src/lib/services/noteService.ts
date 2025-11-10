import { apiFetch } from '$lib/utils/api';

const API_BASE = import.meta.env.VITE_BACKEND || 'http://localhost:8080';

export interface Note {
	id: string;
	candidateId: string;
	content: string;
	createdBy: string;
	createdByName?: string;
	createdAt: string;
	updatedAt: string;
}

export interface CreateNoteRequest {
	candidateId: string;
	content: string;
}

export interface UpdateNoteRequest {
	content: string;
}

export class NoteService {
	/**
	 * Get all notes for a candidate
	 */
	static async getNotesByCandidate(candidateId: string): Promise<Note[]> {
		try {
			// Mock data for development since API endpoint not merged yet
			const mockNotes: Note[] = [
				{
					id: '1',
					candidateId,
					content: 'Strong technical background with React and Node.js experience. Good communication skills during the initial screening.',
					createdBy: 'hr1',
					createdByName: 'Sarah Johnson (HR)',
					createdAt: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
					updatedAt: new Date(Date.now() - 86400000).toISOString()
				},
				{
					id: '2',
					candidateId,
					content: 'Follow up needed - candidate mentioned availability for next week.',
					createdBy: 'manager1',
					createdByName: 'Mike Chen (Hiring Manager)',
					createdAt: new Date(Date.now() - 43200000).toISOString(), // 12 hours ago
					updatedAt: new Date(Date.now() - 43200000).toISOString()
				}
			];

			// Simulate API delay
			await new Promise(resolve => setTimeout(resolve, 300));
			return mockNotes;

			// Real implementation when API is ready:
			// const response = await apiFetch(`${API_BASE}/notes/candidate/${candidateId}`);
			// if (!response.ok) {
			// 	throw new Error(`Failed to fetch notes: ${response.status}`);
			// }
			// return await response.json();
		} catch (error) {
			console.error('Error fetching candidate notes:', error);
			throw error;
		}
	}

	/**
	 * Create a new note for a candidate
	 */
	static async createNote(noteData: CreateNoteRequest): Promise<Note> {
		try {
			// Mock implementation
			const newNote: Note = {
				id: Date.now().toString(),
				candidateId: noteData.candidateId,
				content: noteData.content,
				createdBy: 'current-user',
				createdByName: 'You',
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString()
			};

			// Simulate API delay
			await new Promise(resolve => setTimeout(resolve, 500));
			return newNote;

			// Real implementation when API is ready:
			// const response = await apiFetch(`${API_BASE}/notes`, {
			// 	method: 'POST',
			// 	headers: {
			// 		'Content-Type': 'application/json'
			// 	},
			// 	body: JSON.stringify(noteData)
			// });
			// if (!response.ok) {
			// 	throw new Error(`Failed to create note: ${response.status}`);
			// }
			// return await response.json();
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
			// Mock implementation
			const updatedNote: Note = {
				id: noteId,
				candidateId: 'mock-candidate',
				content: noteData.content,
				createdBy: 'current-user',
				createdByName: 'You',
				createdAt: new Date(Date.now() - 86400000).toISOString(),
				updatedAt: new Date().toISOString()
			};

			// Simulate API delay
			await new Promise(resolve => setTimeout(resolve, 400));
			return updatedNote;

			// Real implementation when API is ready:
			// const response = await apiFetch(`${API_BASE}/notes/${noteId}`, {
			// 	method: 'PUT',
			// 	headers: {
			// 		'Content-Type': 'application/json'
			// 	},
			// 	body: JSON.stringify(noteData)
			// });
			// if (!response.ok) {
			// 	throw new Error(`Failed to update note: ${response.status}`);
			// }
			// return await response.json();
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
			// Mock implementation - just simulate delay
			await new Promise(resolve => setTimeout(resolve, 300));

			// Real implementation when API is ready:
			// const response = await apiFetch(`${API_BASE}/notes/${noteId}`, {
			// 	method: 'DELETE'
			// });
			// if (!response.ok) {
			// 	throw new Error(`Failed to delete note: ${response.status}`);
			// }
		} catch (error) {
			console.error('Error deleting note:', error);
			throw error;
		}
	}
}