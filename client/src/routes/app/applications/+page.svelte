<script lang="ts">
	import { onMount } from 'svelte';
	import { Download, MapPin, Eye, EllipsisVertical } from 'lucide-svelte';
	import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
	import { goto } from '$app/navigation';

	let applications: any[] = [];
	let activities: any[] = [];
	let filterStatus = 'All';
	const DEFAULT_COMPANY_LOGO =
		'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8M3x8YnVpbGRpbmd8ZW58MHx8MHx8fDA%3D';

	const DEFAULT_COMPANY_NAME = 'Unknown Company';

	async function fetchJob(app: any) {
		const jobRes = await fetch(`/jobs/${app.jobApplication.jobID}`, {
			credentials: 'include'
		});
		if (!jobRes.ok) throw new Error('Failed to fetch job details');
		const jobData = await jobRes.json();
		return jobData;
	}

	async function fetchCompany(companyID: string) {
		let companyName = DEFAULT_COMPANY_NAME;
		let companyLogo = DEFAULT_COMPANY_LOGO;

		try {
			const companyRes = await fetch(`/users/${companyID}`, {
				credentials: 'include'
			});
			if (companyRes.ok) {
				const company = await companyRes.json();
				const infoArray = company.userInfo || [];
				const info = Object.fromEntries(infoArray.map((item: any) => [item.Key, item.Value]));

				companyName = info.name || company.name || companyName;
				companyLogo = info.logo || company.avatarURL || companyLogo;
			}
		} catch (err) {
			console.warn(`Failed to load company info for ID ${companyID}:`, err);
		}

    return [companyName, companyLogo];
	}

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

			// fetch job application where status = status and applicantID = userID.
			const res = await fetch(`/apply/query?${queryParams.toString()}`, {
				credentials: 'include'
			});

			if (!res.ok) throw new Error('Failed to fetch applications');
			const appData = await res.json();

			const appPromises = appData.map(async (app: any) => {
				// Fetch job details
				const job = await fetchJob(app);

				// Default values
				let [companyName, companyLogo] = await fetchCompany(job.companyID); 

				const created = new Date(app.jobApplication.createdAt);
				const diffDays = Math.floor((Date.now() - created.getTime()) / (1000 * 60 * 60 * 24));

				return {
					id: app.id,
					companyLogo,
					jobTitle: job?.title || 'Unknown Position',
					companyName,
					location: job?.location || 'N/A',
					daysAgo: diffDays,
					status: app.jobApplication.status.toUpperCase()
				};
			});

			const apps = await Promise.all(appPromises);

			// Update activity timeline
			activities = apps
				.map((app) => ({
					time: app.daysAgo === 0 ? 'Today' : `${app.daysAgo} days ago`,
					title: app.jobTitle,
					description: `Your application to ${app.companyName} for ${app.jobTitle} is ${app.status}`
				}))
				.slice(0, 3);

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

<div class="bg-slate-50 p-2">
	<h1 class="text-2xl font-bold text-gray-900">My Applications</h1>
	<p class="my-2 text-base text-gray-600">
		Track and manage all your job applications in one place
	</p>
	<div class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<div class="rounded-lg bg-white p-4 shadow">
			<h2 class="text-gray-600">Total Applications</h2>
			<p class="text-xl font-bold text-gray-800">24</p>
			<p class="text-sm text-green-500">↑ 12% from last month</p>
		</div>

		<div class="rounded-lg bg-white p-4 shadow">
			<h2 class="text-gray-600">In Review</h2>
			<p class="text-xl font-bold text-gray-800">8</p>
			<p class="text-sm text-green-500">↑ 2 new this week</p>
		</div>

		<div class="rounded-lg bg-white p-4 shadow">
			<h2 class="text-gray-600">Offer Received</h2>
			<p class="text-xl font-bold text-gray-800">3</p>
			<p class="text-sm text-green-500">↑ 2 pending response</p>
		</div>

		<div class="rounded-lg bg-white p-4 shadow">
			<h2 class="text-gray-600">Response Rate</h2>
			<p class="text-xl font-bold text-gray-800">67%</p>
			<p class="text-sm text-green-500">↑ Above average</p>
		</div>
	</div>
	<div class="flex justify-between py-4">
		<div class="flex gap-2">
			{#each ['All', 'Pending', 'In Review', 'Accepted', 'Rejected'] as status, i (i)}
				<button
					class={`rounded-full border border-gray-300 px-3 py-1 text-sm transition
            ${
							filterStatus === status
								? 'border-blue-600 bg-green-600 font-semibold text-white'
								: 'bg-white text-gray-700 hover:bg-gray-100'
						}`}
					on:click={() => handleFilter(status)}
				>
					{status}
				</button>
			{/each}
		</div>
		<button
			class="item-center border-1 flex gap-1 rounded-lg border-gray-300 bg-white px-3 py-1 text-sm font-semibold"
			><Download class="h-5 w-5" /> Export</button
		>
	</div>

	<div class="w-full p-4">
		<div class="overflow-hidden rounded-lg border border-gray-300 bg-white shadow-md">
			<div class="border-b border-gray-300 bg-gray-100 p-3">
				<h1 class="px-2 font-bold">Recent Applications</h1>
			</div>

			{#if applications.length === 0}
				<div class="p-6 text-center text-gray-500">
					No applications found for "{filterStatus}"
				</div>
			{:else}
				{#each applications as app, i (i)}
					<div
						class={`flex items-center justify-between border-b border-gray-300 p-4 ${i === applications.length - 1 ? 'border-b-0' : ''}`}
					>
						<!-- Left part block -->
						<div class="flex items-center gap-4">
							<img
								src={app.companyLogo}
								alt="Company Logo"
								class="h-10 w-10 rounded-full object-cover"
							/>
							<div class="flex flex-col">
								<span class="font-semibold">{app.jobTitle}</span>
								<div class="flex items-center gap-1.5 text-sm text-gray-500">
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
								class={`badge rounded-full border p-2 text-xs font-semibold ${
									app.status === 'IN REVIEW'
										? 'border-yellow-300 bg-yellow-100 text-yellow-800'
										: app.status === 'ACCEPTED'
											? 'border-green-300 bg-green-100 text-green-800'
											: app.status === 'PENDING'
												? 'border-blue-300 bg-blue-100 text-blue-800'
												: app.status === 'REJECTED'
													? 'border-red-300 bg-red-100 text-red-800'
													: 'badge-outline'
								}`}
							>
								{app.status}
							</span>

							<div class="flex items-center gap-2">
								<button class="btn btn-sm btn-ghost rounded-lg border border-gray-300 p-1">
									<Eye class="h-5 w-5" />
								</button>
								<button class="btn btn-sm btn-ghost rounded-lg border border-gray-300 p-1">
									<EllipsisVertical class="h-5 w-5" />
								</button>
							</div>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>

	<div class="mt-8">
		<h2 class="mb-4 text-lg font-semibold text-gray-900">Recent Activity</h2>
		<div class="relative">
			<div
				class="absolute left-4 w-px bg-gray-300"
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
							class="absolute left-4 top-1/2 h-5 w-5 -translate-x-1/2 -translate-y-1/2 rounded-full border-4 border-green-500 bg-white"
						></span>

						<div class="ml-10 rounded-lg bg-white p-3 shadow">
							<p class="text-xs uppercase tracking-wide text-gray-500">{activity.time}</p>
							<p class="text-sm font-semibold text-gray-800">{activity.title}</p>
							<p class="text-sm text-gray-600">{activity.description}</p>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
