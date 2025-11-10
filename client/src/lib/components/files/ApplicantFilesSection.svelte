<script lang="ts">
  // client/src/lib/components/files/ApplicantFilesSection.svelte
  import FilePreviewModal from './FilePreviewModal.svelte';
  import { FileText, Eye, Download, Loader, Loader2 } from 'lucide-svelte';
  import { fileService, type FileMetadata } from '$lib/services/fileService';
  import { toast } from 'svelte-french-toast';
  
  let {
    applicationId
  } = $props();
  
  let files = $state<FileMetadata[]>([]);
  let isLoading = $state(false);
  let error = $state<string | null>(null);
  let isPreviewModalOpen = $state(false);
  let selectedFile = $state<FileMetadata | null>(null);
  
  async function loadApplicantFiles() {
    if (!applicationId) return;
    
    isLoading = true;
    error = null;
    
    try {
      const response = await fileService.getApplicantFiles(applicationId);
      files = response.files || [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load applicant documents';
      console.error('Failed to load applicant files:', err);
      toast.error('Failed to load applicant documents');
    } finally {
      isLoading = false;
    }
  }
  
  function handlePreview(file: FileMetadata) {
    selectedFile = file;
    isPreviewModalOpen = true;
  }
  
  async function handleDownload(file: FileMetadata) {
    try {
      const blob = await fileService.downloadApplicantFile(applicationId, file.id);
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = file.filename;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      console.error('Download error:', err);
      toast.error('Failed to download file');
    }
  }
  
  // Load files when component mounts or applicationId changes
  $effect(() => {
    loadApplicantFiles();
  });
  
  const getCategoryColor = (category: string) => fileService.getCategoryColor(category);
  const getCategoryLabel = (category: string) => fileService.getCategoryLabel(category);
  const formatFileSize = (size: number) => fileService.formatFileSize(size);
  const formatDate = (dateString: string) => new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  });
</script>

<div class="mb-2">
  <h3 class="flex gap-1 font-medium text-gray-800 text-md mt-2">
    Documents
  </h3>
  
  <div class="mt-3">
    {#if isLoading}
      <div class="flex items-center justify-center py-8">
        <Loader2 class="w-6 h-6 text-gray-400 animate-spin" />
      </div>
    {:else if error}
      <div class="bg-red-50 border border-red-200 p-4 rounded-lg">
        <p class="text-sm text-red-700">{error}</p>
      </div>
    {:else if files.length === 0}
      <div class="pt-2 text-center">
        <p class="text-sm text-gray-500">No documents uploaded by this applicant</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each files as file (file.id)}
          <div class="flex items-center justify-between bg-white border border-gray-200 rounded-lg p-4 hover:bg-gray-100 transition-colors group">
            <div class="flex items-center gap-4 flex-1 min-w-0">
              <!-- Icon -->
              <div class="w-12 h-12 bg-{getCategoryColor(file.category)}-100 rounded-lg flex items-center justify-center flex-shrink-0">
                <FileText class="w-6 h-6 text-{getCategoryColor(file.category)}-600" />
              </div>
              
              <!-- Info -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <p class="text-sm font-semibold text-gray-800 truncate">{file.filename}</p>
                  <span class="px-2 py-0.5 text-xs font-medium bg-{getCategoryColor(file.category)}-100 text-{getCategoryColor(file.category)}-700 rounded">
                    {getCategoryLabel(file.category)}
                  </span>
                </div>
                <p class="text-xs text-gray-500">
                  {formatFileSize(file.size)} • Uploaded {formatDate(file.uploadDate)}
                </p>
              </div>
            </div>
            
            <!-- Actions -->
            <div class="flex items-center gap-2">
              <button
                onclick={() => handlePreview(file)}
                class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                title="Preview document"
              >
                <Eye class="w-5 h-5" />
              </button>
              <button
                onclick={() => handleDownload(file)}
                class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
                title="Download document"
              >
                <Download class="w-5 h-5" />
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<!-- Preview Modal -->
{#if selectedFile}
  <FilePreviewModal
    bind:isOpen={isPreviewModalOpen}
    fileId={selectedFile.id}
    filename={selectedFile.filename}
    applicationId={applicationId}
  />
{/if}