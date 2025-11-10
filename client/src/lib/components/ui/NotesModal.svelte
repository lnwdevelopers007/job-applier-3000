<script lang="ts">
  import Modal from './Modal.svelte';
  import { NoteService, type Note } from '$lib/services/noteService';
  import { Plus, CircleAlert, LoaderCircle, Pencil, Trash2, Save, X } from 'lucide-svelte';
  import { formatRelativeTime } from '$lib/utils/datetime';
  import toast from 'svelte-french-toast';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    candidateId: string;
    candidateName: string;
    candidateAvatar?: string;
  }

  let {
    isOpen = $bindable(),
    onClose,
    candidateId,
    candidateName,
    candidateAvatar
  }: Props = $props();

  let notes = $state<Note[]>([]);
  let newNoteContent = $state('');
  let loading = $state(false);
  let submitting = $state(false);
  let error = $state<string | null>(null);
  
  // Editing state
  let editingNoteId = $state<string | null>(null);
  let editContent = $state('');
  let updating = $state(false);

  async function loadNotes() {
    if (!candidateId) return;
    
    try {
      loading = true;
      error = null;
      notes = await NoteService.getNotesByCandidate(candidateId);
    } catch (err) {
      console.error('Failed to load notes:', err);
      error = 'Failed to load notes. Please try again.';
    } finally {
      loading = false;
    }
  }

  async function addNote() {
    if (!newNoteContent.trim()) return;

    try {
      submitting = true;
      const newNote = await NoteService.createNote({
        candidateId,
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
    editingNoteId = note.id;
    editContent = note.content;
  }

  function cancelEditing() {
    editingNoteId = null;
    editContent = '';
  }

  async function saveEdit(noteId: string) {
    if (!editContent.trim()) return;

    try {
      updating = true;
      const updatedNote = await NoteService.updateNote(noteId, {
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

  async function deleteNote(noteId: string) {
    if (!confirm('Are you sure you want to delete this note?')) return;

    try {
      await NoteService.deleteNote(noteId);
      notes = notes.filter(note => note.id !== noteId);
      toast.success('Note deleted successfully');
    } catch (err) {
      console.error('Failed to delete note:', err);
      toast.error('Failed to delete note. Please try again.');
    }
  }

  function handleClose() {
    newNoteContent = '';
    error = null;
    onClose();
  }

  // Load notes when modal opens
  $effect(() => {
    if (isOpen && candidateId) {
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

    <!-- Error State -->
    {#if error}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2">
        <CircleAlert class="w-4 h-4 text-red-500" />
        <span class="text-sm text-red-700">{error}</span>
      </div>
    {/if}

    <!-- Notes List -->
    <div class="space-y-4 max-h-62 overflow-y-auto mb-6">
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
          <div class="bg-gray-100 rounded-lg p-4">
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
                      class="flex items-center px-3 py-1.5 text-xs bg-white border border-gray-200 rounded hover:bg-gray-50"
                      onclick={cancelEditing}
                      disabled={updating}
                    >
                      Cancel
                    </button>
                    <button
                      class="flex items-center px-3 py-1.5 text-xs bg-green-600 font-medium text-white rounded hover:bg-green-700 disabled:opacity-50"
                      disabled={!editContent.trim() || updating}
                      onclick={() => saveEdit(note.id)}
                    >
                      Save
                    </button>
            
                  </div>
                </div>
              </div>
            {:else}
              <!-- View mode -->
              <div class="group">
                <div class="flex items-start justify-between">
                  <p class="text-sm text-gray-700 leading-relaxed flex-1 pr-2">{note.content}</p>
                  <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                      class="p-1 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded"
                      onclick={() => startEditing(note)}
                      title="Edit note"
                    >
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button
                      class="p-1 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded"
                      onclick={() => deleteNote(note.id)}
                      title="Delete note"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <div class="mt-2 flex items-center justify-between">
                  <div class="text-xs text-gray-400">
                    <span>{formatRelativeTime(note.createdAt)}</span>
                    {#if note.updatedAt !== note.createdAt}
                      <span> • Edited {formatRelativeTime(note.updatedAt)}</span>
                    {/if}
                  </div>
                </div>
              </div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>

    <!-- Add Note Form -->
    <div class="mb-6">
      <p class="text-sm font-medium text-gray-600 mb-2">Add a note</p>
      <div class="border border-gray-200 rounded-lg">
        
        <textarea
          bind:value={newNoteContent}
          placeholder="Enter your notes about this candidate..."
          rows="3"
          maxlength="1000"
          class="w-full p-3 border-0 rounded-t-lg text-sm resize-none focus:outline-none focus:ring-1 focus:ring-gray-400"
          disabled={submitting}
          onkeydown={handleKeydown}
        ></textarea>
        <div class="flex items-center justify-between p-3 bg-gray-50 border-t border-gray-200 rounded-b-lg">
          <span class="text-xs text-gray-500">
            {newNoteContent.length}/1000 • Press Cmd+Enter to submit
          </span>
          <button
            class="flex items-center px-3 py-1.5 text-sm bg-green-600 font-medium text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={!newNoteContent.trim() || submitting}
            onclick={addNote}
          >
            <Plus class="w-4 h-4 mr-1" />
            {submitting ? 'Adding...' : 'Add Note'}
          </button>
        </div>
      </div>
    </div>
  </div>
</Modal>