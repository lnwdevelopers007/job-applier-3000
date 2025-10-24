<script>
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { get } from 'svelte/store';
  import { onMount } from 'svelte';
  import { ChevronLeft } from 'lucide-svelte';
  import BasicInfoForm from '$lib/components/job-post-creation/forms/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/job-post-creation/forms/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/job-post-creation/forms/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/job-post-creation/forms/PostSettingForm.svelte';
  import FormSidebar from '$lib/components/job-post-creation/FormSidebar.svelte';
	import { getUserInfo} from '$lib/utils/auth';
	import JobPreviewDrawer from '$lib/components/job-post-creation/JobPreviewDrawer.svelte';
  import { JobService } from '$lib/services/jobService';

  let jobId = $derived($page.params.id);
  let isPreviewOpen = $state(false);
  let activeSection = $state('basic-info');
  let formData = $state(JobService.createEmptyFormData(getUserInfo()?.userID));

  let validationErrors = $state({});
  let showValidationErrors = $state(false);

  const sections = [
    { id: 'basic-info', name: 'Basic Info' },
    { id: 'description', name: 'Description' },
    { id: 'requirements', name: 'Requirement' },
    { id: 'settings', name: 'Post Settings' }
  ];

  function handlePreview() {
    isPreviewOpen = true;
  }

  function updateValidation(validation) {
    validationErrors = validation.errors;
    showValidationErrors = validation.showErrors;
  }

  async function handlePublish() {
    jobId = get(page).params.id;
    await JobService.handleFormSubmit(formData, updateValidation, true, jobId);
  }

  function scrollToSection(sectionId) {
    activeSection = sectionId;
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }

  function handleScroll() {
    const scrollY = window.scrollY;
    const offset = 100;

    for (const section of sections) {
      const element = document.getElementById(section.id);
      if (element) {
        const elementTop = element.offsetTop - offset;
        const elementBottom = elementTop + element.offsetHeight;

        if (scrollY >= elementTop && scrollY < elementBottom) {
          activeSection = section.id;
          break;
        }
      }
    }
  }

  onMount(async () => {
    const result = await JobService.loadJob(jobId);
    if (result.success && result.data) {
      formData = { ...formData, ...result.data };
    } else {
      console.error('Error loading job:', result.error);
    }
  });

</script>

<svelte:window on:scroll={handleScroll} />

<div>
  <h1 class="text-2xl font-semibold text-gray-900 mb-8">Edit Job Post</h1>

  <div class="flex gap-10">
    <FormSidebar 
      {sections}
      {activeSection}
      onSectionClick={scrollToSection}
      onPreview={handlePreview}
      onPublish={handlePublish}
      actionLabel="Save Changes"
    />

    <div class="flex-1 min-w-0 pl-8 pb-10">
      
      <form class="space-y-8">
        <section id="basic-info" class="mb-10 pb-10 border-b border-gray-200">
          <BasicInfoForm bind:formData {validationErrors} {showValidationErrors} />
        </section>

        <section id="description" class="mb-10 pb-10 border-b border-gray-200">
          <DescriptionForm bind:formData {validationErrors} {showValidationErrors} />
        </section>

        <section id="requirements" class="mb-10 pb-10 border-b border-gray-200">
          <RequirementForm bind:formData {validationErrors} {showValidationErrors} />
        </section>

        <section id="settings" class="mb-10">
          <PostSettingForm bind:formData {validationErrors} {showValidationErrors} />
        </section>
      </form>
    </div>
  </div>
</div>

<JobPreviewDrawer bind:isOpen={isPreviewOpen} {formData} />