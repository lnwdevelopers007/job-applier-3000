<script lang="ts">
  import Modal from '../ui/Modal.svelte';
  import { FileText } from 'lucide-svelte';
  import { getUserInfo } from '$lib/utils/auth';
  import { toast } from 'svelte-french-toast';
  
  interface Job {
    id: string;
    title: string;
    company: string;
    location: string;
    workType?: string;
    workArrangement?: string;
  }
  
  interface Document {
    id: string;
    name: string;
    type: string;
    uploadDate: string;
  }
  
  let {
    isOpen = $bindable(false),
    job
  }: {
    isOpen: boolean;
    job: Job;
  } = $props();
  
  let documents = $state<Document[]>([]);
  let isSubmitting = $state(false);
  const userInfo = getUserInfo();
  const userId = userInfo?.userID || '';
  
  async function loadUserDocuments() {
    // TODO: Replace with actual API call to fetch user's documents
    documents = [];
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
      
      // Check if already applied
      const queryParams = new URLSearchParams({ 
        applicantID: userId,
        jobID: job.id 
      });
      
      const checkResponse = await fetch(`/apply/query?${queryParams.toString()}`, {
        credentials: 'include'
      });
      
      if (checkResponse.ok) {
        const applications = await checkResponse.json();
        if (applications.length > 0) {
          toast.error('You have already applied to this job');
          return;
        }
      }

      // Prepare application data matching the server format
      const applicationData = {
        applicantID: userId,
        jobID: job.id,
        status: 'PENDING',
        createdAt: new Date().toISOString(),
        documents: documents.map(doc => doc.id)
      };
      
      // Submit application
      const response = await fetch('/apply', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(applicationData)
      });
      
      if (!response.ok) {
        if (response.status === 409) {
          toast.error('You have already applied to this job');
        } else {
          throw new Error(`Failed to apply: ${response.status}`);
        }
        return;
      }
      
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
</script>

<Modal bind:isOpen onClose={handleClose} size="lg" closeOnBackdrop={false}>
  <div class="p-6">
    <!-- Application Form -->
    <div>
      <!-- Header -->
      <h2 class="text-xl font-medium text-gray-900 mb-2">Apply for this position</h2>
      <p class="text-sm text-gray-600 mb-4">Submit your application to {job.title} at {job.company}</p>
      
      <!-- Documents Section -->
      {#if documents.length === 0}
        <div class="mb-6 bg-gray-100 p-3 rounded-md">
          <p class="text-sm text-gray-600">
            No documents found. You can still apply or 
            <button 
              type="button"
              class="text-green-600 underline hover:text-green-700 hover:cursor-pointer"
              onclick={() => window.location.href = '/app/settings'}
            >
              add documents
            </button> 
            to improve your chances.
          </p>
        </div>
      {:else}
        <div class="mb-6">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium text-gray-900">Documents to be sent</h3>
            <button 
              type="button"
              class="text-sm text-green-600 underline hover:text-green-700"
              onclick={() => window.location.href = '/app/settings'}
            >
              Edit in settings
            </button>
          </div>
          <div class="space-y-2">
            {#each documents as doc}
              <div class="flex items-center gap-3 p-3 border border-gray-200 rounded-lg bg-gray-50">
                <FileText class="w-5 h-5 text-gray-400" />
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">{doc.name}</p>
                  <p class="text-xs text-gray-500">Uploaded {doc.uploadDate}</p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
      
      
      <!-- Actions -->
      <div class="flex justify-center">
        <button
          type="button"
          onclick={handleSubmit}
          disabled={isSubmitting}
          class="flex-1 items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-green-600 border border-transparent rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed hover:cursor-pointer"
        >
          {#if isSubmitting}
            Submitting...
          {:else}
            Submit Application
          {/if}
        </button>
      </div>
    </div>
  </div>
</Modal>