<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { MapPin, Bookmark, Share2, Clock, Banknote } from 'lucide-svelte';
  import SafeHTML from '$lib/utils/SafeHTML.svelte';
  import SkillTag from '$lib/components/job/SkillTag.svelte';
  import ApplyButton from '$lib/components/job/ApplyButton.svelte';
  import CompanyCard from '$lib/components/job/CompanyCard.svelte';
  import JobCard from '$lib/components/job/JobCard.svelte';
  import JobDetailSkeleton from '$lib/components/job/JobDetailSkeleton.svelte';
  import ApplyModal from '$lib/components/job/ApplyModal.svelte';
  import FloatingJobHeader from '$lib/components/job/FloatingJobHeader.svelte';
  import { fetchJob, fetchCompany, fetchCompanyNameLogo, DEFAULT_COMPANY_LOGO } from '$lib/utils/fetcher';
  import { isAuthenticated, getUserInfo } from '$lib/utils/auth';
	import { formatDateDMY } from '$lib/utils/datetime';
  import { bookmarkService } from '$lib/services/bookmarkService';
  import { authStore } from '$lib/stores/auth.svelte';

  let { data }: { data: { jobId: string } } = $props();

  let job = $state<any>(null);
  let company = $state<any>(null);
  let similarJobs = $state<any[]>([]);
  let otherCompanyJobs = $state<any[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let isBookmarked = $state(false);
  let appliedJobs = $state(new Set<string>());
  let userInfo: any = null;
  let showFloatingCard = $state(false);
  let applyButtonRef = $state<HTMLElement | undefined>();
  let showApplyModal = $state(false);
  let showAllSkills = $state(false);
  
  // Subscribe to bookmark changes
  $effect(() => {
    const unsubscribe = bookmarkService.subscribe((jobs) => {
      if (job?.id) {
        isBookmarked = jobs.has(job.id);
      }
    });
    
    return unsubscribe;
  });
  
  // Initialize bookmarks when job loads
  $effect(() => {
    if (job?.id) {
      const user = getUserInfo();
      if (user?.userID) {
        bookmarkService.initializeBookmarks(user.userID).then(() => {
          isBookmarked = bookmarkService.isBookmarked(job.id);
        });
      }
    }
  });

  function getRelativeTime(dateString: string): string {
    if (!dateString || dateString === 'Unknown') return 'Unknown';
    
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (date > now) {
      // Future date (for deadline)
      if (diffDays === 0) return 'Today';
      if (diffDays === 1) return 'Tomorrow';
      if (diffDays < 7) return `In ${diffDays} days`;
      if (diffDays < 30) return `In ${Math.floor(diffDays / 7)} weeks`;
      return date.toLocaleDateString();
    } else {
      // Past date (for posted)
      if (diffDays === 0) return 'Today';
      if (diffDays === 1) return 'Yesterday';
      if (diffDays < 7) return `${diffDays} days ago`;
      if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
      if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
      return date.toLocaleDateString();
    }
  }

  function formatSalary(min: number | null, max: number | null, currency: string): string {
    if (!min && !max) return '';
    
    const formatNumber = (num: number) => {
      return new Intl.NumberFormat('en-US').format(Math.round(num));
    };
    
    if (min && max) {
      return `${formatNumber(min)} - ${formatNumber(max)} ${currency}/month`;
    } else if (min) {
      return `From ${formatNumber(min)} ${currency}/month`;
    } else if (max) {
      return `Up to ${formatNumber(max)} ${currency}/month`;
    }
    return '';
  }

  async function fetchAppliedJobs() {
    if (!userInfo?.userID) return;
    
    try {
      const res = await fetch(`/apply/query?applicantID=${userInfo.userID}`, {
        credentials: 'include'
      });
      if (!res.ok) throw new Error('Failed to fetch applied jobs');
      
      const data = await res.json();
      appliedJobs = new Set(data.map((app: any) => app.jobApplication.jobID));
    } catch (err) {
      console.error('Error fetching applied jobs:', err);
      appliedJobs = new Set();
    }
  }

  async function loadJobData() {
    try {
      loading = true;
      error = null;

      // Check authentication and fetch applied jobs first
      if (isAuthenticated()) {
        userInfo = getUserInfo();
        if (userInfo?.userID) {
          await fetchAppliedJobs();
        }
      }

      const jobData = await fetchJob(data.jobId);
      
      job = {
        id: jobData.id,
        title: jobData.title,
        company: '',
        companyId: jobData.companyID,
        location: jobData.location || 'N/A',
        workType: jobData.workType || 'onsite',  // onsite/remote/hybrid
        workArrangement: jobData.workArrangement || 'full-time',  // full-time/part-time/contract
        salary: formatSalary(jobData.minSalary, jobData.maxSalary, jobData.currency || 'THB'),
        posted: jobData.postOpenDate ? getRelativeTime(jobData.postOpenDate) : 'Unknown',
        postedDate: jobData.postOpenDate || null,
        closeDate: jobData.applicationDeadline ? formatDateDMY(jobData.applicationDeadline) : 'Open until filled',
        closeDateRaw: jobData.applicationDeadline || null,
        description: jobData.jobDescription || 'No description provided.',
        skills: jobData.requiredSkills ? jobData.requiredSkills.split(',').map((skill: string) => skill.trim()) : [],
        logo: DEFAULT_COMPANY_LOGO
      };

      if (jobData.companyID) {
        try {
          const [companyName, companyLogo] = await fetchCompanyNameLogo(jobData.companyID);
          job.company = companyName;
          job.logo = companyLogo;

          try {
            const companyData = await fetchCompany(jobData.companyID);
            const infoArray = companyData.userInfo || [];
            const info = Object.fromEntries(infoArray.map((item: any) => [item.Key, item.Value]));

            company = {
              name: companyName,
              logo: companyLogo,
              industry: info.industry || 'Software Development',
              employees: info.size || '1000+ employees',
              description: info.aboutUs || "Company information not available."
            };
          } catch (companyErr) {
            console.warn('Failed to fetch detailed company info:', companyErr);
            // Set basic company info even if detailed fetch fails
            company = {
              name: companyName,
              logo: companyLogo,
              industry: 'Software Development',
              employees: '1000+ employees',
              description: "Company information not available."
            };
          }

          await loadSimilarJobs();
          await loadOtherCompanyJobs(jobData.companyID);
        } catch (companyErr) {
          console.warn('Failed to fetch company name/logo:', companyErr);
          // Use default values if company fetch fails
          job.company = 'Unknown Company';
          job.logo = DEFAULT_COMPANY_LOGO;
        }
      } else {
        job.company = 'Unknown Company';
        job.logo = DEFAULT_COMPANY_LOGO;
      }

    } catch (err) {
      console.error('Error loading job:', err);
      error = 'Failed to load job details. Please try again.';
    } finally {
      loading = false;
    }
  }

  async function loadSimilarJobs() {
    try {
      // Fetch all jobs first
      const res = await fetch('/jobs/query');
      if (res.ok) {
        const data = await res.json();
        
        // Filter out current job and calculate similarity scores
        const otherJobs = data.filter((j: any) => j.id !== job?.id);
        const jobsWithScores = otherJobs.map((jobData: any) => {
          let score = 0;
          
          // Score based on title similarity (keyword matching)
          if (job?.title && jobData.title) {
            const currentTitle = job.title.toLowerCase();
            const otherTitle = jobData.title.toLowerCase();
            
            // Extract keywords from titles (remove common words)
            const commonWords = ['the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by'];
            const currentKeywords = currentTitle.split(/\s+/).filter((word: string) => word.length > 2 && !commonWords.includes(word));
            const otherKeywords = otherTitle.split(/\s+/).filter((word: string) => word.length > 2 && !commonWords.includes(word));
            
            // Calculate keyword overlap
            const matchingKeywords = currentKeywords.filter((keyword: string) => 
              otherKeywords.some((otherKeyword: string) => 
                otherKeyword.includes(keyword) || keyword.includes(otherKeyword)
              )
            );
            score += (matchingKeywords.length / Math.max(currentKeywords.length, 1)) * 3;
          }
          
          // Score based on skills similarity
          if (job?.skills && jobData.requiredSkills) {
            const currentSkills = job.skills.map((s: string) => s.toLowerCase().trim());
            const otherSkills = jobData.requiredSkills.split(',').map((s: string) => s.toLowerCase().trim());
            
            const matchingSkills = currentSkills.filter((skill: string) => 
              otherSkills.some((otherSkill: string) => otherSkill.includes(skill) || skill.includes(otherSkill))
            );
            score += (matchingSkills.length / Math.max(currentSkills.length, 1)) * 2;
          }
          
          // Minor bonus for same work type
          if (job?.workType && jobData.workType === job.workType) {
            score += 0.5;
          }
          
          return { ...jobData, similarityScore: score };
        });
        
        // Sort by similarity score and take top 3
        const filteredJobs = jobsWithScores
          .sort((a: any, b: any) => b.similarityScore - a.similarityScore)
          .slice(0, 3);
        
        similarJobs = await Promise.all(filteredJobs.map(async (jobData: any) => {
          try {
            const [companyName, companyLogo] = await fetchCompanyNameLogo(jobData.companyID || '');
            return {
              id: jobData.id,
              title: jobData.title,
              company: companyName,
              location: jobData.location || 'N/A',
              logo: companyLogo,
              workType: jobData.workType || 'FULL-TIME',
              workArrangement: jobData.workArrangement || 'ON-SITE',
              posted: jobData.postOpenDate ? new Date(jobData.postOpenDate).toLocaleDateString() : 'Unknown'
            };
          } catch (err) {
            // Only log for non-404 company errors
            if (err instanceof Error && !err.message.includes('Company not found')) {
              console.warn(`Failed to fetch company info for similar job ${jobData.id}:`, err);
            }
            return {
              id: jobData.id,
              title: jobData.title,
              company: 'Unknown Company',
              location: jobData.location || 'N/A',
              logo: DEFAULT_COMPANY_LOGO,
              workType: jobData.workType || 'FULL-TIME',
              workArrangement: jobData.workArrangement || 'ON-SITE',
              posted: jobData.postOpenDate ? new Date(jobData.postOpenDate).toLocaleDateString() : 'Unknown'
            };
          }
        }));
      }
    } catch (err) {
      console.error('Error loading similar jobs:', err);
    }
  }

  async function loadOtherCompanyJobs(companyId: string) {
    try {
      const res = await fetch(`/jobs/query?companyID=${companyId}`);
      if (res.ok) {
        const data = await res.json();
        const filteredJobs = data.filter((j: any) => j.id !== job?.id).slice(0, 2);
        
        otherCompanyJobs = await Promise.all(filteredJobs.map(async (jobData: any) => {
          try {
            const [companyName, companyLogo] = await fetchCompanyNameLogo(jobData.companyID || '');
            return {
              id: jobData.id,
              title: jobData.title,
              company: companyName,
              location: jobData.location || 'N/A',
              logo: companyLogo,
              workType: jobData.workType || 'FULL-TIME',
              workArrangement: jobData.workArrangement || 'ON-SITE',
              posted: jobData.postOpenDate ? new Date(jobData.postOpenDate).toLocaleDateString() : 'Unknown'
            };
          } catch (err) {
            // Only log for non-404 company errors
            if (err instanceof Error && !err.message.includes('Company not found')) {
              console.warn(`Failed to fetch company info for other company job ${jobData.id}:`, err);
            }
            return {
              id: jobData.id,
              title: jobData.title,
              company: 'Unknown Company',
              location: jobData.location || 'N/A',
              logo: DEFAULT_COMPANY_LOGO,
              workType: jobData.workType || 'FULL-TIME',
              workArrangement: jobData.workArrangement || 'ON-SITE',
              posted: jobData.postOpenDate ? new Date(jobData.postOpenDate).toLocaleDateString() : 'Unknown'
            };
          }
        }));
      }
    } catch (err) {
      console.error('Error loading other company jobs:', err);
    }
  }


  async function toggleBookmark() {
    if (!job?.id) return;
    
    const user = getUserInfo();
    if (user?.userID) {
      const newState = await bookmarkService.toggleBookmark(job.id, user.userID);
      isBookmarked = newState;
    } else {
      console.warn('User must be logged in to bookmark jobs');
    }
  }

  async function toggleJobBookmark(jobId: string) {
    const user = getUserInfo();
    if (user?.userID) {
      await bookmarkService.toggleBookmark(jobId, user.userID);
    }
  }

  function handleShare() {
    if (navigator.share) {
      navigator.share({
        title: job?.title || 'Job Opportunity',
        text: `Check out this job: ${job?.title} at ${job?.company}`,
        url: window.location.href
      });
    } else {
      navigator.clipboard.writeText(window.location.href);
      alert('Job link copied to clipboard!');
    }
  }

  function navigateToJob(jobId: string) {
    goto(`/app/jobs/${jobId}`);
  }

  function handleApplyClick() {
    showApplyModal = true;
  }


  // React to job ID changes
  $effect(() => {
    if (data.jobId) {
      loadJobData();
    }
  });

  onMount(() => {
    // Handle scroll detection for floating card
    const handleScroll = () => {
      if (applyButtonRef) {
        const rect = applyButtonRef.getBoundingClientRect();
        showFloatingCard = rect.bottom < 0;
      }
    };

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  });
</script>

{#if loading}
  <JobDetailSkeleton />
{:else if error}
  <div class="flex items-center justify-center min-h-[400px]">
    <div class="text-center">
      <p class="text-red-600 mb-4">{error}</p>
      <button 
        onclick={loadJobData}
        class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
      >
        Try Again
      </button>
    </div>
  </div>
{:else if job}
  <div class="mx-auto pb-10">
    <div class="flex gap-6">
      
      <!-- Main Content Container -->
      <div class="flex-1">
        
        <!-- Content Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          <!-- Left Column - Combined Job Header and Details -->
          <div class="lg:col-span-2">
            <div class="rounded-xl p-8">
              
              <!-- Job Header Section -->
              <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-5 pb-5 border-b border-gray-100">
                <div class="flex-1">
                  <div class="flex items-center gap-3 mb-3">
                    <img src={job.logo} alt={job.company} class="w-10 h-10 rounded-lg">
                    <div>
                      <div class="flex items-center gap-2 text-sm text-gray-600">
                        <span class="font-medium text-lg">{job.company}</span>
                      </div>
                    </div>
                  </div>
                  <h1 class="text-2xl font-semibold tracking-tight text-gray-900 mb-3">{job.title}</h1>

                  <div class="flex flex-col gap-2.5 text-sm text-gray-600 mt-4">
                    <div class="flex items-center gap-2">
                      <MapPin class="w-4 h-4 text-gray-400" />
                      <span>{job.location}{job.workType === 'remote' ? ' (Remote)' : job.workType === 'hybrid' ? ' (Hybrid)' : ''}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <Clock class="w-4 h-4 text-gray-400" />
                      <span>{job.workArrangement.replace(/-/g, ' ').toLowerCase().replace(/^\w/, (c: string) => c.toUpperCase())}</span>
                    </div>
                    {#if job.salary}
                      <div class="flex items-center gap-2">
                        <Banknote class="w-4 h-4 text-gray-400" />
                        <span>{job.salary}</span>
                      </div>
                    {/if}
                    <div class="flex items-center gap-2 mt-2">
                      <span>Posted {job.posted}</span>
                    </div>
                  </div>
                </div>

                <div class="flex gap-2">
                  <button 
                    onclick={toggleBookmark}
                    class="px-2 py-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center hover:cursor-pointer {isBookmarked ? 'text-green-700' : 'text-gray-700'}"
                  >
                    <Bookmark class="w-4 h-4 {isBookmarked ? 'fill-current' : ''}" />
                  </button>
                  <button 
                    onclick={handleShare}
                    class="px-2 py-2 text-sm text-gray-700 font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors hover:cursor-pointer"
                  >
                    <Share2 class="w-4 h-4" />
                  </button>
                </div>
              </div>

              <!-- Job Description Section -->
              <h2 class="text-lg font-semibold text-gray-900 mb-4">Job Description</h2>
              <div class="space-y-5 text-sm text-gray-700 leading-7">
                <SafeHTML html={job.description} />
              </div>

              <!-- Skills Section -->
              {#if job.skills && job.skills.length > 0}
                <div class="mt-6 pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-4">Skills</h2>
                  <div class="flex flex-wrap gap-2">
                    {#each (showAllSkills ? job.skills : job.skills) as skill, index (index)}
                      <SkillTag {skill} />
                    {/each}
                  </div>
                </div>
              {/if}

              <!-- About the Company Section -->
              {#if company}
                <div class="mt-6 pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-4">About the Company</h2>
                  <CompanyCard {company} />
                </div>
              {/if}
            </div>
          </div>

          <!-- Right Column -->
          <div class="lg:col-span-1">
            <div bind:this={applyButtonRef} class="bg-gray-100 rounded-xl overflow-hidden p-6 mb-6">
              <h2 class="text-lg font-medium text-gray-900 mb-2">Ready to apply?</h2>
              <p class="text-sm text-gray-600 mb-4">
                Take the next step in your career. Your application will be sent directly to the employer for review.
              </p>
              <div class="text-xs text-gray-500 mb-3">
                Application deadline: {job.closeDate || 'Open until filled'}
              </div>
              <ApplyButton
                isApplied={appliedJobs.has(job.id)}
                closeDateRaw={job.closeDateRaw}
                postedDate={job.postedDate}
                onClick={handleApplyClick}
                size="lg"
                fullWidth={true}
                userRole={authStore.role}
              />
            </div>
            <!-- Similar Jobs -->
            {#if similarJobs.length > 0}
              <div class="bg-white border border-gray-200 rounded-xl overflow-hidden p-4 mb-6">
                <h2 class="text-lg font-medium text-gray-900 mb-3">Similar Jobs</h2>
                <div class="space-y-3">
                  {#each similarJobs as similarJob (similarJob.id)}
                    <JobCard 
                      job={similarJob} 
                      onclick={() => navigateToJob(similarJob.id)}
                      onBookmark={() => toggleJobBookmark(similarJob.id)}
                    />
                  {/each}
                </div>
              </div>
            {/if}

            <!-- Other Jobs from Company -->
            {#if otherCompanyJobs.length > 0}
              <div class="bg-white border border-gray-200 rounded-xl overflow-hidden p-4">
                <h2 class="text-lg font-medium text-gray-900 mb-3">Other jobs from {job.company}</h2>
                <div class="space-y-3">
                  {#each otherCompanyJobs as otherJob (otherJob.id)}
                    <JobCard 
                      job={otherJob} 
                      onclick={() => navigateToJob(otherJob.id)}
                      onBookmark={() => toggleJobBookmark(otherJob.id)}
                    />
                  {/each}
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
{:else}
  <div class="flex items-center justify-center min-h-[400px]">
    <p class="text-gray-600">Job not found.</p>
  </div>
{/if}

<!-- Floating Job Header -->
<FloatingJobHeader 
  show={showFloatingCard}
  {job}
  onApply={handleApplyClick}
  onBookmark={toggleBookmark}
  onShare={handleShare}
  {isBookmarked}
  isApplied={job && appliedJobs.has(job.id)}
  userRole={authStore.role}
/>

<!-- Apply Modal -->
{#if job}
  <ApplyModal 
    bind:isOpen={showApplyModal}
    {job}
  />
{/if}