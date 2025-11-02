<script lang="ts">
  import Modal from './Modal.svelte';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    onConfirm: (reason: string) => Promise<void>;
    title: string;
    itemName: string;
    description?: string;
    reasonLabel?: string;
    reasonPlaceholder?: string;
    confirmButtonText?: string;
    isDeleting?: boolean;
  }

  let {
    isOpen = $bindable(),
    onClose,
    onConfirm,
    title,
    itemName,
    description,
    reasonLabel = "Reason for deletion",
    reasonPlaceholder = "Please provide a detailed reason...",
    confirmButtonText = "Delete",
    isDeleting = false
  }: Props = $props();

  let deleteReason = $state('');

  async function handleConfirm() {
    if (!deleteReason.trim()) return;
    await onConfirm(deleteReason);
    deleteReason = ''; // Reset after successful deletion
  }

  function handleClose() {
    if (!isDeleting) {
      deleteReason = '';
      onClose();
    }
  }

  // Reset reason when modal closes
  $effect(() => {
    if (!isOpen) {
      deleteReason = '';
    }
  });
</script>

<Modal bind:isOpen={isOpen} size="md" closeOnBackdrop={false} onClose={handleClose}>
  <div class="p-6">
    <!-- Header -->
    <div class="flex items-start space-x-4 mb-6">
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-medium text-gray-900 mb-2">
          {title}
        </h3>
        <p class="text-sm text-gray-600">
          {#if description}
            {description}
          {:else}
            You're about to delete <span class="font-medium text-gray-900">"{itemName}"</span>. This action cannot be undone.
          {/if}
        </p>
      </div>
    </div>

    <!-- Reason input -->
    <div class="mb-6">
      <label for="delete-reason" class="block text-sm font-medium text-gray-700 mb-3">
        {reasonLabel} <span class="text-red-500">*</span>
      </label>
      <textarea
        id="delete-reason"
        bind:value={deleteReason}
        rows="4"
        maxlength="500"
        placeholder={reasonPlaceholder}
        class="w-full border border-gray-300 rounded-lg text-sm px-3 py-3 text-gray-800 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent resize-none transition-all"
        disabled={isDeleting}
      ></textarea>
      <div class="mt-2 flex justify-between items-center">
        <span class="text-xs text-gray-500">This reason will be logged for audit purposes</span>
        <span class="text-xs {deleteReason.length >= 480 ? 'text-red-500 font-medium' : 'text-gray-500'}">
          {deleteReason.length}/500
        </span>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-end space-x-3">
      <button
        class="px-4 py-2 bg-white border border-gray-300 text-gray-700 text-sm rounded-lg hover:bg-gray-50 font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
        onclick={handleClose}
        disabled={isDeleting}
      >
        Cancel
      </button>
      <button
        class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 text-sm disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 flex items-center"
        disabled={!deleteReason.trim() || isDeleting}
        onclick={handleConfirm}
      >
        {#if isDeleting}
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Deleting...
        {:else}
          {confirmButtonText}
        {/if}
      </button>
    </div>
  </div>
</Modal>