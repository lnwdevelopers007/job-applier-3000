<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { apiFetch } from '$lib/utils/api';
	import { getCompanyAnalytics } from '$lib/utils/companyStats';
	import { formatDateDMY } from '$lib/utils/datetime';
	import { Search } from 'lucide-svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import { createCompanyJobColumns, type CompanyJobDisplay } from '$lib/components/tables/columns/companyJobColumns';
	import DeleteModal from '$lib/components/ui/DeleteModal.svelte';

	let jobs = $state<CompanyJobDisplay[]>([]);
	let loading = $state(true);
	let searchQuery = $state('');
	let statusFilter = $state('all');

	// Modal states
	let showDeleteModal = $state(false);
	let jobToDelete = $state<CompanyJobDisplay | null>(null);
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

	const statusOptions = [
		{ value: 'all', label: 'All Status' },
		{ value: 'active', label: 'Active' },
		{ value: 'closed', label: 'Closed' }
	];

	function handleView(job: CompanyJobDisplay) {
		goto(`/app/jobs/${job.id}`);
	}

	function handleEdit(job: CompanyJobDisplay) {
		goto(`/company/edit/${job.id}`);
	}

	function handleDelete(job: CompanyJobDisplay) {
		jobToDelete = job;
		showDeleteModal = true;
	}

	function handleViewApplicants(job: CompanyJobDisplay) {
		// Navigate to applicants page with job filter
		goto(`/company/applicants?jobId=${job.id}`);
	}

	async function confirmDelete(reason: string) {
		if (!jobToDelete) return;
		deleting = true;
		try {
			const res = await apiFetch(`/jobs/${jobToDelete.id}`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ reason })
			});

			if (!res.ok) throw new Error(`Failed to delete job: ${res.status}`);

			jobs = jobs.filter((j) => j.id !== jobToDelete.id);
			showDeleteModal = false;
			jobToDelete = null;
		} catch (err) {
			console.error(err);
			throw err;
		} finally {
			deleting = false;
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
			loading = true;
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
						const applyRes = await apiFetch(`/apply/query?jobID=${jobID}`);
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
						posted: job.postOpenDate ? formatDateDMY(job.postOpenDate) : 'None',
						expires: job.applicationDeadline
							? formatDateDMY(job.applicationDeadline)
							: 'None'
					};
				})
			);

			jobs = jobWithCounts;
		} catch {
			jobs = [];
		} finally {
			loading = false;
		}
	}

	const columns = $derived(
		createCompanyJobColumns({
			onView: handleView,
			onEdit: handleEdit,
			onDelete: handleDelete,
			onViewApplicants: handleViewApplicants
		})
	);

	const filteredJobs = $derived(
		jobs.filter((job) => {
			const matchesSearch =
				!searchQuery ||
				job.title.toLowerCase().includes(searchQuery.toLowerCase());

			const matchesStatus =
				statusFilter === 'all' ||
				job.status.toLowerCase() === statusFilter;

			return matchesSearch && matchesStatus;
		})
	);

	function handleSearch(event: Event) {
		searchQuery = (event.target as HTMLInputElement).value;
	}

	function handleStatusFilter(event: Event) {
		statusFilter = (event.target as HTMLSelectElement).value;
	}
</script>

<div class="">
	<div class="mb-8">
		<h1 class="text-2xl font-semibold text-gray-900">Company Dashboard</h1>
		<p class="mt-2 text-gray-600">Welcome back! Here's what's happening with your recruitment.</p>
	</div>

	<!-- Summary cards -->
	<div class="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		<!-- Active Jobs -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Active Jobs</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.activeJobs}</p>
			<p class={`mt-1 text-sm ${stats.trend.activeJobs >= 0 ? 'text-green-600' : 'text-red-600'}`}>
				{stats.trend.activeJobs >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.activeJobs)} from last month
			</p>
		</div>

		<!-- Total Applicants -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Total Applicants</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.totalApplicants}</p>
			<p class={`mt-1 text-sm ${stats.trend.totalApplicants >= 0 ? 'text-green-600' : 'text-red-600'}`}>
				{stats.trend.totalApplicants >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.totalApplicants)} from last month
			</p>
		</div>

		<!-- Pending Review -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Pending Review</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.pendingReview}</p>
			<p class={`mt-1 text-sm ${stats.trend.pendingReview >= 0 ? 'text-green-600' : 'text-red-600'}`}>
				{stats.trend.pendingReview >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.pendingReview)} new today
			</p>
		</div>

		<!-- Offers Accepted -->
		<div class="rounded-xl border border-gray-200 bg-white p-4">
			<h2 class="mb-2 text-sm text-gray-500">Offers Accepted</h2>
			<p class="text-2xl font-medium text-gray-900">{stats.offersAccepted}</p>
			<p class={`mt-1 text-sm ${stats.trend.offersAccepted >= 0 ? 'text-green-600' : 'text-red-600'}`}>
				{stats.trend.offersAccepted >= 0 ? '↑' : '↓'}
				{Math.abs(stats.trend.offersAccepted)} from last month
			</p>
		</div>
	</div>

	<!-- Jobs Table -->
	<div class="mb-8">
		<div class="mb-6 flex items-center justify-between">
			<h2 class="text-lg font-medium text-gray-900">Posted Jobs</h2>
			<div class="flex items-center gap-3 text-sm text-gray-600">
				<span>Showing {filteredJobs.length} of {jobs.length} job{jobs.length === 1 ? '' : 's'}</span>
			</div>
		</div>

		<!-- Search and Filters -->
		<div class="mb-6">
			<div class="flex gap-3">
				<div class="relative">
					<Search
						class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5"
					/>
					<input
						type="text"
						placeholder="Search jobs by title..."
						class="w-full pl-10 pr-4 py-2.5 border border-gray-200 rounded-lg min-w-lg text-sm placeholder:text-gray-500 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400 transition-all"
						bind:value={searchQuery}
						oninput={handleSearch}
					/>
				</div>

				<select
					class="appearance-none px-4 pr-10 py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700 font-medium hover:bg-gray-100 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400 transition-all cursor-pointer"
					bind:value={statusFilter}
					onchange={handleStatusFilter}
				>
					{#each statusOptions as option (option.value)}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
			</div>
		</div>

		<DataTable data={filteredJobs} {columns} pageSize={10} {loading} />
	</div>
</div>

<DeleteModal
	bind:isOpen={showDeleteModal}
	onClose={() => {
		showDeleteModal = false;
		jobToDelete = null;
	}}
	onConfirm={confirmDelete}
	title="Delete Job"
	itemName={jobToDelete?.title || ''}
	description="You're about to delete this job posting. This action cannot be undone and will remove the job along with all associated applications."
	reasonPlaceholder="Please provide a detailed reason for deleting this job posting..."
	confirmButtonText="Delete Job"
	isDeleting={deleting}
/>
