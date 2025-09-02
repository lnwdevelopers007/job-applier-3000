<script>
  import { ArrowLeft, ArrowRight, ChevronLeft } from 'lucide-svelte';
  import Stepper from '$lib/components/Stepper.svelte';
  import BasicInfoForm from '$lib/components/forms/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/forms/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/forms/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/forms/PostSettingForm.svelte';
  import PreviewDrawer from '$lib/components/PreviewDrawer.svelte';
  
  let currentStep = $state(1); // Current active step (1-4)
  let isPreviewOpen = $state(false);
  /** @type {any} */
  let stepper;
  
  const steps = [
    { id: 1, name: 'Basic Info' },
    { id: 2, name: 'Description' },
    { id: 3, name: 'Requirement' },
    { id: 4, name: 'Post Setting' }
  ];

  let formData = $state({
    // Basic Info
    jobTitle: '',
    location: '',
    category: '',
    workType: 'full-time',
    workArrangement: 'on-site',
    currency: 'THB',
    minSalary: '',
    maxSalary: '',
    
    // Description
    jobDescription: '',
    companyDescription: '',
    
    // Requirements
    requiredSkills: '',
    experienceLevel: '',
    education: '',
    niceToHave: '',
    
    // Post Settings
    applicationDeadline: '',
    numberOfPositions: 1,
    visibility: 'public',
    emailNotifications: true,
    autoReject: false
  });

  function handleBackNavigation() {
    // Handle back navigation to previous page
  }

  function handleSaveDraft() {
    // Handle save draft functionality
  }

  function handlePreview() {
    isPreviewOpen = true;
  }

  const handleProgress = (stepIncrement) => {
    stepper.handleProgress(stepIncrement);
  }

</script>

<div class="max-w-4xl mx-auto py-8">
  <!-- Header -->
  <div class="px-8 mb-2">
    <div class="flex items-center gap-4">
      <button onclick={handleBackNavigation} class="p-1 hover:bg-gray-100 rounded">
        <ChevronLeft class="w-5 h-5 text-gray-600" />
      </button>
      <h1 class="text-xl font-medium text-gray-900">Create Job Post</h1>
    </div>
  </div>

  <Stepper {steps} bind:currentStep bind:this={stepper} />
  
  <form class="px-8 pb-8">
    {#if currentStep === 1}
      <BasicInfoForm bind:formData />
    {:else if currentStep === 2}
      <DescriptionForm bind:formData />
    {:else if currentStep === 3}
      <RequirementForm bind:formData />
    {:else if currentStep === 4}
      <PostSettingForm bind:formData />
    {/if}

    <!-- Action Buttons -->
    <div class="flex justify-between mt-8 pt-6">
      <div>
        {#if currentStep > 1}
          <button
            type="button"
            onclick={() => handleProgress(-1)}
            class="text-sm flex items-center px-4 py-2 text-gray-600 hover:text-gray-800 font-medium"
          >
            <ArrowLeft class="w-3 h-3 mr-1" />
            Back
          </button>
        {/if}
      </div>
      <div class="flex gap-3">
        <button
          type="button"
          onclick={handleSaveDraft}
          class="text-sm px-4 py-2 text-gray-600 hover:text-gray-800 font-medium"
        >
          Save Draft
        </button>
        <button
          type="button"
          onclick={handlePreview}
          class="text-sm px-4 py-2 border border-gray-300 text-gray-700 hover:bg-gray-50 rounded-md font-medium"
        >
          Preview
        </button>
        <button
          type="button"
          onclick={() => handleProgress(1)}
          class="text-sm flex items-center px-4 py-2 bg-green-600 text-white hover:bg-green-600/90 rounded-md font-medium"
        >
          {currentStep < 4 ? 'Continue' : 'Publish'}
          <ArrowRight class="w-3 h-3 ml-1"/>
        </button>
      </div>
    </div>
  </form>
</div>

<!-- Preview Drawer -->
<PreviewDrawer bind:isOpen={isPreviewOpen} {formData} />

