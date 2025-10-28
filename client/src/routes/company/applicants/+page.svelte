<script lang="ts">
  import { onMount } from 'svelte';
  import { Search } from 'lucide-svelte';
  import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
  import { authStore } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import SeekerProfileCard from '$lib/components/profile/SeekerProfileCard.svelte';
  import Badge from '$lib/components/job/Badge.svelte';
  
  type BadgeVariant = 'primary' | 'secondary' | 'warning' | 'danger' | 'success' | 'info' | 'purple';

  let candidates: any[] = [];
  let selectedCandidate: any = null;
  let currentStatusFilter: string = 'all';
  let searchQuery: string = '';
  let selectedJobFilter: string = 'all';
  let companyJobs: any[] = [];

  let company: any = null;
  let pendingCount: number = 0;
  let currentPage: number = 1;
  let itemsPerPage: number = 5;
  let totalPages: number = 1;
  let isUpdatingStatus: boolean = false;

  $: filteredCandidates = candidates.filter(c => {
    // Status filter
    if (currentStatusFilter !== 'all' && c.status.toLowerCase() !== currentStatusFilter.toLowerCase()) {
      return false;
    }
    
    // Job filter
    if (selectedJobFilter !== 'all' && c.applied !== selectedJobFilter) {
      return false;
    }
    
    // Search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      return (
        c.name?.toLowerCase().includes(query) ||
        c.email?.toLowerCase().includes(query) ||
        c.role?.toLowerCase().includes(query) ||
        c.skills?.some((skill: string) => skill.toLowerCase().includes(query))
      );
    }
    
    return true;
  });

  $: paginatedCandidates = filteredCandidates.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  $: totalPages = Math.ceil(filteredCandidates.length / itemsPerPage);


  function getStatusVariant(status: string): BadgeVariant {
    switch (status) {
      case 'Accepted': return 'success';  // green
      case 'Rejected': return 'danger';   // red
      case 'Pending': return 'warning';   // yellow
      default: return 'secondary';
    }
  }

  function normalizeUser(user: any) {
    const infoArray = user.userInfo || [];
    const info = Object.fromEntries(infoArray.map((i: any) => [i.Key, i.Value]));
    
    // Parse skills from comma-separated string
    const skillsString = info.skills || '';
    const skillsArray = skillsString ? skillsString.split(',').map((skill: string) => skill.trim()).filter((skill: string) => skill.length > 0) : [];
    
    return {
      id: user.id,
      userID: user.userID,
      email: user.email || '',
      avatar: user.avatarURL || 'https://cdn-icons-png.flaticon.com/512/149/149071.png',
      name: info.fullName || user.name || 'Unknown User',
      fullName: info.fullName || user.name || 'Unknown User',
      role: info.desiredRole || 'Unknown Role',
      phone: info.phone || '-',
      linkedIn: info.linkedIn || '-',
      location: info.location || '-',
      github: info.github || '',
      portfolio: info.portfolio || '',
      aboutMe: info.aboutMe || '',
      dateOfBirth: info.dateOfBirth || '',
      education: user.education || [],
      skills: skillsArray,
      documents: user.documents || []
    };
  }

  async function fetchCompanyJobs() {
    if (!company?.userID) return [];
    
    try {
      const jobsRes = await fetch(`/jobs/query?companyID=${company.userID}`);
      if (!jobsRes.ok) return [];
      
      const jobsData = await jobsRes.json();
      return jobsData.map((job: any) => ({
        id: job.id || job._id,
        title: job.title
      }));
    } catch (err) {
      console.error('Error fetching company jobs:', err);
      return [];
    }
  }

  async function fetchCandidates(status: string = '') {
    if (!isAuthenticated()) {
      goto('/login');
      return [];
    }

    if (!company?.userID) {
      console.error('Company not found');
      goto('/login');
      return [];
    }

    try {
      // Get all jobs owned by this company
      const jobsRes = await fetch(`/jobs/query?companyID=${company.userID}`, {
      });
      if (!jobsRes.ok) throw new Error('Failed to fetch company jobs');

      const jobsData = await jobsRes.json();
      if (!Array.isArray(jobsData) || jobsData.length === 0) return [];

      const allApplications: any[] = [];

      // For each job, fetch applications
      for (const job of jobsData) {
        const jobID = job.id || job._id;
        if (!jobID) continue;

        const applyRes = await fetch(`/apply?jobID=${jobID}`, {
        });
        if (!applyRes.ok) continue;

        const applyData = await applyRes.json();

        // Filter by status
        const filtered = status
          ? applyData.filter(
              (a: any) => a.jobApplication.status.toUpperCase() === status.toUpperCase()
            )
          : applyData;

        for (const app of filtered) {
          allApplications.push({
            ...app,
            jobTitle: job.title || 'Unknown Job',
            jobID
          });
        }
      }

      // Fetch applicant info
      const candidatePromises = allApplications.map(async (app: any) => {
        const applicantID = app.jobApplication.applicantID;

        const userRes = await fetch(`/users/query?id=${applicantID}`, {
        });
        const userData = await userRes.json();
        console.log('Fetched user data for applicant ID', applicantID, ':', userData);
        const user = normalizeUser(userData && userData[0] ? userData[0] : {});
        console.log('Fetched user data for applicant:', user);

        const created = new Date(app.jobApplication.createdAt);
        const daysAgo = Math.floor((Date.now() - created.getTime()) / (1000 * 60 * 60 * 24));

        const rawStatus = app.jobApplication.status || 'PENDING';
        const displayStatus =
          rawStatus.toUpperCase() === 'PENDING'
            ? 'Pending'
            : rawStatus.charAt(0) + rawStatus.slice(1).toLowerCase();

        return {
          id: app.jobApplication.id,
          applicantID: app.jobApplication.applicantID,
          jobID: app.jobApplication.jobID,
          createdAt: app.jobApplication.createdAt,
          name: user.fullName,
          fullName: user.fullName,
          role: user.role,
          applied: app.jobTitle,
          status: displayStatus,
          avatar: user.avatar,
          time: daysAgo === 0 ? 'Today' : `${daysAgo} days ago`,
          email: user.email,
          phone: user.phone,
          address: user.location,
          linkedin: user.linkedIn,
          github: user.github,
          portfolio: user.portfolio,
          aboutMe: user.aboutMe,
          dateOfBirth: user.dateOfBirth,
          education: user.education,
          skills: user.skills,
          documents: user.documents
        };
      });

      const results = await Promise.all(candidatePromises);
      pendingCount = results.filter(c => c.status === 'Pending').length;
      return results;

    } catch (err) {
      console.error('Error fetching candidates:', err);
      return [];
    }
  }

  async function updateStatus(candidateID: string, newStatus: string) {
    if (isUpdatingStatus) return;
    isUpdatingStatus = true;
    try {
      const candidate = candidates.find(c => c.id === candidateID);
      if (!candidate) throw new Error('Candidate not found');
      const res = await fetch(`/apply/${candidateID}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: candidate.id,
          applicantID: candidate.applicantID,
          jobID: candidate.jobID,
          status: newStatus.toUpperCase(),
          createdAt: candidate.createdAt
        })
      });

      if (!res.ok) throw new Error('Failed to update status');

      // Fetch all candidates (no filter) to keep the full list, then apply client-side filtering
      candidates = await fetchCandidates();
      if (selectedCandidate)
        selectedCandidate = candidates.find((c) => c.id === candidateID) || null;
    } catch (err) {
      console.error('Error updating candidate status:', err);
    } finally {
      isUpdatingStatus = false;
    }
  }

  function handleStatusFilter(status: string) {
    currentStatusFilter = status.toLowerCase();
    currentPage = 1; // Reset to first page when filtering
    // Select first candidate from filtered results
    if (filteredCandidates.length > 0) {
      selectedCandidate = filteredCandidates[0];
    }
  }

  function changePage(page: number) {
    if (page < 1 || page > totalPages) return;
    currentPage = page;
  }

  onMount(async () => {
    // Wait a bit for auth store to be initialized from layout
    await new Promise(resolve => setTimeout(resolve, 100));
    
    company = getUserInfo();
    console.log('Company info after getUserInfo():', company);
    console.log('Auth store state:', { isAuthenticated: isAuthenticated(), authStore });
    
    if (company) {
      console.log('Fetching candidates for company:', company.userID);
      // Fetch jobs and candidates in parallel
      const [jobsData, candidatesData] = await Promise.all([
        fetchCompanyJobs(),
        fetchCandidates()
      ]);
      
      companyJobs = jobsData;
      candidates = candidatesData;
      
      console.log('Fetched jobs:', companyJobs);
      console.log('Fetched candidates:', candidates);
      if (candidates.length > 0) selectedCandidate = candidates[0];
    } else {
      console.log('No company info found - user may need to log in');
      // Try to redirect to login if not authenticated
      if (!isAuthenticated()) {
        goto('/login');
      }
    }
  });
</script>

<div class="h-[calc(100vh-110px)] flex flex-col">
  <h1 class="text-2xl font-semibold text-gray-900 mb-1">
    All Applicants
  </h1>
  <p class="mb-6 text-base text-gray-600">
    Manage candidates across all job postings
  </p>

  <div class="mb-6">
    <div class="flex gap-3">
      <div class="flex-1 relative">
        <div class="absolute left-3 top-1/2 -translate-y-1/2">
          <Search class="w-5 h-5 text-gray-400" />
        </div>
        <input
          type="text"
          placeholder="Search candidates..."
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-full text-sm placeholder:text-gray-500 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400 transition-all"
          bind:value={searchQuery}
        />
      </div>

      <div class="relative">
        <select 
          class="appearance-none px-4 pr-10 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700 font-medium hover:bg-gray-100 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400 transition-all cursor-pointer"
          bind:value={selectedJobFilter}
        >
          <option value="all">All Jobs</option>
          {#each companyJobs as job}
            <option value={job.title}>{job.title}</option>
          {/each}
        </select>
        <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-3 gap-4 flex-1 overflow-hidden">
    <div class="col-span-1 rounded-xl border border-gray-200 overflow-hidden flex flex-col">
      <div class="bg-slate-50 p-4 border-b border-gray-200">
        <div class="flex items-center justify-between mb-2">
          <h2 class="text-lg font-semibold text-gray-800">Applicants</h2>
          <span class="text-sm text-gray-500">{filteredCandidates.length} total</span>
        </div>
        <div class="flex gap-1">
          <button 
            class="px-3 py-1.5 rounded-full text-xs transition-colors {currentStatusFilter === 'all' ? 'bg-green-600 text-white' : 'bg-white border border-gray-200 hover:bg-gray-50'}"
            onclick={() => handleStatusFilter('all')}
          >
            All
          </button>
          <button 
            class="px-3 py-1.5 rounded-full text-xs transition-colors {currentStatusFilter === 'pending' ? 'bg-green-600 text-white' : 'bg-white border border-gray-200 hover:bg-gray-50'}"
            onclick={() => handleStatusFilter('pending')}
          >
            Pending
          </button>
          <button 
            class="px-3 py-1.5 rounded-full text-xs transition-colors {currentStatusFilter === 'accepted' ? 'bg-green-600 text-white' : 'bg-white border border-gray-200 hover:bg-gray-50'}"
            onclick={() => handleStatusFilter('accepted')}
          >
            Accepted
          </button>
          <button 
            class="px-3 py-1.5 rounded-full text-xs transition-colors {currentStatusFilter === 'rejected' ? 'bg-green-600 text-white' : 'bg-white border border-gray-200 hover:bg-gray-50'}"
            onclick={() => handleStatusFilter('rejected')}
          >
            Rejected
          </button>
        </div>
      </div>
      
      <div class="bg-white p-4 space-y-3 flex-1 overflow-y-auto">
        {#each paginatedCandidates as candidate, i (i)}
          <button
            class="relative flex items-start gap-4 px-4 py-5 border rounded-xl bg-white hover:bg-gray-50 transition-all duration-200 cursor-pointer w-full
              {selectedCandidate?.id === candidate.id 
                ? 'bg-gradient-to-r from-green-50 to-emerald-50 border-green-600 shadow-green-200/50' 
                : selectedCandidate?.name === candidate.name 
                  ? 'border-green-600' 
                  : 'border-gray-200'}"
            onclick={() => (selectedCandidate = candidate)}
          >
            <!-- Status Badge - Top Right -->
            <div class="absolute top-3 right-3">
              <Badge 
                variant={getStatusVariant(candidate.status) as any} 
                size="sm" 
                text={candidate.status} 
              />
            </div>
            
            <div class="relative">
              <img src={candidate.avatar} alt={candidate.name} class="w-12 h-12 rounded-full object-cover ring-2 ring-white shadow-sm" />
            </div>
            <div class="flex-1 flex flex-col items-start gap-1 pr-16">
              <h3 class="font-semibold text-gray-900 text-md tracking-tight">{candidate.name}</h3>
              <p class="text-start text-sm font-medium text-gray-700">{candidate.role}</p>
              <p class="text-start text-xs text-gray-500 font-light">Applied for: <span class="font-medium text-gray-600">{candidate.applied}</span></p>
              <div class="flex items-center text-xs text-gray-500 gap-3 mt-1">
                <span class="flex items-center gap-1">
                  <span class="font-medium">{candidate.time}</span>
                </span>
              </div>
            </div>
          </button>
        {/each}
        {#if totalPages > 1}
          <div class="flex items-center justify-center gap-3 py-3">
            <button
              aria-label="Previous page"
              class="flex items-center gap-2 px-2 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              onclick={() => changePage(Math.max(1, currentPage - 1))}
              disabled={currentPage === 1}
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>

            <div class="flex items-center gap-2">
              <span class="text-sm text-gray-600">Page</span>
              <span class="text-sm text-gray-600 font-medium">
                {currentPage}
              </span>
              <span class="text-sm text-gray-600">of</span>
              <span class="text-sm text-gray-600 font-medium">
                {totalPages}
              </span>
            </div>

            <button
              aria-label="Next page"
              class="flex items-center gap-2 px-2 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              onclick={() => changePage(Math.min(totalPages, currentPage + 1))}
              disabled={currentPage === totalPages}
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        {/if}
      </div>
    </div>

    <div class="col-span-2 bg-white rounded-xl border border-gray-200 overflow-y-auto">
      {#if selectedCandidate}
        <SeekerProfileCard 
          userData={{
            fullName: selectedCandidate.fullName,
            name: selectedCandidate.name,
            email: selectedCandidate.email,
            phone: selectedCandidate.phone,
            location: selectedCandidate.address,
            linkedin: selectedCandidate.linkedin,
            github: selectedCandidate.github,
            portfolio: selectedCandidate.portfolio,
            avatar: selectedCandidate.avatar,
            desiredRole: selectedCandidate.role,
            aboutMe: selectedCandidate.aboutMe,
            skills: selectedCandidate.skills,
            documents: selectedCandidate.documents,
            dateOfBirth: selectedCandidate.dateOfBirth
          }} 
          showApplicationInfo={true}
          isPreviewMode={false}
          appliedJobTitle={selectedCandidate.applied}
          onAccept={(id: string) => updateStatus(id, 'ACCEPTED')}
          onReject={(id: string) => updateStatus(id, 'REJECTED')}
          onNotes={() => console.log('Notes clicked')}
          isUpdatingStatus={isUpdatingStatus}
          candidateStatus={selectedCandidate.status}
          candidateId={selectedCandidate.id}
        />
      {:else}
        <div class="flex items-center justify-center h-full text-gray-500">
          <p>Select a candidate to view their profile</p>
        </div>
      {/if}
    </div>
  </div>
</div>


