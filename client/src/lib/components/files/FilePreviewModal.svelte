<script lang="ts">
  import Modal from '../ui/Modal.svelte';
  import { X, Download, Loader } from 'lucide-svelte';
  import { fileService } from '$lib/services/fileService';
  
  let {
    isOpen = $bindable(false),
    fileId = '',
    filename = 'Document'
  } = $props();
  
  let previewUrl = $state<string | null>(null);
  let isLoading = $state(false);
  let error = $state<string | null>(null);
  
  async function loadPreview() {
    if (!fileId) return;
    
    isLoading = true;
    error = null;
    
    try {
      previewUrl = await fileService.getFilePreviewUrl(fileId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load preview';
      console.error('Preview error:', err);
    } finally {
      isLoading = false;
    }
  }
  
  async function handleDownload() {
    try {
      const blob = await fileService.downloadFile(fileId);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Download error:', err);
    }
  }
  
  function handleClose() {
    isOpen = false;
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
      previewUrl = null;
    }
  }
  
  // Load preview when modal opens
  $effect(() => {
    if (isOpen && fileId) {
      loadPreview();
    }
  });
  
  // Cleanup on unmount
  $effect(() => {
    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
    };
  });
</script>

<Modal bind:isOpen size="full" onClose={handleClose}>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center justify-between p-4 border-b border-gray-200 bg-white">
      <h3 class="text-lg font-medium text-gray-900 truncate">{filename}</h3>
      <div class="flex items-center gap-2">
        <button
          onclick={handleDownload}
          class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
          title="Download"
        >
          <Download class="w-5 h-5" />
        </button>
        <button
          onclick={handleClose}
          class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
          title="Close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>
    </div>
    
    <!-- Content -->
    <div class="flex-1 bg-gray-100 overflow-hidden">
      {#if isLoading}
        <div class="flex items-center justify-center h-full">
          <div class="text-center">
            <Loader class="w-12 h-12 text-gray-400 animate-spin mx-auto mb-4" />
            <p class="text-sm text-gray-600">Loading preview...</p>
          </div>
        </div>
      {:else if error}
        <div class="flex items-center justify-center h-full">
          <div class="text-center max-w-md">
            <svg class="w-16 h-16 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="text-sm text-gray-900 font-medium mb-2">Failed to load preview</p>
            <p class="text-sm text-gray-600 mb-4">{error}</p>
            <button
              onclick={handleDownload}
              class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
            >
              Download File Instead
            </button>
          </div>
        </div>
      {:else if previewUrl}
        <embed
          src={previewUrl}
          type="application/pdf"
          class="w-full h-full"
          title={filename}
        />
      {/if}
    </div>
  </div>
</Modal>