<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import { JobService } from '$lib/services/jobService';
	import { Search } from 'lucide-svelte';
	import DeleteModal from '$lib/components/ui/DeleteModal.svelte';
	import DataTable from './tables/DataTable.svelte';
	import { createJobColumns } from './tables/columns';
	import type { JobDisplay, JobFilters } from '$lib/types';

	interface Props {
		title?: string;
		description?: string;
		companyID?: string;
		showCompanyColumn?: boolean;
		editUrlPrefix?: string;
	}

	let {
		title = 'Job Management',
		description = 'Manage and monitor all job postings across your platform',
		companyID,
		showCompanyColumn = true,
		editUrlPrefix = '/company/edit'
	}: Props = $props();

	let jobs = $state<JobDisplay[]>([]);
	let loading = $state(true);
	let searchQuery = $state('');
	let statusFilter = $state('all');

	let showDeleteModal = $state(false);
	let jobToDelete = $state<JobDisplay | null>(null);
	let deleting = $state(false);

	const statusOptions = [
		{ value: 'all', label: 'All Statuses' },
		{ value: 'active', label: 'Active' },
		{ value: 'closed', label: 'Closed' },
	];

	function handleEdit(job: JobDisplay) {
		goto(`${editUrlPrefix}/${job.id}`);
	}

	function handleView(job: JobDisplay) {
		goto(`/app/jobs/${job.id}`);
	}

	function handleDelete(job: JobDisplay) {
		jobToDelete = job;
		showDeleteModal = true;
	}

	async function confirmDelete(reason: string) {
		if (!jobToDelete) return;
		deleting = true;
		try {
			await JobService.deleteJobWithReason(jobToDelete.id, reason);
			jobs = jobs.filter((j) => j.id !== jobToDelete?.id);
			showDeleteModal = false;
			jobToDelete = null;
		} catch (err) {
			console.error('Error deleting job:', err);
			throw err;
		} finally {
			deleting = false;
		}
	}

	const filteredJobs = $derived(
		jobs.filter((job) => {
			const matchesSearch =
				!searchQuery ||
				job.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
				(job.company && job.company.toLowerCase().includes(searchQuery.toLowerCase()));

			const matchesStatus = statusFilter === 'all' || 
				(job.status && job.status.toLowerCase() === statusFilter);

			return matchesSearch && matchesStatus;
		})
	);

	const columns = $derived(
		createJobColumns({
			showCompanyColumn,
			onView: handleView,
			onEdit: handleEdit,
			onDelete: handleDelete
		})
	);

	$effect(() => {
		if (authStore.isAuthenticated && authStore.user) {
			loadAllJobs();
		}
	});

	async function loadAllJobs() {
		try {
			loading = true;
			const queryParams: JobFilters | undefined = companyID ? { companyID } : undefined;
			jobs = await JobService.loadJobsForDisplay(queryParams);
		} catch (err) {
			console.error('Error loading jobs:', err);
			jobs = [];
		} finally {
			loading = false;
		}
	}

	function handleSearch(event: Event) {
		searchQuery = (event.target as HTMLInputElement).value;
	}

	function handleStatusFilter(event: Event) {
		statusFilter = (event.target as HTMLSelectElement).value;
	}
</script>

<div>
	<div>
		<div class="mb-8">
			<h1 class="text-2xl font-semibold text-gray-900 mb-1">{title}</h1>
			<p class="text-base text-gray-600 mb-6">{description}</p>
		</div>

		<div class="mb-6">
			<div class="flex gap-3">
        <div class="relative">
          <Search
            class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5"
          />
          <input
            type="text"
            placeholder="Search jobs{showCompanyColumn ? ' or companies' : ''}..."
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

		<DataTable data={filteredJobs} {columns} pageSize={15} {loading} />
	</div>
</div>

<DeleteModal
	bind:isOpen={showDeleteModal}
	onClose={() => {
		showDeleteModal = false;
		jobToDelete = null;
	}}
	onConfirm={confirmDelete}
	title="Delete Job Posting"
	itemName={jobToDelete?.title || ''}
	description="You're about to delete this job posting. This action cannot be undone and will remove the job posting along with all associated applications."
	reasonPlaceholder="Please provide a detailed reason for deleting this job posting..."
	confirmButtonText="Delete Job"
	isDeleting={deleting}
/>