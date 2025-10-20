<script>
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth.svelte';
  import { apiFetch } from '$lib/utils/api';

  let jobs = $state([]);
  let selectedJob = $state(null);

  function selectRow(index) {
    selectedJob = selectedJob === index ? null : index;
  }

  function handleAction(action, job) {
    if (action.label === 'Edit') {
      goto(`/company/edit/${job.id}`);
    }
    if (action.label === 'View') {
      goto(`/app/jobs/${job.id}`);
    }
    if (action.label === 'Manage') {
      // TODO: Implement manage job posting
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
      
      const res = await apiFetch(`/jobs/query?companyID=${companyID}`);
      
      if (!res.ok) {
        if (res.status === 404) {
          jobs = [];
          return;
        }
        throw new Error(`Failed to load jobs: ${res.status}`);
      }
      
      const data = await res.json();
      
      jobs = data.map(job => ({
        id: job.id,
        title: job.title,
        status: job.status || 'Active',
        applicants: job.applicants || 0,
        views: job.views || 0,
        posted: job.postOpenDate
          ? new Date(job.postOpenDate).toLocaleDateString()
          : "None",
        expires: job.applicationDeadline
          ? new Date(job.applicationDeadline).toLocaleDateString()
          : "None",
        actions: [
          { label: 'View', disabled: false },
          { label: 'Edit', disabled: false },
          { label: 'Manage', disabled: false }
        ]
      }));
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
  <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Active Jobs</h2>
      <p class="text-gray-800 text-xl font-bold">8</p>
      <p class="text-green-500 text-sm">↑ 2 from last month</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Total Applicants</h2>
      <p class="text-gray-800 text-xl font-bold">156</p>
      <p class="text-green-500 text-sm">↑ 23 this week</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Pending Review</h2>
      <p class="text-gray-800 text-xl font-bold">18</p>
      <p class="text-red-500 text-sm">↓ 5 from yesterday</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Offers Extended</h2>
      <p class="text-gray-800 text-xl font-bold">12.</p>
      <p class="text-red-500 text-sm">↓ 3 this week</p>
    </div>
  </div>
  <div class="mt-8 flex space-x-4">
    <button class="px-4 py-2 bg-gray-100 text-gray-700 border-1 border-gray-400 font-medium rounded hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:bg-gray-300">
      Posted Jobs
    </button>
    <button class="px-4 py-2 bg-gray-100 text-gray-700 border-1 border-gray-400 font-medium rounded hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:bg-gray-300">
      Recent Applicants
    </button>
  </div>
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
              {#each job.actions as action, i (i)}
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