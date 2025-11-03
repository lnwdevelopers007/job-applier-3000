<script lang="ts">
  import FileItem from '$lib/components/files/FileItem.svelte';
  import FilePreviewModal from '$lib/components/files/FilePreviewModal.svelte';
  import FileUploadModal from '$lib/components/files/FileUploadModal.svelte';
  import DeleteConfirmModal from '$lib/components/files/DeleteConfirmModal.svelte';
  import { fileService, type FileMetadata } from '$lib/services/fileService';
  import { getUserInfo } from '$lib/utils/auth';
  import { toast } from 'svelte-french-toast';
  import { Loader } from 'lucide-svelte';
  
  let {
    documents = $bindable([])
  } = $props();
  
  let files = $state<FileMetadata[]>([]);
  let isLoading = $state(true);
  let isUploadModalOpen = $state(false);
  let isPreviewModalOpen = $state(false);
  let isDeleteModalOpen = $state(false);
  let isDeleting = $state(false);
  
  let selectedFile = $state<FileMetadata | null>(null);
  let fileToDelete = $state<FileMetadata | null>(null);
  
  const userInfo = getUserInfo();
  const userId = userInfo?.userID || '';
  const userRole = userInfo?.role || 'jobSeeker';
  
  async function loadFiles() {
    if (!userId) {
      console.error('User ID not found');
      return;
    }
    
    isLoading = true;
    try {
      files = await fileService.listUserFiles(userId);
      // Update parent documents binding
      documents = files;
    } catch (err) {
      console.error('Failed to load files:', err);
      toast.error('Failed to load documents');
    } finally {
      isLoading = false;
    }
  }
  
  function handleUploadClick() {
    isUploadModalOpen = true;
  }
  
  async function handleUploadSuccess(file: FileMetadata) {
    // Reload files from server to ensure we have the latest data
    await loadFiles();
    toast.success(`${fileService.getCategoryLabel(file.category)} uploaded successfully`);
  }
  
  function handlePreview(file: FileMetadata) {
    selectedFile = file;
    isPreviewModalOpen = true;
  }
  
  function handleDeleteClick(file: FileMetadata) {
    fileToDelete = file;
    isDeleteModalOpen = true;
  }
  
  async function handleDeleteConfirm() {
    if (!fileToDelete) return;
    
    isDeleting = true;
    try {
      await fileService.deleteFile(fileToDelete.id);
      files = files.filter(f => f.id !== fileToDelete.id);
      documents = files;
      toast.success('Document deleted successfully');
      isDeleteModalOpen = false;
      fileToDelete = null;
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete document');
    } finally {
      isDeleting = false;
    }
  }
  
  function handleDeleteCancel() {
    fileToDelete = null;
  }
  
  // Load files on mount
  $effect(() => {
    loadFiles();
  });
    $effect(() => {
    const handleOpenModal = () => {
      isUploadModalOpen = true;
    };
    
    window.addEventListener('openUploadModal', handleOpenModal);
    
    return () => {
      window.removeEventListener('openUploadModal', handleOpenModal);
    };
  });
  
  // Expose upload trigger for parent component
  $effect(() => {
    const uploadInput = document.getElementById('file-upload');
    if (uploadInput) {
      uploadInput.addEventListener('click', handleUploadClick);
      return () => {
        uploadInput.removeEventListener('click', handleUploadClick);
      };
    }
  });
</script>

<div class="divide-y divide-gray-200">
  <div class="py-5">
    <div class="grid grid-cols-4 gap-8 items-start">
      <div>
       <label for="file-upload" class="text-sm font-medium text-gray-700">Documents</label>
        <p class="text-xs text-gray-500 mt-1">
          Upload documents such as resume, transcript, or certificates (PDF only, max 10MB)
        </p>
      </div>
      <div class="col-span-2">
        {#if isLoading}
          <div class="flex items-center justify-center py-12">
            <Loader class="w-8 h-8 text-gray-400 animate-spin" />
          </div>
        {:else if files.length === 0}
          <div class="text-center py-12 border-2 border-dashed border-gray-200 rounded-lg">
            <svg class="w-12 h-12 text-gray-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
            <h3 class="text-sm font-medium text-gray-900 mb-1">No documents</h3>
            <p class="text-sm text-gray-500 mb-4">Get started by uploading your first document.</p>
          </div>
        {:else}
          <div class="space-y-3">
            {#each files as file (file.id)}
              <FileItem
                {file}
                onPreview={() => handlePreview(file)}
                onDelete={() => handleDeleteClick(file)}
              />
            {/each}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<!-- Hidden file input for parent component trigger -->
<input type="file" id="file-upload" class="hidden" />

<!-- Upload Modal -->
<FileUploadModal
  bind:isOpen={isUploadModalOpen}
  {userRole}
  onUploadSuccess={handleUploadSuccess}
/>

<!-- Preview Modal -->
{#if selectedFile}
  <FilePreviewModal
    bind:isOpen={isPreviewModalOpen}
    fileId={selectedFile.id}
    filename={selectedFile.filename}
  />
{/if}

<!-- Delete Confirmation Modal -->
{#if fileToDelete}
  <DeleteConfirmModal
    bind:isOpen={isDeleteModalOpen}
    filename={fileToDelete.filename}
    {isDeleting}
    onConfirm={handleDeleteConfirm}
    onCancel={handleDeleteCancel}
  />
{/if}
