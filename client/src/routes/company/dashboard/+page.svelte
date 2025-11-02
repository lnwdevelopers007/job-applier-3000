<script>
  import { goto } from '$app/navigation';
	import ConfirmActionWithReason from '$lib/components/modals/ConfirmActionWithReason.svelte';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';
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

<div>
  <h1 class="text-2xl font-semibold text-gray-900 mb-1">
    Company Dashboard
  </h1>
  <p class="mb-6 text-base text-gray-600">
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

  <TableWithAction
    things={jobs}
    tableHeader={['Job Title', 'Status', 'Applicants', 'Views', 'Posted', 'Expires']}
		rowAttributes={['title', 'status', 'applicants', 'views', 'posted', 'expires']}
		{handleAction}
  />

</div>

<ConfirmActionWithReason
  bind:isVisible={showDeleteModal}
  actionName={'Delete'}
  actOnKind={'Job'}
  actOnIndividual={jobToDelete?.title}
  bind:isActionInProgress={deleting}
  reasonForAction={deleteReason}
  action={confirmDelete}
/>

