<script>
  import TiptapEditor from '$lib/components/TextEditor.svelte';
  
  let { formData = $bindable() } = $props();
</script>

<div class="space-y-6">
  <h2 class="text-lg font-medium text-gray-900 mb-6">Job Description</h2>
  
  <!-- Job Description -->
  <div>
    <label for="jobDescription" class="block text-sm font-medium text-gray-700 mb-2">
      Job Description <span class="text-red-500">*</span>
    </label>
    <TiptapEditor 
      bind:value={formData.jobDescription}
      placeholder="Describe the job role and responsibilities..."
    />
  </div>

  <!-- Job Summary -->
  <div>
    <label for="jobSummary" class="block text-sm font-medium text-gray-700 mb-2">
      Job Summary
    </label>
    <textarea
      id="jobSummary"
      bind:value={formData.jobSummary}
      rows="4"
      placeholder="Brief summary of the position..."
      class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 resize-none"
    ></textarea>
  </div>
</div>