<script lang="ts">
  import { X, MapPin, Users, ExternalLink } from 'lucide-svelte';
  import { fade, fly } from 'svelte/transition';
  
  let { 
    isOpen = $bindable(false),
    formData = {}
  } = $props();
  
  let activeTab = $state('detail'); // 'detail' or 'search'
  
  function closeDrawer() {
    isOpen = false;
  }
  
  function formatSalary() {
    if (!formData.salaryMin || !formData.salaryMax) return 'Not specified';
    return `${formData.salaryMin.toLocaleString()}-${formData.salaryMax.toLocaleString()}`;
  }
  
  function formatDate(dateString: string) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  }
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div 
    class="fixed inset-0 bg-black/50 z-60"
    onclick={closeDrawer}
    transition:fade={{ duration: 300 }}
  ></div>
  
  <!-- Drawer -->
  <div 
    class="fixed right-0 top-0 h-full w-[820px] bg-white shadow-xl z-70 overflow-y-auto"
    transition:fly={{ x: 820, duration: 300, opacity: 1 }}
  >
    <div class="sticky top-0 bg-white px-6 py-4 flex items-center justify-between">
      <div>
        <h2 class="text-lg font-semibold text-gray-900">Preview</h2>
        <p class="text-sm text-gray-500 mt-1">This is a preview of what your job post will look like to appliers</p>
      </div>
      <button
        onclick={closeDrawer}
        class="p-1 hover:bg-gray-100 rounded-md transition-colors"
      >
        <X class="w-5 h-5 text-gray-500" />
      </button>
    </div>
    
    <div class="p-6">
      <!-- Tab Navigation -->
      <div class="flex gap-4 mb-3 pl-6">
        <div class="border-b border-gray-200">
          <button 
            onclick={() => activeTab = 'detail'}
            class="pb-3 px-1 text-sm font-medium transition-colors {activeTab === 'detail' ? 'border-b-2 border-green-600 text-gray-900' : 'text-gray-500 hover:text-gray-700'}"
          >
            Job detail
          </button>
          <button 
            onclick={() => activeTab = 'search'}
            class="pb-3 px-1 text-sm font-medium transition-colors {activeTab === 'search' ? 'border-b-2 border-green-600 text-gray-900' : 'text-gray-500 hover:text-gray-700'}"
          >
            Search result
          </button>
        </div>
      </div>
      
      {#if activeTab === 'detail'}
        <!-- Job Detail View -->
        <div class="p-6">
          <div class="p-10 border border-gray-200 rounded-xl">
          <div class="mb-6">
            <div class="flex items-start gap-4">
              <div class="w-12 h-12 rounded-lg bg-white border border-gray-200 flex items-center justify-center">
                {#if formData.companyLogo}
                  <img src={formData.companyLogo} alt={formData.companyName} class="w-10 h-10 object-contain" />
                {:else}
                  <span class="text-lg font-bold text-gray-400">
                    {formData.companyName ? formData.companyName.charAt(0) : 'C'}
                  </span>
                {/if}
              </div>
              <div class="flex-1">
                <h3 class="text-xs font-medium text-gray-600">{formData.companyName || 'Company Name'}</h3>
                <h1 class="text-xl font-semibold text-gray-900 mt-1 flex items-center gap-2">
                  {formData.jobTitle || 'Job Title'}
                  <ExternalLink class="w-4 h-4 text-gray-400" />
                </h1>
              </div>
            </div>
          </div>
          
          <div class="flex items-center gap-4 text-sm text-gray-600 mb-6">
            <span class="text-sm">{formatDate(formData.postingOpenDate) || '1 day ago'}</span>
            <span class="text-gray-400">·</span>
            <span>Be the first to apply</span>
            <span class="text-gray-400">·</span>
            <span>Salary: {formatSalary()}</span>
          </div>
          
          <div class="flex flex-wrap gap-2 mb-6">
            <span class="px-3 py-1 bg-gray-100 text-gray-700 text-sm rounded-full">
              {formData.employmentType || 'Full-time'}
            </span>
            <span class="px-3 py-1 bg-gray-100 text-gray-700 text-sm rounded-full">
              {formData.workLocation || 'On-site'}
            </span>
            <span class="px-3 py-1 bg-gray-100 text-gray-700 text-sm rounded-full">
              {formData.location || 'Bangkok, Thailand'}
            </span>
          </div>
          
          <div class="mb-8">
            <button class="flex-1 px-4 py-2 bg-primary text-white text-sm rounded-md hover:bg-primary-600 transition-colors font-medium mr-1">
              Apply
            </button>
            <button class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 text-sm rounded-md hover:bg-gray-50 transition-colors font-medium">
              Bookmark
            </button>
          </div>
          
          <div class="mb-8">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Job Description</h3>
            {#if formData.jobDescription}
              <div class="prose prose-sm text-gray-600 max-w-full overflow-hidden line-clamp-5">
                {@html formData.jobDescription}
              </div>
              <button class="text-green-600 hover:text-green-700 text-sm font-medium mt-2">
                Show more
              </button>
            {:else}
              <p class="text-gray-500 italic">No job description provided</p>
            {/if}
          </div>
          
          {#if formData.requiredSkills?.length > 0 || formData.yearsOfExperience || formData.educationLevel}
            <div class="mb-8">
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Requirements</h3>
              
              {#if formData.yearsOfExperience}
                <div class="mb-3">
                  <h4 class="text-sm font-medium text-gray-700 mb-1">Experience</h4>
                  <p class="text-sm text-gray-600">{formData.yearsOfExperience}</p>
                </div>
              {/if}
              
              {#if formData.educationLevel}
                <div class="mb-3">
                  <h4 class="text-sm font-medium text-gray-700 mb-1">Education</h4>
                  <p class="text-sm text-gray-600">{formData.educationLevel}</p>
                </div>
              {/if}
              
              {#if formData.requiredSkills?.length > 0}
                <div class="mb-3">
                  <h4 class="text-sm font-medium text-gray-700 mb-2">Skills & Technologies</h4>
                  <div class="flex flex-wrap gap-2">
                    {#each formData.requiredSkills as skill}
                      <span class="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded">
                        {skill}
                      </span>
                    {/each}
                  </div>
                </div>
              {/if}
            </div>
          {/if}
          
          <div class="mb-8">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">About the Company</h3>
            
            <div class="flex items-start gap-4 bg-slate-50 p-6 rounded-xl">
              <div class="w-12 h-12 rounded-lg bg-white border border-gray-200 flex items-center justify-center">
                {#if formData.companyLogo}
                  <img src={formData.companyLogo} alt={formData.companyName} class="w-10 h-10 object-contain" />
                {:else}
                  <span class="text-lg font-bold text-gray-400">
                    {formData.companyName ? formData.companyName.charAt(0) : 'C'}
                  </span>
                {/if}
              </div>
              <div class="flex-1">
                <h4 class="font-semibold text-gray-900">{formData.companyName || 'Company Name'}</h4>
                <p class="text-sm text-gray-600">{formData.industry || 'Software Development'}</p>
                
                <div class="flex items-center gap-4 mt-3 text-sm text-gray-600">
                  <div class="flex items-center gap-1">
                    <Users class="w-4 h-4" />
                    <span>{formData.companySize || '10,000+'} employees</span>
                  </div>
                </div>
                
                {#if formData.companyDescription}
                  <p class="text-sm text-gray-600 mt-3 line-clamp-3">
                    {formData.companyDescription}
                  </p>
                {:else}
                  <p class="text-sm text-gray-500 mt-3 line-clamp-3">
                    A problem isn't truly solved until it's solved for all. Googlers build products 
                    that help create opportunities for everyone, whether down the street or 
                    across the globe. Bring your insight, imagination and a healthy disregard for...
                  </p>
                {/if}
                
                <button class="text-green-600 hover:text-green-700 text-sm font-medium mt-2">
                  Show more
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      {:else}
        <!-- Search Result View -->
        <div class="space-y-4 p-6">
          <div class="p-4 border border-gray-200 rounded-lg animate-pulse">
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 bg-gray-200 rounded"></div>
              <div class="flex-1">
                <div class="h-5 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-1/2 mb-3"></div>
                <div class="flex gap-2">
                  <div class="h-6 bg-gray-200 rounded-full w-20"></div>
                  <div class="h-6 bg-gray-200 rounded-full w-16"></div>
                </div>
                <div class="h-3 bg-gray-200 rounded w-full mt-3"></div>
                <div class="h-3 bg-gray-200 rounded w-5/6 mt-2"></div>
              </div>
            </div>
          </div>

          <!-- Active Job Card -->
          <div class="p-4 border border-primary-300 rounded-lg bg-white">
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 rounded bg-white border border-gray-200 flex items-center justify-center flex-shrink-0">
                {#if formData.companyLogo}
                  <img src={formData.companyLogo} alt={formData.companyName} class="w-8 h-8 object-contain" />
                {:else}
                  <span class="text-sm font-bold text-gray-400">
                    {formData.companyName ? formData.companyName.charAt(0) : 'G'}
                  </span>
                {/if}
              </div>
              <div class="flex-1">
                <h3 class="font-semibold text-gray-900">{formData.jobTitle || 'Software Engineer, AI Acceleration'}</h3>
                <p class="text-sm text-gray-600">{formData.companyName || 'Google'}</p>
                <div class="flex items-center gap-2 mt-2">
                  <MapPin class="w-3 h-3 text-gray-400" />
                  <span class="text-sm text-gray-600">{formData.location || 'Bangkok, Thailand'}</span>
                </div>
              </div>
            </div>
            <div class="flex gap-2 mt-3">
              <span class="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full uppercase font-medium">
                {formData.employmentType || 'Full-time'}
              </span>
              <span class="px-2 py-1 bg-orange-100 text-orange-700 text-xs rounded-full uppercase font-medium">
                {formData.workLocation || 'On-site'}
              </span>
            </div>
            <p class="text-xs text-gray-500 mt-3">{formatDate(formData.postingOpenDate) || '1 day ago'}</p>
          </div>

          <!-- Skeleton Card 2 -->
          <div class="p-4 border border-gray-200 rounded-lg animate-pulse">
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 bg-gray-200 rounded"></div>
              <div class="flex-1">
                <div class="h-5 bg-gray-200 rounded w-2/3 mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-1/3 mb-3"></div>
                <div class="flex gap-2">
                  <div class="h-6 bg-gray-200 rounded-full w-24"></div>
                </div>
                <div class="h-3 bg-gray-200 rounded w-full mt-3"></div>
              </div>
            </div>
          </div>

          <!-- Skeleton Card 3 -->
          <div class="p-4 border border-gray-200 rounded-lg animate-pulse">
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 bg-gray-200 rounded"></div>
              <div class="flex-1">
                <div class="h-5 bg-gray-200 rounded w-4/5 mb-2"></div>
                <div class="h-4 bg-gray-200 rounded w-2/5 mb-3"></div>
                <div class="flex gap-2">
                  <div class="h-6 bg-gray-200 rounded-full w-20"></div>
                  <div class="h-6 bg-gray-200 rounded-full w-20"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}