<script lang="ts">
	import { onMount } from 'svelte';
	import { MapPin, EllipsisVertical } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { fetchJob, fetchCompanyNameLogo } from '$lib/utils/fetcher';
	import Dropdown from '$lib/components/ui/Dropdown.svelte';
	import { formatRelativeTime, getDaysBetween } from '$lib/utils/datetime';
	import { toast } from 'svelte-french-toast';
	import Modal from '$lib/components/ui/Modal.svelte';
	import type { PageData } from './$types';
	import { getJobseekerStats, type JobseekerStats } from '$lib/utils/jobseekerStats';

	let { data }: { data: PageData } = $props();

	let applications = $state<any[]>([]);
	let activities = $state<any[]>([]);
	let filterStatus = $state('All');
	let openDropdown = $state<string | null>(null);
	let dropdownTriggers = $state<Record<string, HTMLElement>>({});
	let showAllApplications = $state(false);
	let stats: JobseekerStats = $state({
		totalApplications: 0,
		inReview: 0,
		offerReceived: 0,
		responseRate: 0,
		trend: {
			totalApplications: 0,
			inReview: 0,
			offerReceived: 0,
			responseRate: 0
		}
	});

	// Cancel modal state
	let showCancelModal = $state(false);
	let isCancelling = $state(false);
	let applicationToCancel = $state<any>(null);

	async function fetchApplications(status: string = 'All') {
		if (!data.user?.userID) {
			console.error('User not found');
			return [];
		}

		try {
			const queryParams = new URLSearchParams({ applicantID: data.user.userID });
			if (status !== 'All') queryParams.append('status', status.toUpperCase());

			stats = await getJobseekerStats(data.user.userID);
			// fetch job application where status = status and applicantID = userID.
			const res = await fetch(`/apply/query?${queryParams.toString()}`, {
				credentials: 'include'
			});

			if (!res.ok) throw new Error('Failed to fetch applications');
			const appData = await res.json();

			const appPromises = appData.map(async (app: any, index: number) => {
				// Handle different API response structures
				const jobApplication = app.jobApplication || app;

				// Fetch job details
				const job = await fetchJob(jobApplication.jobID);

				// fetch company name and logo.
				let [companyName, companyLogo] = await fetchCompanyNameLogo(job.companyID || '');

				const createdAt = jobApplication.createdAt || new Date().toISOString();
				const created = new Date(createdAt);
				const diffDays = getDaysBetween(created, new Date());

				return {
					id: `${app.id || jobApplication.jobID}_${index}_${Date.now()}`,
					jobId: jobApplication.jobID,
					jobApplication: jobApplication, // Keep reference for activity timeline
					companyLogo,
					jobTitle: job?.title || 'Unknown Position',
					companyName,
					location: job?.location || 'N/A',
					daysAgo: diffDays,
					status: jobApplication.status?.toUpperCase() || 'PENDING'
				};
			});

			const apps = await Promise.all(appPromises);

			if (status === 'All') {
				activities = apps
					.map((app) => {
						let eventTitle = '';
						let eventDescription = '';
						switch (app.status) {
							case 'PENDING':
								eventTitle = 'Application Submitted';
								eventDescription = `Your application for ${app.jobTitle} at ${app.companyName} was successfully submitted`;
								break;
							case 'IN REVIEW':
								eventTitle = 'Application Viewed';
								eventDescription = `Your application for ${app.jobTitle} was viewed by the hiring team at ${app.companyName}`;
								break;
							case 'ACCEPTED':
								eventTitle = 'Offer Received';
								eventDescription = `${app.companyName} extended an offer for the ${app.jobTitle} position`;
								break;
							case 'REJECTED':
								eventTitle = 'Application Declined';
								eventDescription = `Your application for ${app.jobTitle} at ${app.companyName} was not selected for this role`;
								break;
							default:
								eventTitle = 'Application Updated';
								eventDescription = `Your application status for ${app.jobTitle} at ${app.companyName} has been updated`;
						}
						return {
							time: formatRelativeTime(app.jobApplication.createdAt || new Date().toISOString()),
							title: eventTitle,
							description: eventDescription
						};
					})
					.slice(0, 3);
			}

			return apps;
		} catch (err) {
			console.error('Error fetching applications:', err);
			return [];
		}
	}

	// When user clicks a status filter button
	async function handleFilter(status: string) {
		filterStatus = status;
		applications = await fetchApplications(status);
	}

	function toggleDropdown(appId: string) {
		openDropdown = openDropdown === appId ? null : appId;
	}

	function viewJob(jobId: string) {
		goto(`/app/jobs/${jobId}`);
		openDropdown = null;
	}

	function initiateCancel(appId: string) {
		const app = applications.find((a) => a.id === appId);
		if (!app) {
			toast.error('Application not found');
			return;
		}
		applicationToCancel = app;
		showCancelModal = true;
		openDropdown = null;
	}

	async function confirmCancelApplication() {
		if (!applicationToCancel?.jobApplication?.id) {
			toast.error('Application not found');
			showCancelModal = false;
			return;
		}

		isCancelling = true;

		try {
			const res = await fetch(`/apply/${applicationToCancel.jobApplication.id}`, {
				method: 'DELETE',
				credentials: 'include'
			});

			if (!res.ok) {
				throw new Error(`Failed to cancel application: ${res.status}`);
			}

			toast.success('Application cancelled successfully');

			// Remove from the local list
			applications = applications.filter((a) => a.id !== applicationToCancel.id);

			showCancelModal = false;
			applicationToCancel = null;
		} catch (err) {
			console.error('Error cancelling application:', err);
			toast.error('Failed to cancel application. Please try again.');
		} finally {
			isCancelling = false;
		}
	}

	onMount(async () => {
		applications = await fetchApplications('All');
	});
</script>

<div class="">
	<div class="mb-8">
		<h1 class="text-2xl font-semibold text-gray-900">My Applications</h1>
		<p class="mt-2 text-gray-600">Track and manage all your job applications in one place</p>
	</div>
	<div class="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<!-- Total Applications -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Total Applications</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.totalApplications}</p>
			<p
				class={`mt-1 text-sm ${stats.trend.totalApplications >= 0 ? 'text-green-600' : 'text-red-600'}`}
			>
				{stats.trend.totalApplications >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.totalApplications)} from last month
			</p>
		</div>

		<!-- In Review -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">In Review</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.inReview}</p>
			<p class={`mt-1 text-sm ${stats.trend.inReview >= 0 ? 'text-green-600' : 'text-red-600'}`}>
				{stats.trend.inReview >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.inReview)} new this week
			</p>
		</div>

		<!-- Offer Received -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Offer Received</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.offerReceived}</p>
			<p
				class={`mt-1 text-sm ${stats.trend.offerReceived >= 0 ? 'text-green-600' : 'text-red-600'}`}
			>
				{stats.trend.offerReceived >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.offerReceived)} new this week
			</p>
		</div>

		<!-- Response Rate -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Response Rate</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.responseRate}%</p>
			<p
				class={`mt-1 text-sm ${stats.trend.responseRate >= 0 ? 'text-green-600' : 'text-red-600'}`}
			>
				{stats.trend.responseRate >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.responseRate)}% from last month
			</p>
		</div>
	</div>
	<div class="mb-6 flex items-center justify-between">
		<div class="flex gap-2">
			{#each ['All', 'Pending', 'In Review', 'Accepted', 'Rejected'] as status, i (i)}
				<button
					class={`rounded-full px-4 py-1.5 text-sm font-medium transition-colors hover:cursor-pointer ${
						filterStatus === status ? 'bg-green-600 text-white' : 'text-gray-700 hover:bg-gray-200'
					}`}
					onclick={() => handleFilter(status)}
				>
					{status}
				</button>
			{/each}
		</div>
		<!-- <button
			class="flex items-center gap-2 rounded-lg bg-white border border-gray-200 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors hover:cursor-pointer"
		>
			<Download class="h-4 w-4" />
			Export
		</button> -->
	</div>

	<div class="mb-8 overflow-hidden rounded-xl border border-gray-200 bg-white">
		<div class="border-b border-gray-200 bg-slate-50 px-6 py-4">
			<div class="flex items-center justify-between">
				<h2 class="text-md font-medium text-gray-900">Recent Applications</h2>
				<div class="flex items-center gap-3 text-sm text-gray-600">
					<span
						>Showing {showAllApplications ? applications.length : Math.min(6, applications.length)} of
						{applications.length}</span
					>

					<button
						class="font-medium text-gray-600 hover:cursor-pointer hover:text-gray-700"
						disabled
						title="Coming soon"
					>
						View all
					</button>
				</div>
			</div>
		</div>

		{#if applications.length === 0}
			<div class="p-8 text-center text-gray-500">
				No applications found for "{filterStatus}"
			</div>
		{:else}
			{#each showAllApplications ? applications : applications.slice(0, 6) as app, i (app.id)}
				<div
					class={`flex items-center justify-between px-6 py-4 ${i === (showAllApplications ? applications : applications.slice(0, 6)).length - 1 ? '' : 'border-b border-gray-100'}`}
				>
					<!-- Left part block -->
					<div class="flex items-center gap-4">
						<img
							src={app.companyLogo}
							alt="Company Logo"
							class="h-12 w-12 rounded-lg object-cover"
						/>
						<div class="flex flex-col">
							<span class="font-medium text-gray-900">{app.jobTitle}</span>
							<div class="mt-1 flex items-center gap-2 text-sm text-gray-500">
								<span>{app.companyName}</span>
								<span class="text-gray-300">•</span>
								<div class="flex items-center gap-1">
									<MapPin class="h-3 w-3" />
									<span>{app.location}</span>
								</div>
							</div>
						</div>
					</div>

					<!-- Right part block -->
					<div class="flex items-center gap-6">
						<div class="text-right">
							<span class="text-xs uppercase tracking-wide text-gray-500">Applied</span>
							<div class="text-sm font-medium text-gray-600">
								{formatRelativeTime(new Date(Date.now() - app.daysAgo * 24 * 60 * 60 * 1000))}
							</div>
						</div>

						<div
							class={`inline-flex items-center rounded-full border px-3 py-1.5 text-xs font-medium ${
								app.status === 'IN REVIEW'
									? 'border-amber-200 bg-amber-50 text-amber-700'
									: app.status === 'ACCEPTED'
										? 'border-emerald-200 bg-emerald-50 text-emerald-700'
										: app.status === 'PENDING'
											? 'border-blue-200 bg-blue-50 text-blue-700'
											: app.status === 'REJECTED'
												? 'border-red-200 bg-red-50 text-red-700'
												: 'border-gray-200 bg-gray-50 text-gray-700'
							}`}
						>
							{app.status}
						</div>

						<div class="relative">
							<button
								bind:this={dropdownTriggers[app.id]}
								onclick={() => toggleDropdown(app.id)}
								class="rounded-lg border border-gray-200 p-2 transition-colors hover:cursor-pointer hover:bg-gray-50"
							>
								<EllipsisVertical class="h-4 w-4 text-gray-600" />
							</button>

							<Dropdown
								isOpen={openDropdown === app.id}
								onClose={() => (openDropdown = null)}
								triggerElement={dropdownTriggers[app.id]}
								items={[
									{
										label: 'View Job',
										action: () => viewJob(app.jobId)
									},
									{
										label: 'Cancel',
										action: () => initiateCancel(app.id),
										variant: 'danger'
									}
								]}
							/>
						</div>
					</div>
				</div>
			{/each}
		{/if}
	</div>

	<div class="mt-8">
		<h2 class="mb-4 text-lg font-medium text-gray-900">Recent Activity</h2>
		<div class="relative">
			<div
				class="absolute left-4 w-px border-l border-gray-300 bg-gray-300"
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
							class="absolute left-4 top-1/2 h-5 w-5 -translate-x-1/2 -translate-y-1/2 rounded-full border-4 border-green-600 bg-white"
						></span>

						<div class="ml-10 rounded-xl border border-gray-200 bg-white p-4">
							<p class="mb-1 text-xs uppercase tracking-wide text-gray-500">{activity.time}</p>
							<p class="text-sm font-medium text-gray-800">{activity.title}</p>
							<p class="text-sm text-gray-500">{activity.description}</p>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>

<!-- Cancel Application Modal -->
<Modal bind:isOpen={showCancelModal} size="sm" closeOnBackdrop={!isCancelling}>
	<div class="p-6">
		<h3 class="mb-4 text-lg font-medium text-gray-900">Cancel Application</h3>
		{#if applicationToCancel}
			<p class="mb-6 text-sm text-gray-600">
				Are you sure you want to cancel your application for <span class="font-medium"
					>{applicationToCancel.jobTitle}</span
				>
				at <span class="font-medium">{applicationToCancel.companyName}</span>?
			</p>
		{/if}
		<div class="flex items-center justify-end gap-3">
			<button
				onclick={() => (showCancelModal = false)}
				disabled={isCancelling}
				class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:cursor-pointer hover:bg-gray-50 disabled:opacity-50"
			>
				Keep
			</button>
			<button
				onclick={confirmCancelApplication}
				disabled={isCancelling}
				class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:cursor-pointer hover:bg-red-700 disabled:opacity-50"
			>
				{isCancelling ? 'Cancelling...' : 'Cancel'}
			</button>
		</div>
	</div>
</Modal>
