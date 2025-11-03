<script lang="ts">
  // client/src/lib/components/files/DeleteConfirmModal.svelte
  import Modal from '../ui/Modal.svelte';
  import { AlertTriangle, Loader } from 'lucide-svelte';
  
  let {
    isOpen = $bindable(false),
    filename = '',
    isDeleting = false,
    onConfirm = () => {},
    onCancel = () => {}
  } = $props();
  
  function handleConfirm() {
    onConfirm();
  }
  
  function handleCancel() {
    if (!isDeleting) {
      onCancel();
      isOpen = false;
    }
  }
</script>

<Modal bind:isOpen size="sm" onClose={handleCancel} closeOnBackdrop={!isDeleting}>
  <div class="p-6">
    <!-- Icon -->
    <div class="flex items-center justify-center w-12 h-12 mx-auto mb-4 bg-red-100 rounded-full">
      <AlertTriangle class="w-6 h-6 text-red-600" />
    </div>
    
    <!-- Title -->
    <h3 class="text-lg font-medium text-gray-900 text-center mb-2">
      Delete Document
    </h3>
    
    <!-- Message -->
    <div class="mb-6 text-center">
      <p class="text-sm text-gray-600 mb-2">
        Are you sure you want to delete
      </p>
      <p class="text-sm font-medium text-gray-900 mb-2">
        "{filename}"?
      </p>
      <p class="text-sm text-red-600">
        This action cannot be undone.
      </p>
    </div>
    
    <!-- Actions -->
    <div class="flex items-center gap-3">
      <button
        onclick={handleCancel}
        disabled={isDeleting}
        class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors"
      >
        Cancel
      </button>
      <button
        onclick={handleConfirm}
        disabled={isDeleting}
        class="flex-1 px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
      >
        {#if isDeleting}
          <Loader class="w-4 h-4 animate-spin" />
          Deleting...
        {:else}
          Delete
        {/if}
      </button>
    </div>
  </div>
</Modal>