<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    GraduationCap,
    Briefcase,
    Clock,
    StickyNote,
    Check,
    X,
    User,
    CodeXml,
    FileText,
    Github,
    Link
  } from 'lucide-svelte';
  import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
  import { goto } from '$app/navigation';

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
    return {
      id: user.id,
      userID: user.userID,
      email: user.email || '',
      avatar: user.avatarURL || 'https://cdn-icons-png.flaticon.com/512/149/149071.png',
      name: user.name || info.fullName || 'Unknown User',
      role: info.desiredRole || 'Unknown Role',
      phone: info.phone || '-',
      linkedIn: info.linkedIn || '-',
      location: info.location || '-',
      github: info.github || '',
      aboutMe: info.aboutMe || '',
      education: user.education || [],
      skills: user.skills || [],
      documents: user.documents || []
    };
  }

  async function fetchCandidates(status: string = '') {
    if (!isAuthenticated()) {
      goto('/login');
      return [];
    }

    if (!company?.userID) {
      console.error('Company not found in localStorage');
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
        const user = normalizeUser(userData[0] || {});

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
          name: user.name,
          role: user.role,
          applied: app.jobTitle,
          status: displayStatus,
          avatar: user.avatar,
          time: daysAgo === 0 ? 'Today' : `${daysAgo} days ago`,
          email: user.email,
          phone: user.phone,
          address: user.location,
          linkedin: user.linkedIn,
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
    company = getUserInfo();
    candidates = await fetchCandidates();
    if (candidates.length > 0) selectedCandidate = candidates[0];
  });
</script>

<div class="p-2 bg-gray-100">
  <h1 class="text-2xl font-bold text-gray-900">
    All Applicants
  </h1>
  <p class="my-3 text-base text-gray-600">
    Manage candidates across all job postings
  </p>
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search candidates by name, email, skills or position..."
          class="flex-1 ml-2 outline-none border-none text-sm"
        />
      </div>

      <div class="flex gap-2 my-2">      
      <select class="pr-8 pl-3 py-1 bg-white border border-gray-200 rounded text-sm text-left">
        <option>All Jobs</option>
        <option>AI Researcher</option>
        <option>Developer</option>
        <option>Tester</option>
      </select>

      <select
        class="pr-8 pl-3 py-1 bg-white border border-gray-200 rounded text-sm text-left"
        on:change={(e) => handleStatusFilter(e.target.value.toUpperCase())}
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

  <div class="grid grid-cols-3 gap-4 mt-4">
    <div class="col-span-1 bg-white p-4 rounded-lg shadow">
      <h2 class="text-lg font-semibold text-gray-800">All Candidates</h2>
      <div class="flex gap-2 my-2">
        <button class="px-2 py-1 bg-white border-1 border-gray-200 rounded-full text-xs">All</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">New</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">Shortlisted</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">Rejected</button>
      </div>
    <div class="mt-3 space-y-3">
        {#each paginatedCandidates as candidate, i (i)}
          <button
            class="flex items-start gap-3 p-3 border rounded-lg shadow-sm hover:shadow-md transition cursor-pointer w-full
              {selectedCandidate?.name === candidate.name ? 'bg-green-100 border-green-600' : ''}"
            on:click={() => (selectedCandidate = candidate)}
          >
            <img src={candidate.avatar} alt={candidate.name} class="w-10 h-10 rounded-full object-cover" />
            <div class="flex-1 flex flex-col items-start">
              <h3 class="font-semibold text-gray-800 text-sm">{candidate.name}</h3>
              <p class="text-start text-xs text-gray-800">{candidate.role}</p>
              <p class="text-start text-[10px] text-gray-500">Applied for: {candidate.applied}</p>
              <p class="flex text-start text-[10px] text-gray-500 gap-0.5">
                <GraduationCap class="h-3 w-3 m-0.5" />{candidate.university}
                <Briefcase class="w-3 h-3 m-0.5" />{candidate.year}
                <Clock class="w-3 h-3 m-0.5" />{candidate.time}</p>
              <span
                class="inline-block mt-1 text-xs px-2 py-0.5 rounded-lg font-semibold
                  {candidate.status === 'Pending' ? 'bg-blue-100 text-blue-700 border-1 border-blue-700' : ''}
                  {candidate.status === 'Accepted' ? 'bg-green-100 text-green-700 border-1 border-green-700' : ''}
                  {candidate.status === 'Shortlisted' ? 'bg-purple-100 text-purple-700 border-1 border-purple-700' : ''}
                  {candidate.status === 'Rejected' ? 'bg-red-100 text-red-700 border-1 border-red-700' : ''}
                  {candidate.status === 'Reviewed' ? 'bg-yellow-100 text-yellow-700 border-1 border-yellow-700' : ''}"
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

    <div class="col-span-2 bg-white rounded-lg">
      {#if selectedCandidate}
        <div class="flex items-start gap-3 p-3 rounded-lg mt-3">
          <img src={selectedCandidate.avatar} alt={selectedCandidate.name} class="w-20 h-20 rounded-full object-cover" />
          <div class="flex-1 flex flex-col items-start">
            <div class="flex items-center justify-between w-full">
              <h3 class="font-semibold text-gray-800 text-lg">{selectedCandidate.name}</h3>
              <div class="flex gap-2 font-semibold">
                <button class="flex px-4 py-1 text-xs bg-gray-100 rounded border border-gray-300 hover:bg-gray-500"><StickyNote class="w-4 h-4 mx-1" />Notes</button>
                <button class="flex px-4 py-1 text-xs bg-green-500 text-white rounded border border-green-300 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
                on:click={() => updateStatus(selectedCandidate.id, 'ACCEPTED')}
                disabled={isUpdatingStatus || selectedCandidate.status === 'Accepted' || selectedCandidate.status === 'Rejected'}
                ><Check class="w-4 h-4 mx-1" />Accept</button>
                <button class="flex px-4 py-1 text-xs bg-red-500 text-white rounded border border-red-300 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
                on:click={() => updateStatus(selectedCandidate.id, 'REJECTED')}
                disabled={isUpdatingStatus || selectedCandidate.status === 'Accepted' || selectedCandidate.status === 'Rejected'}
                ><X class="w-4 h-4 mx-1" />Reject</button>
              </div>
            </div>
            <p class="text-sm text-gray-800">{selectedCandidate.role}</p>
            <p class="text-sm text-gray-500">Applied for: {selectedCandidate.applied}</p>
            <p class="flex text-sm text-gray-500 gap-1 mt-1">
              <GraduationCap class="h-4 w-4 m-0.5" />{selectedCandidate.university}
              <Briefcase class="w-4 h-4 m-0.5" />{selectedCandidate.year}
              <Clock class="w-4 h-4 m-0.5" />{selectedCandidate.time}
            </p>
            <span
              class="inline-block mt-2 text-xs px-2 py-0.5 rounded-lg font-semibold
                {selectedCandidate.status === 'Pending' ? 'bg-blue-100 text-blue-700 border-1 border-blue-700' : ''}
                {selectedCandidate.status === 'Accepted' ? 'bg-green-100 text-green-700 border-1 border-green-700' : ''}
                {selectedCandidate.status === 'Shortlisted' ? 'bg-purple-100 text-purple-700 border-1 border-purple-700' : ''}
                {selectedCandidate.status === 'Rejected' ? 'bg-red-100 text-red-700 border-1 border-red-700' : ''}
                {selectedCandidate.status === 'Reviewed' ? 'bg-yellow-100 text-yellow-700 border-1 border-yellow-700' : ''}"
            >
              {selectedCandidate.status}
            </span>
          </div>
        </div>
        <div class="p-4">       
          <hr class="w-full my-4 border-gray-300" />
            <h3 class="flex gap-1 font-semibold text-gray-800 text-lg my-4"><User class="mt-0.5" />Contact Information</h3>
            <div class="grid grid-cols-2 gap-4 w-full text-sm text-gray-700">
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Email</span>
                <span class="text-gray-800">{selectedCandidate.email || "alice.johnson@gmail.com"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Phone</span>
                <span class="text-gray-800">{selectedCandidate.phone || "+66 123 4567"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Address</span>
                <span class="text-gray-800">{selectedCandidate.address || "Bangkok, Thailand"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">LinkedIn</span>
                <span class="text-gray-800 break-all whitespace-normal">{selectedCandidate.linkedin || "None"}</span>
              </div>
            </div>
        </div>
        <div class="px-4 py-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2"><GraduationCap class="mt-0.5" />Education</h3>
          <div class="grid grid-cols-1 gap-2 w-full text-sm mt-3">
            {#if selectedCandidate.education}
              {#each selectedCandidate.education as edu, i (i)}
                <div class="flex flex-col">
                  <span class="font-semibold text-sm">{edu.degree}</span>
                  <span class="text-gray-800 mt-1">{edu.university} {edu.period}</span>
                  <span class="text-gray-800 mt-1">GPA {edu.gpa}</span>
                </div>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No education details provided</span>
            {/if}
          </div>
        </div>
        <div class="px-4 py-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2">
            <CodeXml class="mt-0.5" />Skills & Technologies
          </h3>
          <div class="flex flex-wrap gap-2 mt-3">
            {#if selectedCandidate.skills && selectedCandidate.skills.length > 0}
              {#each selectedCandidate.skills as skill, i (i)}
                <span class="px-2 py-1 text-xs bg-gray-200 text-gray-800 rounded-lg border-1 border-gray-300">{skill}</span>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No skills provided</span>
            {/if}
          </div>
        </div>
        <div class="px-4 py-2 mb-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2">
            <FileText class="mt-0.5" />Documents & Portfolio
          </h3>
          <div class="flex flex-col gap-4 mt-3">
            {#if selectedCandidate.documents && selectedCandidate.documents.length > 0}
              {#each selectedCandidate.documents as doc, i (i)}
                <a href={doc.name.startsWith('http') ? doc.name : '#'} target="_blank" class="flex items-center gap-4 bg-gray-100 border border-gray-300 rounded p-3 hover:bg-gray-200 transition">
                  {#await Promise.resolve(getDocIcon(doc.name)) then Icon}
                    <Icon class="w-8 h-8 text-gray-600" />
                  {/await}
                  <div class="flex flex-col">
                    <span class="text-sm font-semibold text-gray-800">{doc.name}</span>
                    <span class="text-xs text-gray-500">{doc.description}</span>
                  </div>
                </a>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No documents uploaded</span>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>


