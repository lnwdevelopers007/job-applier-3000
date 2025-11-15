<script lang="ts">
  import Modal from './Modal.svelte';
  import DeleteModal from './DeleteModal.svelte';
  import { NoteService } from '$lib/services/noteService';
  import type { Note } from '$lib/types';
  import { Plus, LoaderCircle, Pencil, Trash2 } from 'lucide-svelte';
  import toast from 'svelte-french-toast';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    jobApplicationId: string;
    candidateName: string;
    candidateAvatar?: string;
  }

  let {
    isOpen = $bindable(),
    onClose,
    jobApplicationId,
    candidateName,
    candidateAvatar
  }: Props = $props();

  let notes = $state<Note[]>([]);
  let newNoteContent = $state('');
  let loading = $state(false);
  let submitting = $state(false);
  
  // Editing state
  let editingNoteId = $state<string | null>(null);
  let editContent = $state('');
  let updating = $state(false);
  
  // Delete confirmation state
  let deleteConfirmOpen = $state(false);
  let noteToDelete = $state<Note | null>(null);
  let isDeleting = $state(false);

  async function loadNotes() {
    if (!jobApplicationId) return;
    
    try {
      loading = true;
      const result = await NoteService.getNotesByJobApplication(jobApplicationId);
      notes = result || []; // Ensure notes is always an array
    } catch (err) {
      console.error('Failed to load notes:', err);
      toast.error('Failed to load notes. Please try again.');
      notes = []; // Set to empty array on error
    } finally {
      loading = false;
    }
  }

  async function addNote() {
    if (!newNoteContent.trim()) return;

    try {
      submitting = true;
      const newNote = await NoteService.createNote({
        jobApplicationID: jobApplicationId,
        content: newNoteContent.trim()
      });
      
      notes = [newNote, ...notes]; // Add to top
      newNoteContent = '';
      toast.success('Note added successfully');
    } catch (err) {
      console.error('Failed to add note:', err);
      toast.error('Failed to add note. Please try again.');
    } finally {
      submitting = false;
    }
  }

  function startEditing(note: Note) {
    if (!note.id) return;
    editingNoteId = note.id;
    editContent = note.content;
  }

  function cancelEditing() {
    editingNoteId = null;
    editContent = '';
  }

  async function saveEdit(noteId: string) {
    if (!editContent.trim()) return;

    // Find the original note to preserve its jobApplicationID format
    const originalNote = notes.find(note => note.id === noteId);
    if (!originalNote) {
      toast.error('Note not found');
      return;
    }

    try {
      updating = true;
      const updatedNote = await NoteService.updateNote(noteId, {
        jobApplicationID: originalNote.jobApplicationID, // Preserve original format
        content: editContent.trim()
      });
      
      // Update the note in the list
      notes = notes.map(note => 
        note.id === noteId ? updatedNote : note
      );
      
      editingNoteId = null;
      editContent = '';
      toast.success('Note updated successfully');
    } catch (err) {
      console.error('Failed to update note:', err);
      toast.error('Failed to update note. Please try again.');
    } finally {
      updating = false;
    }
  }

  function confirmDeleteNote(note: Note) {
    noteToDelete = note;
    deleteConfirmOpen = true;
  }

  async function handleDeleteConfirm(_reason: string) {
    if (!noteToDelete) return;

    try {
      isDeleting = true;
      await NoteService.deleteNote(noteToDelete.id!, noteToDelete);
      notes = notes.filter(note => note.id !== noteToDelete!.id);
      toast.success('Note deleted successfully');
      deleteConfirmOpen = false;
      noteToDelete = null;
    } catch (err) {
      console.error('Failed to delete note:', err);
      toast.error('Failed to delete note. Please try again.');
    } finally {
      isDeleting = false;
    }
  }

  function cancelDelete() {
    deleteConfirmOpen = false;
    noteToDelete = null;
  }

  function handleClose() {
    newNoteContent = '';
    onClose();
  }

  // Load notes when modal opens
  $effect(() => {
    if (isOpen && jobApplicationId) {
      loadNotes();
    }
  });

  // Handle enter key in textarea
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) {
      event.preventDefault();
      addNote();
    }
  }

  // Format timestamp for display
  function formatRelativeTime(timestamp: string): string {
    return NoteService.formatNoteTimestamp(timestamp);
  }
</script>

<Modal bind:isOpen={isOpen} size="xl" closeOnBackdrop={false} onClose={handleClose}>
  <div class="p-6">
    <!-- Header -->
    <div class="flex items-center justify-between pb-4">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full overflow-hidden">
          {#if candidateAvatar}
            <img src={candidateAvatar} alt={candidateName} class="w-8 h-8 object-cover" />
          {:else}
            <div class="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
              <span class="text-sm font-semibold text-gray-600">
                {candidateName ? candidateName.charAt(0).toUpperCase() : 'U'}
              </span>
            </div>
          {/if}
        </div>
        <h2 class="text-xl font-medium text-gray-900">{candidateName}'s Notes</h2>
      </div>
    </div>


    <!-- Notes List -->
    <div class="space-y-4 h-72 overflow-y-auto">
      {#if loading}
        <div class="flex items-center justify-center py-8">
          <LoaderCircle class="animate-spin text-gray-500 w-8 h-8" />
        </div>
      {:else if notes.length === 0}
        <div class="text-center py-8">
          <p class="text-sm text-gray-500 mb-1">No notes yet</p>
          <p class="text-xs text-gray-400">Add the first note about this candidate</p>
        </div>
      {:else}
        {#each notes as note (note.id)}
          <div class="bg-gray-100 rounded-xl p-4">
            {#if editingNoteId === note.id}
              <!-- Edit mode -->
              <div class="space-y-2">
                <textarea
                  bind:value={editContent}
                  rows="3"
                  maxlength="1000"
                  class="w-full p-2 border border-gray-300 rounded text-sm resize-none focus:outline-none focus:ring-1 focus:ring-gray-400"
                  disabled={updating}
                ></textarea>
                <div class="flex items-center justify-between">
                  <span class="text-xs text-gray-500">{editContent.length}/1000</span>
                  <div class="flex gap-2">
                    <button
                      class="flex items-center px-3 py-1.5 text-xs rounded-md text-gray-600 hover:text-gray-700"
                      onclick={cancelEditing}
                      disabled={updating}
                    >
                      Cancel
                    </button>
                    <button
                      class="flex items-center px-3 py-1.5 text-xs bg-green-600 font-medium text-white rounded-md hover:bg-green-700 disabled:opacity-50"
                      disabled={!editContent.trim() || updating}
                      onclick={() => note.id && saveEdit(note.id)}
                    >
                      Save
                    </button>
            
                  </div>
                </div>
              </div>
            {:else}
              <!-- View mode -->
              <div class="group">
                <div class="flex items-start gap-3">
                  <p class="text-sm text-gray-700 leading-relaxed flex-1 min-w-0 whitespace-pre-wrap break-words">{note.content}</p>
                  <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                    <button
                      class="p-1 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded"
                      onclick={() => startEditing(note)}
                      title="Edit note"
                    >
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button
                      class="p-1 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded"
                      onclick={() => confirmDeleteNote(note)}
                      title="Delete note"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <div class="mt-2 flex items-center justify-between">
                  <div class="text-xs text-gray-400">
                    <span>{formatRelativeTime(note.timestamp)}</span>
                  </div>
                </div>
              </div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Modal Footer - Add Note Form -->
  <div class="border-t border-gray-200 p-4">
    <div class="space-y-3">
      <textarea
        bind:value={newNoteContent}
        placeholder="Add a note about this candidate..."
        rows="3"
        maxlength="1000"
        class="w-full p-3 border border-gray-300 rounded-lg text-sm resize-none focus:outline-none focus:ring-1 focus:ring-gray-400"
        disabled={submitting}
        onkeydown={handleKeydown}
      ></textarea>
      <div class="flex items-center justify-between">
        <span class="text-xs text-gray-500">
          {newNoteContent.length}/1000
          {#if newNoteContent.length > 0}
            • Press Cmd+Enter to submit
          {/if}
        </span>
        <button
          class="flex items-center gap-1.5 px-4 py-2 text-sm bg-green-600 font-medium text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={!newNoteContent.trim() || submitting}
          onclick={addNote}
        >
          <Plus class="w-4 h-4" />
          {submitting ? 'Adding...' : 'Add Note'}
        </button>
      </div>
    </div>
  </div>
</Modal>

<!-- Delete Confirmation Modal -->
{#if deleteConfirmOpen && noteToDelete}
  <DeleteModal
    bind:isOpen={deleteConfirmOpen}
    onClose={cancelDelete}
    onConfirm={handleDeleteConfirm}
    title="Delete Note"
    itemName={noteToDelete.content.slice(0, 50) + (noteToDelete.content.length > 50 ? '...' : '')}
    description="Are you sure you want to delete this note? This action cannot be undone."
    reasonLabel=""
    confirmButtonText="Delete Note"
    {isDeleting}
  />
{/if}