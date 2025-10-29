<script lang="ts">
  // client/src/lib/components/files/FileUploadModal.svelte
  import Modal from '../ui/Modal.svelte';
  import Select from '../forms/Select.svelte';
  import { FileText, Upload, Loader } from 'lucide-svelte';
  import { fileService } from '$lib/services/fileService';
  
  let {
    isOpen = $bindable(false),
    userRole = 'jobSeeker',
    onUploadSuccess = (file: any) => {}
  } = $props();
  
  let selectedFile = $state<File | null>(null);
  let selectedCategory = $state('');
  let isUploading = $state(false);
  let error = $state<string | null>(null);
  let dragOver = $state(false);
  
  // Get valid categories based on user role
  const categories = $derived(fileService.getValidCategories(userRole));
  
  // Set default category when role changes
  $effect(() => {
    if (categories.length > 0 && !selectedCategory) {
      selectedCategory = categories[0].value;
    }
  });
  
  function handleFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files[0]) {
      validateAndSetFile(target.files[0]);
    }
  }
  
  function handleDrop(event: DragEvent) {
    event.preventDefault();
    dragOver = false;
    
    if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
      validateAndSetFile(event.dataTransfer.files[0]);
    }
  }
  
  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    dragOver = true;
  }
  
  function handleDragLeave() {
    dragOver = false;
  }
  
  function validateAndSetFile(file: File) {
    error = null;
    
    // Validate file type
    if (file.type !== 'application/pdf') {
      error = 'Only PDF files are allowed';
      return;
    }
    
    // Validate file size (10MB)
    if (file.size > 10 * 1024 * 1024) {
      error = 'File size must be less than 10MB';
      return;
    }
    
    selectedFile = file;
  }
  
  function removeFile() {
    selectedFile = null;
    error = null;
  }
  
  async function handleUpload() {
    if (!selectedFile || !selectedCategory) return;
    
    isUploading = true;
    error = null;
    
    try {
      const result = await fileService.uploadFile(selectedFile, selectedCategory);
      onUploadSuccess(result);
      handleClose();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to upload file';
    } finally {
      isUploading = false;
    }
  }
  
  function handleClose() {
    if (isUploading) return;
    
    isOpen = false;
    selectedFile = null;
    error = null;
    
    // Reset category to default
    if (categories.length > 0) {
      selectedCategory = categories[0].value;
    }
  }
</script>

<Modal bind:isOpen size="md" onClose={handleClose} closeOnBackdrop={!isUploading}>
  <div class="p-6">
    <h2 class="text-xl font-medium text-gray-900 mb-4">Upload Document</h2>
    
    <!-- Category Selection -->
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700 mb-2">
        Document Category
      </label>
      <Select
        bind:value={selectedCategory}
        options={categories}
        placeholder="Select category..."
        disabled={isUploading}
      />
    </div>
    
    <!-- File Upload Area -->
    <div class="mb-4">
      <label for="file-input" class="block text-sm font-medium text-gray-700 mb-2">
        Select File
        </label>
        <input
        id="file-input"
        type="file"
        accept="application/pdf"
        class="hidden"
        onchange={handleFileSelect}
        disabled={isUploading}
        />
      
      {#if !selectedFile}
        <div
        role="button"
        tabindex="0"
        aria-label="File drop area"
        class="border-2 border-dashed rounded-lg p-8 text-center transition-colors {dragOver ? 'border-green-500 bg-green-50' : 'border-gray-300 hover:border-gray-400'}"
        ondrop={handleDrop}
        ondragover={handleDragOver}
        ondragleave={handleDragLeave}
        >
          <Upload class="w-12 h-12 text-gray-400 mx-auto mb-3" />
          <p class="text-sm text-gray-600 mb-2">
            Drag and drop your PDF here, or
          </p>
          <label class="inline-block px-4 py-2 bg-white border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 transition-colors cursor-pointer">
            Browse Files
            <input
              type="file"
              accept="application/pdf"
              class="hidden"
              onchange={handleFileSelect}
              disabled={isUploading}
            />
          </label>
          <p class="text-xs text-gray-500 mt-3">PDF only, max 10MB</p>
        </div>
      {:else}
        <div class="border border-gray-200 rounded-lg p-4 bg-gray-50">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-red-100 rounded flex items-center justify-center flex-shrink-0">
              <FileText class="w-5 h-5 text-red-600" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">{selectedFile.name}</p>
              <p class="text-xs text-gray-500">{fileService.formatFileSize(selectedFile.size)}</p>
            </div>
            {#if !isUploading}
              <button
                onclick={removeFile}
                class="text-gray-400 hover:text-red-600 transition-colors"
                aria-label="Remove file"
                >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                </button>
            {/if}
          </div>
        </div>
      {/if}
    </div>
    
    <!-- Error Message -->
    {#if error}
      <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md">
        <p class="text-sm text-red-700">{error}</p>
      </div>
    {/if}
    
    <!-- Actions -->
    <div class="flex items-center justify-end gap-3">
      <button
        onclick={handleClose}
        disabled={isUploading}
        class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors"
      >
        Cancel
      </button>
      <button
        onclick={handleUpload}
        disabled={!selectedFile || !selectedCategory || isUploading}
        class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:opacity-50 transition-colors flex items-center gap-2"
      >
        {#if isUploading}
          <Loader class="w-4 h-4 animate-spin" />
          Uploading...
        {:else}
          Upload Document
        {/if}
      </button>
    </div>
  </div>
</Modal>