<div class="p-2 bg-gray-100">
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
        {#each jobs as job}
          <tr class="border-t">
            <td class="px-4 py-3 font-medium text-gray-900">{job.title}</td>
            <td class="px-4 py-3">{job.status}</td>
            <td class="px-4 py-3">{job.applicants}</td>
            <td class="px-4 py-3">{job.views}</td>
            <td class="px-4 py-3">{job.posted}</td>
            <td class="px-4 py-3">{job.expires}</td>
            <td class="px-4 py-3 space-x-2">
              {#each job.actions as action}
                <button
                  class="px-3 py-1 text-sm rounded bg-gray-100 text-gray-900 border-gray-500 border-1 disabled:opacity-50 disabled:cursor-not-allowed enabled:hover:bg-gray-300"
                  disabled={action.disabled}
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
<script>
  let jobs = [
    {
      title: "AI Researcher",
      status: "Active",
      applicants: 34,
      views: 120,
      posted: "Aug 20, 2025",
      expires: "Sep 20, 2025",
      actions: [
        { label: "View", disabled: false },
        { label: "Edit", disabled: false },
        { label: "Manage", disabled: false }
      ]
    },
    {
      title: "Frontend Developer",
      status: "Paused",
      applicants: 12,
      views: 85,
      posted: "Aug 15, 2025",
      expires: "Sep 15, 2025",
      actions: [
        { label: "View", disabled: false },
        { label: "Edit", disabled: false },
        { label: "Manage", disabled: true }
      ]
    },
    {
      title: "QA Tester",
      status: "Closed",
      applicants: 48,
      views: 210,
      posted: "Jul 10, 2025",
      expires: "Aug 10, 2025",
      actions: [
        { label: "View", disabled: false },
        { label: "Edit", disabled: true },
        { label: "Manage", disabled: true }
      ]
    }
  ];
</script>