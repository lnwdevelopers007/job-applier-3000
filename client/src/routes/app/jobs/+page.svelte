<script>
  import { Search, ArrowUpDown, MapPin, Users, Building } from 'lucide-svelte';
  import SafeHTML from '$lib/utils/SafeHTML.svelte';
  import { onMount } from 'svelte';
  import { jobSearchStore } from '$lib/stores/jobSearch';
  import { get } from 'svelte/store';
  import { isAuthenticated, getUserInfo } from '$lib/utils/auth';

  let jobs = [];
  let filteredJobs = [];
  let selectedJob = null;
  let searchQuery = "";
  let currentPage = 1;
  const pageSize = 4;

  let userInfo = null;
  let companyInfo = null;
  let isLoggedIn = false;
  let appliedJobs = new Set();

  $: totalPages = Math.ceil(filteredJobs.length / pageSize);
  $: paginatedJobs = filteredJobs.slice((currentPage - 1) * pageSize, currentPage * pageSize);
  $: if (selectedJob?.companyID) {
    fetchCompanyInfo(selectedJob.companyID);
  }

  let activeFilters = {
    type: null,
    posted: null,
    arrangement: null,
  };
  let sortBy = null;

  const typeCycle = ["Full-time", "Part-time", "Contract", "Casual"];
  const arrangementCycle = ["On-site", "Remote", "Hybrid"];
  const sortOptions = [null, "dateAsc", "dateDesc", "title"];
  const sortLabels = {
    null: "Sort",
    dateAsc: "Oldest",
    dateDesc: "Newest",
    title: "Title A-Z"
  };

  async function fetchAppliedJobs() {
    if (!isLoggedIn || !userInfo?.userID) return;

    try {
      const res = await fetch(`/apply/query?applicantID=${userInfo.userID}`, {
        credentials: 'include',
      });
      if (!res.ok) throw new Error('Failed to fetch applied jobs');

      const data = await res.json();
      appliedJobs = new Set(data.map(app => app.jobApplication.jobID));
    } catch (err) {
      console.error('Error fetching applied jobs:', err);
      appliedJobs = new Set();
    }
  }

  async function fetchCompanyInfo(companyID) {
    try {
      const res = await fetch(`/users/query?id=${companyID}`);
      if (!res.ok) throw new Error("Failed to fetch company info");

      const data = await res.json();
      const raw = data[0];
      if (!raw) {
        companyInfo = null;
        return;
      }

      const infoArray = raw.userInfo || [];
      const info = Object.fromEntries(infoArray.map(item => [item.Key, item.Value]));

      companyInfo = {
        name: info.name || raw.name || "Unknown Company",
        logo: info.logo || raw.avatarURL || "https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000",
        location: info.headquarters || "N/A",
        website: info.website || "",
        aboutUs: info.aboutUs || "",
        industry: info.industry || "",
        size: info.size || "",
        foundedYear: info.foundedYear || "",
        linkedIn: info.linkedIn || ""
      };
    } catch (err) {
      console.error("Error fetching company info:", err);
      companyInfo = null;
    }
  }

  async function fetchJobs(query = "", filters = {}, sort = null) {
    try {
      const params = new URLSearchParams();
      if (query) params.set("title", query);
      if (filters.type) params.set("workType", filters.type);
      if (filters.posted) params.set("postOpenDate", filters.posted);
      if (filters.arrangement) params.set("workArrangement", filters.arrangement);
      if (sort) params.set("sort", sort);
      const res = await fetch(`/jobs/query?${params.toString()}`);

      if (res.status === 404) {
        jobs = [];
        filteredJobs = [];
        selectedJob = null;
        return;
      }

      if (!res.ok) throw new Error(`Failed to load jobs: ${res.status}`);
      const data = await res.json();

      const jobPromises = data.map(async (job) => {
        let companyName = "Unknown Company";
        let companyLogo = "https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000";

      if (job.companyID) {
        try {
          const companyRes = await fetch(`/users/query?id=${job.companyID}`);
          if (companyRes.ok) {
            const companyData = await companyRes.json();
            if (Array.isArray(companyData) && companyData.length > 0) {
              const company = companyData[0];
              const infoArray = company.userInfo || [];
              const info = Object.fromEntries(infoArray.map(item => [item.Key, item.Value]));

              companyName = info.name || company.name || companyName;
              companyLogo = info.logo || company.avatarURL || companyLogo;
            }
          }
        } catch (err) {
          console.warn(`Failed to load company info for ID ${job.companyID}:`, err);
        }
      }

        return {
          id: job.id,
          title: job.title,
          company: companyName,
          companyID: job.companyID,
          location: job.location || "N/A",
          type: job.workType || "Full-time",
          tags: job.requiredSkills ? job.requiredSkills.split(",").map(skill => skill.trim()) : [],
          posted: job.postOpenDate ? new Date(job.postOpenDate).toLocaleDateString() : "Unknown",
          closeDate: job.applicationDeadline ? new Date(job.applicationDeadline).toLocaleDateString() : "Unknown",
          description: job.jobDescription || "No description provided.",
          logo: companyLogo,
        };
      });

      jobs = await Promise.all(jobPromises);
      filteredJobs = jobs;
      selectedJob = jobs[0] || null;
      currentPage = 1;
    } catch (err) {
      console.error("Error fetching jobs:", err);
      jobs = [];
      filteredJobs = [];
      selectedJob = null;
    }
  }

  async function applyJob(job) {
    try {
      if (!isLoggedIn || !userInfo.userID) {
        alert("❌ You must be logged in to apply.");
        return;
      }

      if (appliedJobs.has(job.id)) {
        alert("⚠️ You have already applied to this job.");
        return;
      }

      const payload = {
        applicantID: userInfo.userID,
        jobID: job.id,
        companyID: job.companyID,
        status: "PENDING",
        createdAt: new Date().toISOString(),
      };

      const res = await fetch("/apply", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        if (res.status === 409) {
          alert("⚠️ You already applied to this job.");
        } else {
          throw new Error(`Failed to apply: ${res.status}`);
        }
        return;
      }

      appliedJobs.add(job.id);
      alert(`✅ Successfully applied to ${job.title}`);
    } catch (err) {
      console.error("Error applying to job:", err);
      alert("❌ Failed to apply. Please try again.");
    }
  }

  function toggleCycle(field) {
    if (field === "type") {
      let idx = typeCycle.indexOf(activeFilters.type);
      activeFilters.type = idx === -1 ? typeCycle[0] : (idx + 1 < typeCycle.length ? typeCycle[idx + 1] : null);
    } else if (field === "arrangement") {
      let idx = arrangementCycle.indexOf(activeFilters.arrangement);
      activeFilters.arrangement = idx === -1 ? arrangementCycle[0] : (idx + 1 < arrangementCycle.length ? arrangementCycle[idx + 1] : null);
    }
    activeFilters = { ...activeFilters };
    fetchJobs(searchQuery, activeFilters);
  }

  function toggleFilter(type, value) {
    activeFilters[type] = activeFilters[type] === value ? null : value;
    activeFilters = { ...activeFilters };
    fetchJobs(searchQuery, activeFilters, sortBy);
  }

  function toggleSort() {
    let idx = sortOptions.indexOf(sortBy);
    sortBy = idx === -1 || idx + 1 === sortOptions.length ? sortOptions[0] : sortOptions[idx + 1];
    fetchJobs(searchQuery, activeFilters, sortBy);
  }

  function onSearchInput() {
    fetchJobs(searchQuery, activeFilters);
  }

  onMount(() => {
    isLoggedIn = isAuthenticated();
    if (isLoggedIn) {
      userInfo = getUserInfo();
      fetchAppliedJobs();
    }

    const unsubscribe = jobSearchStore.subscribe(state => {
      if (state.shouldFetch) {
        searchQuery = state.query;
        fetchJobs(state.query, activeFilters).then(() => {
          jobSearchStore.clearFetchFlag();
        });
      }
    });

    const currentState = get(jobSearchStore);
    if (!currentState.shouldFetch) {
      fetchJobs();
    }

    return () => {
      unsubscribe();
    };
  });
</script>

<div class="flex h-screen">
  <main class="flex-1 p-0.5 space-y-6 w-full min-w-0">
    <!-- Search and Filters -->
    <div class="flex items-center gap-3">
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search jobs..."
          class="flex-1 ml-2 outline-none border-none"
          bind:value={searchQuery}
          oninput={onSearchInput}
        />
      </div>

      <div class="flex gap-2">
        <button
          class={`px-3 py-1 rounded-full text-sm ${activeFilters.type && typeCycle.includes(activeFilters.type) ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
          onclick={() => toggleCycle("type")}
        >
          {activeFilters.type || "Work Type"}
        </button>

        <button
          class={`px-3 py-1 rounded-full text-sm ${activeFilters.posted === '1d' ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
          onclick={() => toggleFilter("posted", "1d")}
        >
          1 day ago
        </button>

        <button
          class={`px-3 py-1 rounded-full text-sm ${activeFilters.posted === '6w' ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
          onclick={() => toggleFilter("posted", "6w")}
        >
          6 weeks
        </button>

        <button
          class={`px-3 py-1 rounded-full text-sm ${activeFilters.arrangement && arrangementCycle.includes(activeFilters.arrangement) ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
          onclick={() => toggleCycle("arrangement")}
        >
          {activeFilters.arrangement || "Arrangement"}
        </button>

        <button 
          class={`px-3 py-1 rounded-full text-sm flex items-center ${sortBy ? 'bg-green-600 text-white' : 'bg-gray-200 text-black'}`} 
          onclick={toggleSort}
        >
          <ArrowUpDown class="h-4 w-4 mr-1" />
          {sortBy ? sortLabels[sortBy] : "Sort"}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 w-full">
      <!-- Jobs List -->
      <section class="col-span-1 flex flex-col space-y-4">
        <div class="flex-1 overflow-y-auto w-full space-y-3 p-2">
          {#if paginatedJobs.length === 0}
            <div class="text-center text-gray-500 mt-4">
              No jobs match your search or filters.
            </div>
          {:else}
            {#each paginatedJobs as job (job.id)}
              <button
                onclick={() => (selectedJob = job)}
                class={`w-full flex flex-col cursor-pointer bg-white rounded-lg shadow ring-offset-2 hover:ring-2 hover:ring-green-300 ${selectedJob?.id === job.id ? 'ring-2 ring-green-500' : ''}`}
              >
                <div class="flex items-start gap-3 p-2">
                  <img src={job.logo} alt={job.company} class="w-12 h-12 rounded-full object-cover flex-shrink-0 mt-2" />
                  <div class="flex flex-col text-left flex-1">
                    <div class="font-semibold">{job.title}</div>
                    <div class="text-sm text-gray-600">{job.company}</div>
                    <div class="text-sm text-black flex items-center gap-2 mt-1">
                      <MapPin class="w-4 h-4 mt-0.5" />
                      {job.location}
                    </div>
                  </div>
                </div>

                <div class="flex flex-wrap gap-2 mt-2 p-2">
                  {#each job.tags as tag, i (i)}
                    <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
                  {/each}
                </div>

                <div class="text-xs text-left mx-3 pb-2 text-gray-500 mt-1">{job.posted}</div>
              </button>
            {/each}
          {/if}
        </div>

        <!-- Pagination controls -->
        {#if totalPages > 1}
          <div class="flex justify-center gap-2 my-4">
            <button class="px-3 py-1 rounded bg-gray-200 disabled:opacity-50"
              onclick={() => currentPage = Math.max(1, currentPage - 1)}
              disabled={currentPage === 1}>
              Prev
            </button>

            <span class="px-3 py-1 text-sm">
              Page {currentPage} of {totalPages}
            </span>

            <button class="px-3 py-1 rounded bg-gray-200 disabled:opacity-50"
              onclick={() => currentPage = Math.min(totalPages, currentPage + 1)}
              disabled={currentPage === totalPages}>
              Next
            </button>
          </div>
        {/if}
      </section>

      <!-- Job Detail -->
      <section class="col-span-1 bg-white p-6 rounded-lg shadow space-y-2 min-h-[400px] flex flex-col">
        {#if selectedJob}
          <!-- Job Info -->
          <div class="mb-2">
            <div class="flex items-start gap-3 px-2">
              <img src={selectedJob.logo} alt={selectedJob.company} class="w-4 h-4 rounded-full object-cover flex-shrink-0 mt-0.5" />
              <p class="text-sm font-semibold text-gray-600">{selectedJob.company}</p>
            </div>

            <h2 class="text-xl font-bold">{selectedJob.title}</h2>

            <div class="text-sm text-black flex items-center gap-2 my-2">
              <MapPin class="w-4 h-4 mt-0.5" />
              {selectedJob.location}
            </div>

            <div class="text-sm text-gray-500 flex gap-4 my-2">
              <span>Open: {selectedJob.posted}</span>
              <span>Close: {selectedJob.closeDate}</span>
            </div>
          </div>

          <!-- Tags -->
          <div class="flex gap-2">
            {#each selectedJob.tags as tag, i (i)}
              <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
            {/each}
          </div>

          <!-- Apply / Bookmark -->
          <div class="flex gap-2 mb-4">
            {#if new Date(selectedJob.closeDate) < new Date()}
              <button class="px-4 py-2 bg-gray-400 text-white text-sm rounded-lg cursor-not-allowed" disabled>
                Closed
              </button>
            {:else if new Date(selectedJob.posted) > new Date()}
              <button class="px-4 py-2 bg-gray-400 text-white text-sm rounded-lg cursor-not-allowed" disabled>
                Not Open Yet
              </button>
            {:else}
              <button class="px-4 py-2 bg-green-600 text-white text-sm rounded-lg" onclick={() => applyJob(selectedJob)}>
                Apply
              </button>
            {/if}

            <button class="px-4 py-2 bg-yellow-400 text-sm rounded-lg">Bookmark</button>
          </div>

          <!-- Job Description -->
          <h3 class="font-semibold text-lg">Job Description</h3>
          <div class="space-y-4">
            <SafeHTML html={selectedJob.description} />
          </div>

        <h3 class="font-semibold text-lg">About the company</h3>
        {#if companyInfo}
          <div class="mt-2 p-4 bg-gray-50 rounded-lg shadow space-y-2">
            <div class="flex items-start gap-4">
              <img
                src={companyInfo.logo || "https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000"}
                alt={companyInfo.name}
                class="w-16 h-16 rounded-full object-cover flex-shrink-0"
              />
              <div class="flex-1 flex flex-col space-y-1">
                <h3 class="font-semibold text-lg">{companyInfo.name}</h3>

                {#if companyInfo.location}
                  <p class="text-sm text-gray-600 flex gap-1">
                    <MapPin class="w-4 h-4 mt-0.5" /> {companyInfo.location}
                  </p>
                {/if}
              </div>
            </div>

            {#if companyInfo.size || companyInfo.industry}
              <div class="flex gap-6 text-sm text-gray-700">
                {#if companyInfo.size}
                  <span class="flex items-center gap-1">
                    <Users class="w-4 h-4 mt-0.5" /> {companyInfo.size}
                  </span>
                {/if}
                {#if companyInfo.industry}
                  <span class="flex items-center gap-1">
                    <Building class="w-4 h-4 mt-0.5" /> {companyInfo.industry}
                  </span>
                {/if}
              </div>
            {/if}

            {#if companyInfo.aboutUs}
              <p class="text-sm text-gray-700 mt-1">{companyInfo.aboutUs}</p>
            {/if}
          </div>
        {:else}
          <div class="text-gray-500 mt-4">Company information not available.</div>
        {/if}

        {:else}
          <div class="flex-1 flex items-center justify-center text-gray-500">
            No job selected.
          </div>
        {/if}
      </section>
    </div>
  </main>
</div>
