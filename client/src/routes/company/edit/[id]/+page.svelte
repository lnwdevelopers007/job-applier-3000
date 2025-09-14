<script>
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { get } from 'svelte/store';
  import { onMount } from 'svelte';
  import { ChevronLeft } from 'lucide-svelte';
  import BasicInfoForm from '$lib/components/forms/job-post-creation/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/forms/job-post-creation/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/forms/job-post-creation/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/forms/job-post-creation/PostSettingForm.svelte';
  import PreviewDrawer from '$lib/components/PreviewDrawer.svelte';
  import FormSidebar from '$lib/components/FormSidebar.svelte';

  let jobId = $derived($page.params.id);
  let isPreviewOpen = $state(false);
  let activeSection = $state('basic-info');
  let loading = true;
  let error = $state(null);

  let formData = $state({
    // Basic Info
    jobTitle: '',
    companyID: '',
    location: '',
    category: '',
    workType: 'full-time',
    workArrangement: 'on-site',
    currency: 'THB',
    minSalary: '',
    maxSalary: '',

    // Description
    jobDescription: '',
    jobSummary: '',

    // Requirements
    requiredSkills: [],
    experienceLevel: '',
    education: '',

    // Post Settings
    applicationDeadline: '',
    numberOfPositions: 1,
    visibility: 'public',
    emailNotifications: true,
    autoReject: false
  });

  const sections = [
    { id: 'basic-info', name: 'Basic Info' },
    { id: 'description', name: 'Description' },
    { id: 'requirements', name: 'Requirement' },
    { id: 'settings', name: 'Post Settings' }
  ];

  function handleBackNavigation() {
    goto('/company/dashboard');
  }

  function handlePreview() {
    isPreviewOpen = true;
  }

  async function handlePublish() {
    jobId = get(page).params.id;
    try {
      const res = await fetch(`/jobs/${jobId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      });

      if (!res.ok) {
        throw new Error(`Failed to update job: ${res.status}`);
      }

      const updatedJob = await res.json();
      console.log('Job updated successfully:', updatedJob);

      goto('/company/dashboard');
    } catch (err) {
      console.error('Error updating job:', err);
    }
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
    try {
      const res = await fetch(`/jobs/query?id=${jobId}`);
      if (!res.ok) throw new Error(`Failed to load job: ${res.status}`);
      const data = await res.json();
      const job = Array.isArray(data) ? data[0] : data;

      const formattedDeadline = job.applicationDeadline
        ? new Date(job.applicationDeadline).toISOString().slice(0, 10)
        : '';

      formData = {
        ...formData,
        ...job,
        requiredSkills: typeof job.requiredSkills === 'string'
          ? job.requiredSkills.split(',').map(s => s.trim()).filter(Boolean)
          : job.requiredSkills || [],
        postingCloseDate: formattedDeadline,
        jobTitle: job.title || '',
        screeningQuestions: job.questions || ''
      };

      console.log('Loaded job:', formData);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  });

</script>

<svelte:window on:scroll={handleScroll} />

<div class="max-w-7xl mx-auto py-8">
  <div class="flex gap-10">
    <FormSidebar 
      {sections}
      {activeSection}
      onSectionClick={scrollToSection}
      onPreview={handlePreview}
      onPublish={handlePublish}
      actionLabel="Update Job Post"
    />

    <div class="flex-1 min-w-0">
      <div class="mb-8">
        <div class="flex items-center gap-4">
          <button onclick={handleBackNavigation} class="p-1 hover:bg-gray-100 rounded">
            <ChevronLeft class="w-5 h-5 text-gray-600" />
          </button>
          <h1 class="text-xl font-medium text-gray-900">Edit Job Post</h1>
          <span class="text-sm text-gray-500">ID: {jobId}</span>
        </div>
      </div>

      <form class="space-y-8">
        <section id="basic-info" class="mb-8 pb-10 border-b border-gray-200">
          <BasicInfoForm bind:formData />
        </section>

        <section id="description" class="mb-8 pb-10 border-b border-gray-200">
          <DescriptionForm bind:formData />
        </section>

        <section id="requirements" class="mb-8 pb-10 border-b border-gray-200">
          <RequirementForm bind:formData />
        </section>

        <section id="settings" class="mb-8">
          <PostSettingForm bind:formData />
        </section>
      </form>
    </div>
  </div>
</div>

<PreviewDrawer bind:isOpen={isPreviewOpen} {formData} />