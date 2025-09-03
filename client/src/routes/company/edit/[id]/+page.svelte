<script>
  import { goto } from '$app/navigation';
  import { ChevronLeft } from 'lucide-svelte';
  import BasicInfoForm from '$lib/components/forms/job-post-creation/BasicInfoForm.svelte';
  import DescriptionForm from '$lib/components/forms/job-post-creation/DescriptionForm.svelte';
  import RequirementForm from '$lib/components/forms/job-post-creation/RequirementForm.svelte';
  import PostSettingForm from '$lib/components/forms/job-post-creation/PostSettingForm.svelte';
  import PreviewDrawer from '$lib/components/PreviewDrawer.svelte';
  import FormSidebar from '$lib/components/FormSidebar.svelte';
  
  let isPreviewOpen = $state(false);
  let activeSection = $state('basic-info');

  let jobId = 'demo-123';

  let formData = $state({
    // Basic Info
    jobTitle: 'Senior Software Engineer',
    location: 'Bangkok, Thailand',
    category: 'Technology',
    workType: 'full-time',
    workArrangement: 'hybrid',
    currency: 'THB',
    minSalary: '80000',
    maxSalary: '120000',
    
    // Description
    jobDescription: 'We are looking for a Senior Software Engineer...',
    
    // Requirements
    requiredSkills: ['JavaScript', 'React', 'Node.js'],
    experienceLevel: 'Mid-Level (4-6 years)',
    education: 'Bachelor\'s Degree',
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

  function handlePublish() {
    console.log('Publishing job:', jobId);
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