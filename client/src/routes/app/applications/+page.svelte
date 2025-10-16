<script lang="ts">
  import { onMount } from 'svelte';
  import { Download, MapPin, Eye, EllipsisVertical } from 'lucide-svelte';
  import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
  import { goto } from '$app/navigation';

  let applications = [];
  let activities = [];
  let filterStatus = 'All';

  async function fetchApplications(status: string = 'All') {
    if (!isAuthenticated()) {
      goto('/login');
      return [];
    }

    const user = getUserInfo();
    if (!user?.userID) {
      console.error('User not found');
      goto('/login');
      return [];
    }

    try {
      const queryParams = new URLSearchParams({ applicantID: user.userID });
      if (status !== 'All') queryParams.append('status', status.toUpperCase());

      const res = await fetch(`/apply?${queryParams.toString()}`, {
        credentials: 'include',
      });

      if (!res.ok) throw new Error('Failed to fetch applications');
      const appData = await res.json();

      const appPromises = appData.map(async (app) => {
        // Fetch job details
        const jobRes = await fetch(`/jobs/query?id=${app.jobApplication.jobID}`, {
          credentials: 'include',
        });
        if (!jobRes.ok) throw new Error('Failed to fetch job details');
        const jobData = await jobRes.json();
        const job = jobData[0];

        // Default values
        let companyName = 'Unknown Company';
        let companyLogo =
          'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8YnVpbGRpbmd8ZW58MHx8MHx8fDA%3D';

        // Fetch company details
        if (job?.companyID) {
          try {
            const companyRes = await fetch(`/users/query?id=${job.companyID}`, {
              credentials: 'include',
            });
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

        const created = new Date(app.jobApplication.createdAt);
        const diffDays = Math.floor(
          (Date.now() - created.getTime()) / (1000 * 60 * 60 * 24)
        );

        return {
          id: app.id,
          companyLogo,
          jobTitle: job?.title || 'Unknown Position',
          companyName,
          location: job?.location || 'N/A',
          daysAgo: diffDays,
          status: app.jobApplication.status.toUpperCase(),
        };
      });

      const apps = await Promise.all(appPromises);

      // Update activity timeline
      activities = apps.map((app) => ({
        time: app.daysAgo === 0 ? 'Today' : `${app.daysAgo} days ago`,
        title: app.jobTitle,
        description: `Your application to ${app.companyName} for ${app.jobTitle} is ${app.status}`,
      })).slice(0, 3);

      return apps;
    } catch (err) {
      console.error('Error fetching applications:', err);
      return [];
    }
  }

  onMount(async () => {
    applications = await fetchApplications('All');
  });

  // When user clicks a status filter button
  async function handleFilter(status: string) {
    filterStatus = status;
    applications = await fetchApplications(status);
  }
</script>

<div class="p-2 bg-slate-50">
  <h1 class="text-2xl font-bold text-gray-900">
    My Applications
  </h1>
  <p class="my-2 text-base text-gray-600">
    Track and manage all your job applications in one place
  </p>
  <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Total Applications</h2>
      <p class="text-gray-800 text-xl font-bold">24</p>
      <p class="text-green-500 text-sm">↑ 12% from last month</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">In Review</h2>
      <p class="text-gray-800 text-xl font-bold">8</p>
      <p class="text-green-500 text-sm">↑ 2 new this week</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Offer Received</h2>
      <p class="text-gray-800 text-xl font-bold">3</p>
      <p class="text-green-500 text-sm">↑ 2 pending response</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Response Rate</h2>
      <p class="text-gray-800 text-xl font-bold">67%</p>
      <p class="text-green-500 text-sm">↑ Above average</p>
    </div>
  </div>
  <div class="flex justify-between py-4">
    <div class="flex gap-2">
      {#each ['All', 'Pending', 'In Review', 'Accepted', 'Rejected'] as status, i (i)}
        <button
          class={`px-3 py-1 rounded-full text-sm border border-gray-300 transition
            ${filterStatus === status
              ? 'bg-green-600 text-white font-semibold border-blue-600'
              : 'bg-white text-gray-700 hover:bg-gray-100'}`}
            on:click={() => handleFilter(status)}
        >
          {status}
        </button>
      {/each}
    </div>
    <button class="px-3 py-1 bg-white rounded-lg font-semibold text-sm flex item-center gap-1 border-1 border-gray-300 "><Download class="w-5 h-5" /> Export</button>
  </div>

  <div class="p-4 w-full">
    <div class="bg-white shadow-md rounded-lg overflow-hidden border border-gray-300">
      <div class="p-3 bg-gray-100 border-b border-gray-300">
        <h1 class="px-2 font-bold">Recent Applications</h1>
      </div>

      {#if applications.length === 0}
        <div class="p-6 text-center text-gray-500">
          No applications found for "{filterStatus}"
        </div>
      {:else}
        {#each applications as app, i (i)}
          <div class={`flex items-center justify-between p-4 border-b border-gray-300 ${i === applications.length - 1 ? 'border-b-0' : ''}`}>
            <!-- Left part block -->
            <div class="flex items-center gap-4">
              <img src={app.companyLogo} alt="Company Logo" class="w-10 h-10 rounded-full object-cover" />
              <div class="flex flex-col">
                <span class="font-semibold">{app.jobTitle}</span>
                <div class="flex items-center gap-1.5 text-gray-500 text-sm">
                  <span>{app.companyName}</span>
                  <MapPin class="h-4 w-4 pl-1 text-black" />
                  <span>{app.location}</span>
                </div>
              </div>
            </div>

            <!-- Right part block -->
            <div class="flex items-center gap-4">
              <div class="flex flex-col items-end text-right">
                <span class="text-xs font-bold text-gray-600">APPLIED</span>
                <span class="text-sm text-gray-500">{app.daysAgo} days ago</span>
              </div>

              <span
                class={`badge text-xs font-semibold p-2 rounded-full border ${
                  app.status === 'IN REVIEW'
                    ? 'bg-yellow-100 text-yellow-800 border-yellow-300'
                    : app.status === 'ACCEPTED'
                    ? 'bg-green-100 text-green-800 border-green-300'
                    : app.status === 'PENDING'
                    ? 'bg-blue-100 text-blue-800 border-blue-300'
                    : app.status === 'REJECTED'
                    ? 'bg-red-100 text-red-800 border-red-300'
                    : 'badge-outline'
                }`}
              >
                {app.status}
              </span>

              <div class="flex items-center gap-2">
                <button class="btn btn-sm border border-gray-300 p-1 rounded-lg btn-ghost">
                  <Eye class="w-5 h-5" />
                </button>
                <button class="btn btn-sm border border-gray-300 p-1 rounded-lg btn-ghost">
                  <EllipsisVertical class="w-5 h-5" />
                </button>
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>


  <div class="mt-8">
    <h2 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h2>
    <div class="relative">
      <div
        class="absolute left-4 bg-gray-300 w-px"
        style="top: calc(2rem); bottom: calc(2rem);"
      ></div>

      {#if activities.length === 0}
        <div class="ml-10 p-4 text-gray-500">
          No recent activity for "{filterStatus}"
        </div>
      {:else}
        {#each activities as activity, i (i)}
          <div class="relative mb-8 last:mb-0">
            <span
              class="absolute left-4 top-1/2 w-5 h-5 rounded-full bg-white border-4 border-green-500 -translate-x-1/2 -translate-y-1/2"
            ></span>

            <div class="ml-10 bg-white shadow p-3 rounded-lg">
              <p class="text-xs text-gray-500 uppercase tracking-wide">{activity.time}</p>
              <p class="text-sm font-semibold text-gray-800">{activity.title}</p>
              <p class="text-sm text-gray-600">{activity.description}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
