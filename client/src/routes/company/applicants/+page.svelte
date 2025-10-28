<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    GraduationCap,
    Briefcase,
    Clock,
    FileText,
    Github,
    Link
  } from 'lucide-svelte';
  import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
  import { authStore } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import SeekerProfileCard from '$lib/components/profile/SeekerProfileCard.svelte';

  let candidates: any[] = [];
  let selectedCandidate: any = null;
  let currentStatusFilter: string = '';

  let company: any = null;
  let pendingCount: number = 0;
  let currentPage: number = 1;
  let itemsPerPage: number = 5;
  let totalPages: number = 1;
  let isUpdatingStatus: boolean = false;

  $: paginatedCandidates = candidates.slice(
  (currentPage - 1) * itemsPerPage,
  currentPage * itemsPerPage
  );

  function getDocIcon(doc: string) {
    if (doc.toLowerCase().endsWith('.pdf')) return FileText;
    if (doc.toLowerCase().includes('github')) return Github;
    return Link;
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
      totalPages = Math.ceil(results.length / itemsPerPage);
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

      candidates = await fetchCandidates(currentStatusFilter);
      if (selectedCandidate)
        selectedCandidate = candidates.find((c) => c.id === candidateID) || null;
    } catch (err) {
      console.error('Error updating candidate status:', err);
    } finally {
      isUpdatingStatus = false;
    }
  }

  function handleStatusFilter(status: string) {
    currentStatusFilter = status;
    fetchCandidates(status).then((data) => {
      candidates = data;
      selectedCandidate = data.length > 0 ? data[0] : null;
    });
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
      candidates = await fetchCandidates();
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

<div class="pt-8 h-[calc(100vh-80px)] flex flex-col">
  <h1 class="text-2xl font-semibold text-gray-900">
    All Applicants
  </h1>
  <p class="my-3 text-base text-gray-600">
    Manage candidates across all job postings
  </p>
  
  <div class="flex flex-col gap-3 mb-6">
    <div class="flex items-center bg-white rounded-lg shadow px-3 py-2 max-w-md">
      <Search class="w-5 h-5 text-gray-500" />
      <input
        type="text"
        placeholder="Search candidates by name, email, skills or position..."
        class="flex-1 ml-2 outline-none border-none text-sm"
      />
    </div>

    <div class="flex gap-2">      
      <select class="pr-8 pl-3 py-1 bg-white border border-gray-200 rounded text-sm text-left">
        <option>All Jobs</option>
        <option>AI Researcher</option>
        <option>Developer</option>
        <option>Tester</option>
      </select>

      <select
        class="pr-8 pl-3 py-1 bg-white border border-gray-200 rounded text-sm text-left"
        on:change={(e) => handleStatusFilter((e.target as HTMLSelectElement).value.toUpperCase())}
      >
        <option value="">All Statuses</option>
        <option value="Accepted">Accepted</option>
        <option value="Pending">Pending</option>
        <option value="Rejected">Rejected</option>
      </select>

      <div class="border-l border-gray-300 h-6 py-1"></div>
        <button class="px-3 py-1 bg-white border-1 border-gray-200 rounded-full text-sm">All</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">New ({pendingCount})</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">High Match</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Recent Grads</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Experienced</button>
      </div>
  </div>

  <div class="grid grid-cols-3 gap-4 flex-1 overflow-hidden">
    <div class="col-span-1 rounded-xl border border-gray-200 overflow-hidden flex flex-col">
      <div class="bg-slate-50 p-4 border-b border-gray-200">
        <h2 class="text-lg font-semibold text-gray-800">All Candidates</h2>
        <div class="flex gap-1 mt-2">
          <button class="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-xs">All</button>
          <button class="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-xs">New</button>
          <button class="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-xs">Shortlisted</button>
          <button class="px-3 py-1.5 bg-white border border-gray-200 rounded-full text-xs">Rejected</button>
        </div>
      </div>
      
      <div class="bg-white p-4 space-y-3 flex-1 overflow-y-auto">
        {#each paginatedCandidates as candidate, i (i)}
          <button
            class="flex items-start gap-4 p-4 border rounded-xl bg-white hover:bg-gray-50 shadow-sm hover:shadow-lg transition-all duration-200 cursor-pointer w-full
              {selectedCandidate?.name === candidate.name ? 'bg-gradient-to-r from-green-50 to-emerald-50 border-green-500 shadow-green-200/50' : 'border-gray-200'}"
            on:click={() => (selectedCandidate = candidate)}
          >
            <div class="relative">
              <img src={candidate.avatar} alt={candidate.name} class="w-12 h-12 rounded-full object-cover ring-2 ring-white shadow-sm" />
              <div class="absolute -bottom-1 -right-1 w-3 h-3 bg-green-500 rounded-full border-2 border-white"></div>
            </div>
            <div class="flex-1 flex flex-col items-start gap-1.5">
              <h3 class="font-semibold text-gray-900 text-sm tracking-tight">{candidate.name}</h3>
              <p class="text-start text-xs font-medium text-gray-700">{candidate.role}</p>
              <p class="text-start text-[11px] text-gray-500 font-light">Applied for: <span class="font-medium text-gray-600">{candidate.applied}</span></p>
              <div class="flex items-center text-[10px] text-gray-500 gap-3 mt-1">
                <span class="flex items-center gap-1">
                  <Briefcase class="w-3.5 h-3.5 text-gray-400" />
                  <span class="font-medium">{candidate.year || 'N/A'}</span>
                </span>
                <span class="flex items-center gap-1">
                  <Clock class="w-3.5 h-3.5 text-gray-400" />
                  <span class="font-medium">{candidate.time}</span>
                </span>
              </div>
              <span
                class="inline-block mt-2 text-[11px] px-2.5 py-1 rounded-full font-semibold tracking-wide
                  {candidate.status === 'Pending' ? 'bg-blue-50 text-blue-700 border border-blue-200' : ''}
                  {candidate.status === 'Accepted' ? 'bg-green-50 text-green-700 border border-green-200' : ''}
                  {candidate.status === 'Shortlisted' ? 'bg-purple-50 text-purple-700 border border-purple-200' : ''}
                  {candidate.status === 'Rejected' ? 'bg-red-50 text-red-700 border border-red-200' : ''}
                  {candidate.status === 'Reviewed' ? 'bg-amber-50 text-amber-700 border border-amber-200' : ''}"
              >
                {candidate.status}
              </span>
            </div>
          </button>
        {/each}
        <div class="flex justify-center items-center gap-2 mt-4">
          <button
            on:click={() => changePage(currentPage - 1)}
            class="px-3 py-1 bg-gray-200 border border-gray-300 rounded disabled:opacity-50"
            disabled={currentPage === 1}
          >
            Prev
          </button>

          <span class="text-sm text-gray-600">
            Page {currentPage} of {totalPages}
          </span>

          <button
            on:click={() => changePage(currentPage + 1)}
            class="px-3 py-1 bg-gray-200 border border-gray-300 rounded disabled:opacity-50"
            disabled={currentPage === totalPages}
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <div class="col-span-2 bg-white rounded-xl border border-gray-200 overflow-y-auto">
      {#if selectedCandidate}
        <!-- Profile Card with integrated header -->
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
          onAccept={(id) => updateStatus(id, 'ACCEPTED')}
          onReject={(id) => updateStatus(id, 'REJECTED')}
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


