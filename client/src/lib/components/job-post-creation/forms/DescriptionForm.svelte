<script>
  import TiptapEditor from '$lib/components/job-post-creation/TextEditor.svelte';
  import { AlertCircle } from 'lucide-svelte';
  
  let { 
    formData = $bindable(),
    validationErrors = {},
    showValidationErrors = false
  } = $props();
  
  // Track if user has started typing to hide error
  let hideJobSummaryError = $state(false);
  // Track hover state for error icon
  let isJobSummaryErrorHovered = $state(false);
  
  // Hide error when user starts typing
  function handleJobSummaryInput() {
    hideJobSummaryError = true;
  }
  
  // Reset hideError when new validation occurs
  $effect(() => {
    if (showValidationErrors && validationErrors.jobSummary) {
      hideJobSummaryError = false;
    }
  });
</script>

<div class="space-y-6">
  
  <!-- Job Description -->
  <div>
    <label for="jobDescription" class="block text-sm font-medium text-gray-700 mb-2">
      Job Description <span class="text-red-500">*</span>
    </label>
    <TiptapEditor 
      bind:value={formData.jobDescription}
      placeholder="Describe the job role and responsibilities..."
      height="450px"
      error={validationErrors.jobDescription}
      showError={showValidationErrors}
    />
  </div>

  <!-- Job Summary -->
  <div class="relative">
    <label for="jobSummary" class="block text-sm font-medium text-gray-700 mb-2">
      Job Summary <span class="text-red-500">*</span>
    </label>
    <textarea
      id="jobSummary"
      bind:value={formData.jobSummary}
      rows="4"
      placeholder="Brief summary of the position..."
      oninput={handleJobSummaryInput}
      class="w-full text-sm px-3 py-2 border {showValidationErrors && validationErrors.jobSummary && !hideJobSummaryError ? 'border-red-500' : 'border-gray-300'} rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 resize-none"
    ></textarea>
    
    {#if showValidationErrors && validationErrors.jobSummary && !hideJobSummaryError}
      <!-- Error icon positioned at top-right with hover detection -->
      <div 
        class="absolute top-8 right-2"
        onmouseenter={() => isJobSummaryErrorHovered = true}
        onmouseleave={() => isJobSummaryErrorHovered = false}
        role="img"
        aria-label="Error"
      >
        <AlertCircle class="w-5 h-5 text-white cursor-help" fill="red" />
      </div>
      
      <!-- Floating error tooltip - only show on hover -->
      {#if isJobSummaryErrorHovered}
        <div class="absolute z-50 right-0 top-8 -translate-y-full -mt-2 px-3 py-2 bg-red-500 text-white text-xs rounded-md shadow-sm whitespace-nowrap transition-opacity duration-200">
          <!-- Arrow pointing down to the alert icon -->
          <div class="absolute -bottom-1 right-4 w-2 h-2 bg-red-500 transform rotate-45"></div>
          {validationErrors.jobSummary}
        </div>
      {/if}
    {/if}
  </div>
</div>