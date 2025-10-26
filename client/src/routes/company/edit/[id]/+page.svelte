<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import BasicInfoForm from '$lib/components/job-post-creation/forms/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/job-post-creation/forms/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/job-post-creation/forms/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/job-post-creation/forms/PostSettingForm.svelte';
  import FormSidebar from '$lib/components/job-post-creation/FormSidebar.svelte';
	import { getUserInfo} from '$lib/utils/auth';
	import JobPreviewDrawer from '$lib/components/job-post-creation/JobPreviewDrawer.svelte';
  import { JobService } from '$lib/services/jobService';
  import toast from 'svelte-french-toast';

  let jobId = $derived($page.params.id);
  let isPreviewOpen = $state(false);
  let activeSection = $state('basic-info');
  let formData = $state(JobService.createEmptyFormData(getUserInfo()?.userID));

  let validationErrors = $state({});
  let showValidationErrors = $state(false);

  const sections = [
    { 
      id: 'basic-info', 
      name: 'Basic Info',
      title: 'Basic Information',
      description: 'Set up the fundamental details of your job posting including title, location, and compensation.'
    },
    { 
      id: 'description', 
      name: 'Description',
      title: 'Job Description',
      description: 'Provide a detailed overview of the role, responsibilities, and what makes this opportunity unique.'
    },
    { 
      id: 'requirements', 
      name: 'Requirement',
      title: 'Requirements',
      description: 'Specify the qualifications, skills, and experience needed for this position.'
    },
    { 
      id: 'settings', 
      name: 'Post Settings',
      title: 'Post Settings',
      description: 'Configure application requirements, posting timeline, and notification preferences.'
    }
  ];

  function handlePreview() {
    isPreviewOpen = true;
  }

  function updateValidation(validation: { errors: Record<string, string>; showErrors: boolean }) {
    validationErrors = validation.errors;
    showValidationErrors = validation.showErrors;
    
    // If there are errors, switch to the first tab that has errors
    if (validation.showErrors && Object.keys(validation.errors).length > 0) {
      switchToFirstErrorTab(validation.errors);
    }
  }
  
  function switchToFirstErrorTab(errors: Record<string, string>) {
    // Define which fields belong to which section
    const fieldToSection: Record<string, string> = {
      // Basic Info fields
      jobTitle: 'basic-info',
      location: 'basic-info',
      workType: 'basic-info',
      workArrangement: 'basic-info',
      minSalary: 'basic-info',
      maxSalary: 'basic-info',
      currency: 'basic-info',
      
      // Description fields
      jobDescription: 'description',
      
      // Requirements fields
      yearsExperience: 'requirements',
      educationLevel: 'requirements',
      requiredSkills: 'requirements',
      
      // Settings fields
      postingOpenDate: 'settings',
      postingCloseDate: 'settings',
      applicationRequirements: 'settings'
    };
    
    // Find the first section with errors
    for (const errorField of Object.keys(errors)) {
      const section = fieldToSection[errorField];
      if (section && section !== activeSection) {
        activeSection = section;
        break;
      }
    }
  }

  function handleSectionClick(sectionId: string) {
    activeSection = sectionId;
  }

  async function handleSave() {
    // Get the step number based on active section
    const stepMap: Record<string, number> = {
      'basic-info': 1,
      'description': 2,
      'requirements': 3,
      'settings': 4
    };
    
    const currentStep = stepMap[activeSection];
    
    // Validate current step
    const isValid = await JobService.handleStepValidation(currentStep, formData, updateValidation);
    
    if (isValid) {
      // Save the job
      const result = await JobService.handleFormSubmit(formData, updateValidation, true, jobId);
      if (result) {
        toast.success('Job saved successfully!');
      }
    }
  }

  function getActiveSection() {
    return sections.find(s => s.id === activeSection);
  }

  onMount(async () => {
    if (jobId) {
      const result = await JobService.loadJob(jobId);
      if (result.success && result.data) {
        formData = { ...formData, ...result.data };
      } else {
        console.error('Error loading job:', result.error);
      }
    }
  });

</script>


<div>
  <h1 class="text-2xl font-semibold text-gray-900 mb-8">Edit Job Post</h1>

  <div class="flex gap-10">
    <FormSidebar 
      {sections}
      {activeSection}
      onSectionClick={handleSectionClick}
    />

    <div class="flex-1 min-w-0 pl-8 pb-10">
      <!-- Section Header with Action Buttons -->
      <div class="mb-6 pb-4 border-b border-gray-200">
        <div class="flex justify-between items-start">
          <div>
            <h2 class="text-lg font-medium text-gray-900">{getActiveSection()?.title || ''}</h2>
            <p class="text-sm text-gray-500 mt-1">{getActiveSection()?.description || ''}</p>
          </div>
          <div class="flex gap-3">
            <button
              onclick={handlePreview}
              type="button"
              class="px-4 py-2 text-sm text-gray-600 bg-white border border-gray-200 rounded-md hover:bg-gray-50 transition-colors font-medium"
            >
              Preview
            </button>
            <button
              onclick={handleSave}
              type="button"
              class="px-4 py-2 text-sm bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors font-medium"
            >
              Save Changes
            </button>
          </div>
        </div>
      </div>
      
      <form>
        {#if activeSection === 'basic-info'}
          <BasicInfoForm bind:formData {validationErrors} {showValidationErrors} />
        {:else if activeSection === 'description'}
          <DescriptionForm bind:formData {validationErrors} {showValidationErrors} />
        {:else if activeSection === 'requirements'}
          <RequirementForm bind:formData {validationErrors} {showValidationErrors} />
        {:else if activeSection === 'settings'}
          <PostSettingForm bind:formData {validationErrors} {showValidationErrors} />
        {/if}
      </form>
    </div>
  </div>
</div>

<JobPreviewDrawer bind:isOpen={isPreviewOpen} {formData} />