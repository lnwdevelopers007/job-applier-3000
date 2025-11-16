<script lang="ts">
  // client/src/lib/components/job/ApplyModal.svelte
  import Modal from '../ui/Modal.svelte';
  import FilePreviewModal from '../files/FilePreviewModal.svelte';
  import { FileText, Eye, Loader } from 'lucide-svelte';
  import { getUserInfo } from '$lib/utils/auth';
  import { toast } from 'svelte-french-toast';
  import { fileService, type FileMetadata } from '$lib/services/fileService';
  import { JobApplicationService } from '$lib/services/jobApplicationService';
  
  interface Job {
    id: string;
    title: string;
    company: string;
  }
  
  let {
    isOpen = $bindable(false),
    job
  }: {
    isOpen: boolean;
    job: Job;
  } = $props();
  
  let documents = $state<FileMetadata[]>([]);
  let isLoadingDocs = $state(false);
  let isSubmitting = $state(false);
  let isPreviewModalOpen = $state(false);
  let selectedFile = $state<FileMetadata | null>(null);
  
  const userInfo = getUserInfo();
  const userId = userInfo?.userID || '';
  
  async function loadUserDocuments() {
    if (!userId) return;
    
    isLoadingDocs = true;
    try {
      const allFiles = await fileService.listUserFiles(userId);
      // Filter to show only job seeker relevant documents
      documents = allFiles.filter(file => 
        ['resume', 'transcript', 'certification'].includes(file.category)
      );
    } catch (err) {
      console.error('Failed to load documents:', err);
      toast.error('Failed to load your documents');
    } finally {
      isLoadingDocs = false;
    }
  }
  
  function handlePreview(file: FileMetadata) {
    selectedFile = file;
    isPreviewModalOpen = true;
  }
  
  function navigateToSettings() {
    window.location.href = '/app/settings';
  }
  
  async function handleSubmit() {
    if (isSubmitting) return;
    
    isSubmitting = true;
    
    try {
      // Check if user is authenticated
      if (!userId) {
        toast.error('You must be logged in to apply');
        return;
      }
      
      // Check if already applied using the service
      const hasApplied = await JobApplicationService.hasUserAppliedToJob(userId, job.id);
      if (hasApplied) {
        toast.error('You have already applied to this job');
        return;
      }

      // Prepare application data  
      const applicationData = {
        applicantID: userId,
        jobID: job.id,
        status: 'PENDING',
        createdAt: new Date().toISOString()
      };
      
      
      // Submit application using the service
      await JobApplicationService.createApplication(applicationData);
      
      // Show success toast and close modal
      toast.success(`Successfully applied to ${job.title}`);
      isOpen = false;
      
    } catch (error) {
      console.error('Error submitting application:', error);
      toast.error('Failed to apply. Please try again.');
    } finally {
      isSubmitting = false;
    }
  }
  
  function handleClose() {
    if (!isSubmitting) {
      isOpen = false;
    }
  }
  
  // Load documents when modal opens
  $effect(() => {
    if (isOpen) {
      loadUserDocuments();
    }
  });
  
  const getCategoryColor = (category: string) => fileService.getCategoryColor(category);
  const getCategoryLabel = (category: string) => fileService.getCategoryLabel(category);
  const formatFileSize = (size: number) => fileService.formatFileSize(size);
</script>

<Modal bind:isOpen onClose={handleClose} size="lg" closeOnBackdrop={false}>
  <div class="p-6">
    <!-- Header -->
    <h2 class="text-xl font-medium text-gray-900 mb-2">Apply for this position</h2>
    <p class="text-sm text-gray-600 mb-6">Submit your application to {job.title} at {job.company}</p>
    
    <!-- Documents Section -->
    <div class="mb-6">
      <div class="flex items-center justify-between mb-3">
        <h3 class="font-medium text-gray-900">Your Documents</h3>
        <button 
          type="button"
          onclick={navigateToSettings}
          class="text-sm text-green-600 hover:text-green-700 font-medium"
        >
          Manage in settings →
        </button>
      </div>
      
      {#if isLoadingDocs}
        <div class="flex items-center justify-center py-8">
          <Loader class="w-6 h-6 text-gray-400 animate-spin" />
        </div>
      {:else if documents.length === 0}
        <div class="bg-yellow-50 border border-yellow-200 p-4 rounded-md">
          <p class="text-sm text-yellow-800 mb-2">
            You haven't uploaded any documents yet.
          </p>
          <p class="text-sm text-yellow-700 mb-3">
            While you can still apply, uploading documents will improve your chances.
          </p>
          <button 
            type="button"
            onclick={navigateToSettings}
            class="text-sm text-yellow-900 font-medium hover:text-yellow-800"
          >
            Go to Settings to Upload Documents →
          </button>
        </div>
      {:else}
        <div class="space-y-2">
          {#each documents as doc (doc.id)}
            <div class="flex items-center justify-between p-3 border border-gray-200 rounded-lg bg-gray-50">
              <div class="flex items-center gap-3 flex-1 min-w-0">
                <div class="w-10 h-10 bg-{getCategoryColor(doc.category)}-100 rounded flex items-center justify-center flex-shrink-0">
                  <FileText class="w-5 h-5 text-{getCategoryColor(doc.category)}-600" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-0.5">
                    <p class="text-sm font-medium text-gray-900 truncate">{doc.filename}</p>
                    <span class="px-2 py-0.5 text-xs font-medium bg-{getCategoryColor(doc.category)}-100 text-{getCategoryColor(doc.category)}-700 rounded">
                      {getCategoryLabel(doc.category)}
                    </span>
                  </div>
                  <p class="text-xs text-gray-500">{formatFileSize(doc.size)}</p>
                </div>
              </div>
              <button
                type="button"
                onclick={() => handlePreview(doc)}
                class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors flex-shrink-0"
                title="Preview"
              >
                <Eye class="w-5 h-5" />
              </button>
            </div>
          {/each}
        </div>
        <p class="text-xs text-gray-500 mt-2">
          These documents will be automatically sent with your application
        </p>
      {/if}
    </div>
    
    <!-- Actions -->
    <div class="flex justify-end gap-3">
      <button
        type="button"
        onclick={handleClose}
        disabled={isSubmitting}
        class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors"
      >
        Cancel
      </button>
      <button
        type="button"
        onclick={handleSubmit}
        disabled={isSubmitting}
        class="px-6 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:opacity-50 transition-colors flex items-center gap-2"
      >
        {#if isSubmitting}
          <Loader class="w-4 h-4 animate-spin" />
          Submitting...
        {:else}
          Submit Application
        {/if}
      </button>
    </div>
  </div>
</Modal>

<!-- Preview Modal -->
{#if selectedFile}
  <FilePreviewModal
    bind:isOpen={isPreviewModalOpen}
    fileId={selectedFile.id}
    filename={selectedFile.filename}
  />
{/if}