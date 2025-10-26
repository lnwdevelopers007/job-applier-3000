<script>
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth.svelte';
  import { apiFetch } from '$lib/utils/api';
  import { getCompanyAnalytics } from '$lib/utils/companyStats';

  let jobs = $state([]);
  let selectedJob = $state(null);

  // Modal states
  let showDeleteModal = $state(false);
  let deleteReason = $state('');
  let jobToDelete = $state(null);
  let deleting = $state(false);
  let stats = $state({
    activeJobs: 0,
    totalApplicants: 0,
    pendingReview: 0,
    offersAccepted: 0,
    trend: {
      activeJobs: 0,
      totalApplicants: 0,
      pendingReview: 0,
      offersAccepted: 0
    }
  });

  function selectRow(index) {
    selectedJob = selectedJob === index ? null : index;
  }

  async function handleAction(action, job) {
    if (action.label === 'Edit') {
      goto(`/company/edit/${job.id}`);
    }
    if (action.label === 'View') {
      goto(`/app/jobs/${job.id}`);
    }
    if (action.label === 'Delete') {
      jobToDelete = job;
      deleteReason = '';
      showDeleteModal = true;
    }
  }

  async function confirmDelete() {
    if (!jobToDelete) return;
    deleting = true;
    try {
      const res = await apiFetch(`/jobs/${jobToDelete.id}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ reason: deleteReason })
      });

      if (!res.ok) throw new Error(`Failed to delete job: ${res.status}`);

      jobs = jobs.filter(j => j.id !== jobToDelete.id);
    } catch (err) {
      console.error(err);
    } finally {
      deleting = false;
      showDeleteModal = false;
      jobToDelete = null;
    }
  }


  // Wait for auth to be ready before loading data
  $effect(() => {
    if (authStore.isAuthenticated && authStore.user) {
      loadCompanyData();
    }
  });

  async function loadCompanyData() {
    try {
      const user = authStore.user;
      if (!user?.userID) {
        goto('/login');
        return;
      }
      const companyID = user.userID;
      stats = await getCompanyAnalytics(companyID);
      const res = await apiFetch(`/jobs/query?companyID=${companyID}`);

      if (!res.ok) {
        if (res.status === 404) {
          jobs = [];
          return;
        }
        throw new Error(`Failed to load jobs: ${res.status}`);
      }

      const data = await res.json();

      const jobWithCounts = await Promise.all(
        data.map(async (job) => {
          const jobID = job.id;
          let applicantCount = 0;

          try {
            const applyRes = await apiFetch(`/apply?jobID=${jobID}`);
            if (applyRes.ok) {
              const applications = await applyRes.json();
              if (Array.isArray(applications)) applicantCount = applications.length;
            }
          } catch (err) {
            console.warn(`Failed to fetch applicants for job ${jobID}`, err);
          }

          const now = new Date();
          const deadline = job.applicationDeadline ? new Date(job.applicationDeadline) : null;
          const isClosed = deadline && deadline < now;
          const status = isClosed ? 'Closed' : 'Active';

          return {
            id: job.id,
            title: job.title,
            status: status || 'Active',
            applicants: applicantCount,
            views: job.views || 0,
            posted: job.postOpenDate
              ? new Date(job.postOpenDate).toLocaleDateString()
              : 'None',
            expires: job.applicationDeadline
              ? new Date(job.applicationDeadline).toLocaleDateString()
              : 'None',
            actions: [
              { label: 'View', disabled: false },
              { label: 'Edit', disabled: false },
              { label: 'Delete', disabled: false }
            ]
          };
        })
      );

      jobs = jobWithCounts;
    } catch {
      jobs = [];
    }
  }
</script>

<div class="p-2 bg-slate-50">
  <h1 class="text-2xl font-bold text-gray-900">
    Company Dashboard
  </h1>
  <p class="my-2 text-base text-gray-600">
    Welcome back, company HR Team. Here's what’s happening with your recruitment.
  </p>

  <!-- Summary cards -->
   <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Active Jobs</h2>
      <p class="text-gray-800 text-2xl font-medium">{stats.activeJobs}</p>
      <p class={stats.trend.activeJobs >= 0 ? 'text-green-500 text-sm' : 'text-red-500 text-sm'}>
        {stats.trend.activeJobs >= 0 ? '↑' : '↓'} {Math.abs(stats.trend.activeJobs)} from lastmonth
      </p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Total Applicants</h2>
      <p class="text-gray-800 text-2xl font-medium">{stats.totalApplicants}</p>
      <p class={stats.trend.totalApplicants >= 0 ? 'text-green-500 text-sm' : 'text-red-500 text-sm'}>
        {stats.trend.totalApplicants >= 0 ? '↑' : '↓'} {Math.abs(stats.trend.totalApplicants)} from last month
      </p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Pending Review</h2>
      <p class="text-gray-800 text-2xl font-medium">{stats.pendingReview}</p>
      <p class={stats.trend.pendingReview >= 0 ? 'text-green-500 text-sm' : 'text-red-500 text-sm'}>
        {stats.trend.pendingReview >= 0 ? '↑' : '↓'} {Math.abs(stats.trend.pendingReview)} new today
      </p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Offers Accepted</h2>
      <p class="text-gray-800 text-2xl font-medium">{stats.offersAccepted}</p>
      <p class={stats.trend.offersAccepted >= 0 ? 'text-green-500 text-sm' : 'text-red-500 text-sm'}>
        {stats.trend.offersAccepted >= 0 ? '↑' : '↓'} {Math.abs(stats.trend.offersAccepted)} from last month
      </p>
    </div>
  </div>

  <!-- Buttons -->
  <div class="mt-8 flex space-x-4">
    <button class="px-4 py-2 bg-gray-100 text-gray-700 border-1 border-gray-400 font-medium rounded hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:bg-gray-300">
      Posted Jobs
    </button>
    <button class="px-4 py-2 bg-gray-100 text-gray-700 border-1 border-gray-400 font-medium rounded hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:bg-gray-300">
      Recent Applicants
    </button>
  </div>

  <!-- Jobs Table -->
  <div class="mt-6 overflow-x-auto bg-white shadow rounded-lg">
    <table class="min-w-full text-sm text-left text-gray-600">
      <thead class="bg-gray-100 text-gray-700 text-sm font-semibold">
        <tr>
          <th class="px-4 py-3">Job Title</th>
          <th class="px-4 py-3">Status</th>
          <th class="px-4 py-3">Applicants</th>
          <th class="px-4 py-3">Views</th>
          <th class="px-4 py-3">Posted</th>
          <th class="px-4 py-3">Expires</th>
          <th class="px-4 py-3">Actions</th>
        </tr>
      </thead>
      <tbody>
        {#each jobs as job, i (i)}
          <tr 
            class="border-t cursor-pointer {selectedJob === i ? 'bg-green-200' : ''}" 
            onclick={() => selectRow(i)}
          >
            <td class="px-4 py-3 font-medium text-gray-900">{job.title}</td>
            <td class="px-4 py-3">{job.status}</td>
            <td class="px-4 py-3">{job.applicants}</td>
            <td class="px-4 py-3">{job.views}</td>
            <td class="px-4 py-3">{job.posted}</td>
            <td class="px-4 py-3">{job.expires}</td>
            <td class="px-4 py-3 space-x-2">
              {#each job.actions as action, j (j)}
                <button
                  class="px-3 py-1 text-sm rounded bg-gray-100 text-gray-900 border-gray-500 border-1 disabled:opacity-50 disabled:cursor-not-allowed enabled:hover:bg-gray-300"
                  disabled={action.disabled}
                  onclick={(e) => {
                    e.stopPropagation();
                    handleAction(action, job);
                  }}
                >
                  {action.label}
                </button>
              {/each}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<!-- Delete Reason Modal -->
{#if showDeleteModal}
  <div class="fixed inset-0 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-lg w-96 p-6">
      <h2 class="text-lg font-semibold text-gray-800 mb-4">
        Delete Job: {jobToDelete?.title}
      </h2>
      <label class="block text-sm text-gray-600 mb-2">Reason for deletion:</label>

      <textarea
        bind:value={deleteReason}
        input={(e) => {
          if (e.target.value.length > 500) {
            deleteReason = e.target.value.slice(0, 500);
          } else {
            deleteReason = e.target.value;
          }
        }}
        rows="4"
        maxlength="500"
        placeholder="Enter reason..."
        class="w-full border border-gray-300 rounded-md p-2 text-gray-800 focus:ring-2 focus:ring-red-400 focus:outline-none"
      ></textarea>

      <div class="flex justify-between items-center mt-1 text-sm">
        <span class="{deleteReason.length >= 480 ? 'text-red-500' : 'text-gray-500'}">
          {deleteReason.length}/500 characters
        </span>
        {#if deleteReason.length >= 500}
          <span class="text-red-600 font-medium">Limit reached</span>
        {/if}
      </div>

      <div class="mt-4 flex justify-end space-x-3">
        <button
          class="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200"
          onclick={() => (showDeleteModal = false)}
        >
          Cancel
        </button>
        <button
          class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
          disabled={!deleteReason.trim() || deleting}
          onclick={confirmDelete}
        >
          Confirm Delete
        </button>
      </div>
    </div>
  </div>
{/if}
