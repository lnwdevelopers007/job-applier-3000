<script>
  import { fly } from 'svelte/transition';
  import { ArrowLeft, ArrowRight } from 'lucide-svelte';
  import Stepper from '$lib/components/job-post-creation/Stepper.svelte';
  import BasicInfoForm from '$lib/components/job-post-creation/forms/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/job-post-creation/forms/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/job-post-creation/forms/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/job-post-creation/forms/PostSettingForm.svelte';
  import JobPreviewDrawer from '$lib/components/job-post-creation/JobPreviewDrawer.svelte';
  import { getUserInfo } from '$lib/utils/auth';
  import { JobService } from '$lib/services/jobService';

  let currentStep = $state(1); // Current step (1-4)
  let isPreviewOpen = $state(false);
  let stepper;
  
  const steps = [
    { id: 1, name: 'Basic Info' },
    { id: 2, name: 'Description' },
    { id: 3, name: 'Requirement' },
    { id: 4, name: 'Post Settings' }
  ];

  let formData = $state(JobService.createEmptyFormData(getUserInfo()?.userID));

  let validationErrors = $state({});
  let showValidationErrors = $state(false);

  // Validation functions are now imported from jobValidation.ts


  function handleSaveDraft() {
    // Handle save draft
  }

  function handlePreview() {
    isPreviewOpen = true;
  }

  function updateValidation(validation) {
    validationErrors = validation.errors;
    showValidationErrors = validation.showErrors;
  }


  const handleProgress = async (stepIncrement) => {
    // Moving forward - validate current step
    if (stepIncrement === 1) {
      const isValid = await JobService.handleStepValidation(currentStep, formData, updateValidation);
      if (!isValid) return;
    }
    
    // Final submission - validate all steps and submit
    if (currentStep === 4 && stepIncrement === 1) {
      const success = await JobService.handleFormSubmit(formData, updateValidation);
      if (!success) return;
    } else {
      // Clear validation errors when moving to next step
      showValidationErrors = false;
      stepper.handleProgress(stepIncrement);
    }
  };

</script>

<div>
  <h1 class="text-2xl font-semibold text-gray-900 mb-3">Create Job Post</h1>

  <Stepper {steps} bind:currentStep bind:this={stepper} />
  
  <form class="px-8 pb-8">
    {#if currentStep === 1}
      <div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
        <BasicInfoForm bind:formData {validationErrors} {showValidationErrors} />
      </div>
    {:else if currentStep === 2}
      <div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
        <DescriptionForm bind:formData {validationErrors} {showValidationErrors} />
      </div>
    {:else if currentStep === 3}
      <div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
        <RequirementForm bind:formData {validationErrors} {showValidationErrors} />
      </div>
    {:else if currentStep === 4}
      <div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
        <PostSettingForm bind:formData {validationErrors} {showValidationErrors} />
      </div>
    {/if}

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
          class="text-sm flex items-center px-4 py-2 bg-green-600 text-white hover:bg-green-700 rounded-md font-medium"
        >
          {currentStep < 4 ? 'Continue' : 'Publish'}
          <ArrowRight class="w-3 h-3 ml-1"/>
        </button>
      </div>
    </div>
  </form>
</div>

<JobPreviewDrawer bind:isOpen={isPreviewOpen} {formData} />

