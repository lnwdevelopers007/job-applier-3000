<script lang="ts">
	import { Search, ChevronLeft, ChevronRight } from 'lucide-svelte';
	import ApplyModal from '$lib/components/job/ApplyModal.svelte';
	import JobDetailCard from '$lib/components/job/JobDetailCard.svelte';
	import JobCard from '$lib/components/job/JobCard.svelte';
	import FilterPill from '$lib/components/forms/FilterPill.svelte';
	import { onMount } from 'svelte';
	import { jobSearchStore } from '$lib/stores/jobSearch';
	import { get } from 'svelte/store';
	import { isAuthenticated, getUserInfo } from '$lib/utils/auth';
	import {
		DEFAULT_COMPANY_LOGO,
		DEFAULT_COMPANY_NAME
	} from '$lib/utils/fetcher';
	import { bookmarkService } from '$lib/services/bookmarkService';
	import { JobService } from '$lib/services/jobService';
	import { UserService } from '$lib/services/userService';
	import type { JobUI, UserInfo, JobCompanyInfo, Job, JobFilters, CompanyInfo } from '$lib/types';

	let jobs = $state<JobUI[]>([]);
	let filteredJobs = $state<JobUI[]>([]);
	let selectedJob = $state<JobUI | null>(null);
	let searchQuery = $state('');
	let currentPage = $state(1);
	let isLoading = $state(false);
	const pageSize = 6;

	let userInfo = $state<UserInfo | null>(null);
	let companyInfo = $state<JobCompanyInfo | null>(null);
	let isLoggedIn = $state(false);

	let showApplyModal = $state(false);
	let jobToApply = $state<JobUI | null>(null);
	let selectedJobBookmarked = $state(false);

	const totalPages = $derived(Math.ceil(filteredJobs.length / pageSize));
	const paginatedJobs = $derived(
		filteredJobs.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	$effect(() => {
		if (selectedJob?.companyID) {
			fetchCompanyInfo(selectedJob.companyID);
		}
	});

	// Subscribe to bookmark changes for selected job
	$effect(() => {
		const unsubscribe = bookmarkService.subscribe((bookmarks) => {
			if (selectedJob?.id) {
				selectedJobBookmarked = bookmarks.has(selectedJob.id);
			}
		});
		return unsubscribe;
	});


	let activeFilters = $state<{
		workType: string;
		postTime: string;
		arrangement: string;
	}>({
		workType: '',
		postTime: '',
		arrangement: ''
	});
	let sortBy = $state('');

	const workTypeOptions = [
		{ value: 'Full-time', label: 'Full-time' },
		{ value: 'Part-time', label: 'Part-time' },
		{ value: 'Contract', label: 'Contract' },
		{ value: 'Casual', label: 'Casual' }
	];

	const postTimeOptions = [
		{ value: '1d', label: 'Past 24 hours' },
		{ value: '6w', label: 'Past 6 weeks' }
	];

	const arrangementOptions = [
		{ value: 'On-site', label: 'On-site' },
		{ value: 'Remote', label: 'Remote' },
		{ value: 'Hybrid', label: 'Hybrid' }
	];

	const sortOptions = [
		{ value: 'dateDesc', label: 'Newest first' },
		{ value: 'dateAsc', label: 'Oldest first' },
		{ value: 'title', label: 'Title A-Z' }
	];


	async function fetchCompanyInfo(companyID: string) {
		try {
			const user = await UserService.getUserById(companyID);
			const company = user.userInfo as CompanyInfo;
			
			// Transform user data to company info format
			companyInfo = {
				name: company?.name || user.name || DEFAULT_COMPANY_NAME,
				logo: company?.logo || user.avatarURL || DEFAULT_COMPANY_LOGO,
				location: company?.headquarters || 'N/A',
				website: company?.website || '',
				aboutUs: company?.aboutUs || '',
				industry: company?.industry || '',
				size: company?.size || '',
				foundedYear: company?.foundedYear || '',
				linkedIn: company?.linkedIn || ''
			};
		} catch (err) {
			console.error('Error fetching company info:', err);
			companyInfo = null;
		}
	}

	async function fetchJobs(query = '', filters: typeof activeFilters = activeFilters, sort = '') {
		try {
			isLoading = true;
			
			// Build filters object for JobService
			const jobFilters: JobFilters = {};
			if (query) jobFilters.title = query;
			if (filters.workType) jobFilters.workType = filters.workType;
			if (filters.postTime) jobFilters.postOpenDate = filters.postTime;
			if (filters.arrangement) jobFilters.workArrangement = filters.arrangement;
			if (sort && (sort === 'dateDesc' || sort === 'dateAsc' || sort === 'title')) {
				jobFilters.sort = sort;
			}

			// Use JobService to fetch and transform jobs
			const rawJobs = await JobService.queryJobs(jobFilters);
			
			if (!rawJobs || rawJobs.length === 0) {
				jobs = [];
				filteredJobs = [];
				selectedJob = null;
				return;
			}

			// Transform jobs to UI format
			const transformedJobs = await Promise.all(
				rawJobs.map(async (job: Job): Promise<JobUI> => {
					const displayJob = await JobService.transformJobForDisplay(job);
					console.log('Debug fetchJobs - Raw job object:', job);
					console.log('Debug fetchJobs - Job ID being used:', job.id);
					
					return {
						id: job.id || '',
						title: job.title || '',
						company: displayJob.company || 'Unknown Company',
						companyID: job.companyID || '',
						location: job.location || 'N/A',
						workType: job.workType || 'Full-time',
						workArrangement: job.workArrangement || 'On-site',
						tags: job.requiredSkills
							? job.requiredSkills.split(',').map((skill: string) => skill.trim())
							: [],
						posted: displayJob.posted || 'Unknown',
						closeDate: displayJob.expires || 'Unknown',
						description: job.jobDescription || 'No description provided.',
						logo: displayJob.companyLogo || DEFAULT_COMPANY_LOGO,
						salary:
							job.minSalary && job.maxSalary
								? `${job.currency || 'THB'} ${job.minSalary.toLocaleString()}-${job.maxSalary.toLocaleString()}`
								: undefined
					};
				})
			);

			jobs = transformedJobs;
			filteredJobs = transformedJobs;
			selectedJob = transformedJobs[0] || null;
			currentPage = 1;
		} catch (err) {
			console.error('Error fetching jobs:', err);
			jobs = [];
			filteredJobs = [];
			selectedJob = null;
		} finally {
			isLoading = false;
		}
	}

	function applyJob(job: JobUI) {
		jobToApply = job;
		showApplyModal = true;
	}

	function refreshJobs() {
		fetchJobs(searchQuery, activeFilters, sortBy);
	}

	onMount(() => {
		async function initializeApp() {
			isLoggedIn = isAuthenticated();
			if (isLoggedIn) {
				userInfo = getUserInfo();
				// Initialize bookmark service with user ID
				if (userInfo?.userID) {
					bookmarkService.initializeBookmarks(userInfo.userID);
				}
			}

			const currentState = get(jobSearchStore);
			if (!currentState.shouldFetch) {
				// Only fetch jobs after applied jobs are loaded
				fetchJobs();
			}
		}

		// Start initialization
		initializeApp();

		const unsubscribe = jobSearchStore.subscribe((state) => {
			if (state.shouldFetch) {
				searchQuery = state.query;
				fetchJobs(state.query, activeFilters).then(() => {
					jobSearchStore.clearFetchFlag();
				});
			}
		});

		return () => {
			unsubscribe();
		};
	});
</script>

<div class="flex">
	<main class="w-full min-w-0 flex-1 space-y-6">
		<!-- Search and Filters -->
		<div
			class="-mt-5 pb-8 pt-8"
			style="margin-left: calc(-50vw + 50%); margin-right: calc(-50vw + 50%); padding-left: calc(50vw - 50%); padding-right: calc(50vw - 50%); background: radial-gradient(circle,rgba(245, 255, 252, 1) 0%, rgba(248, 255, 249, 1) 25%, rgba(243, 255, 245, 1) 50%, rgba(237, 254, 244, 1) 75%, rgba(232, 254, 240, 1) 100%);"
		>
			<div class="mx-auto max-w-7xl space-y-4">
				<!-- Search Bar -->
				<div class="relative mx-auto max-w-3xl">
					<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
						<Search class="h-5 w-5 text-gray-400" />
					</div>
					<input
						type="text"
						placeholder="Search for jobs..."
						class="w-full rounded-full border border-gray-200 bg-white py-3 pl-12 pr-4 text-gray-900 placeholder-gray-500 shadow-sm outline-none transition-all duration-200 focus:border-transparent focus:ring-1 focus:ring-gray-300"
						bind:value={searchQuery}
						oninput={refreshJobs}
					/>
				</div>

				<!-- Filter Pills -->
				<div class="flex flex-wrap items-center justify-center gap-3">
					<div class="flex items-center gap-2">
						<span class="text-sm font-medium text-gray-600">Filters:</span>

						<FilterPill
							label="Work Type"
							options={workTypeOptions}
							bind:selectedValue={activeFilters.workType}
							onSelectionChange={refreshJobs}
						/>

						<FilterPill
							label="Post Time"
							options={postTimeOptions}
							bind:selectedValue={activeFilters.postTime}
							onSelectionChange={refreshJobs}
						/>

						<FilterPill
							label="Work Arrangement"
							options={arrangementOptions}
							bind:selectedValue={activeFilters.arrangement}
							onSelectionChange={refreshJobs}
						/>

						<div class="h-4 w-px bg-gray-300"></div>

						<FilterPill
							label="Sort"
							options={sortOptions}
							bind:selectedValue={sortBy}
							onSelectionChange={refreshJobs}
							type="sort"
						/>
					</div>
				</div>
			</div>
		</div>

		<div class="grid w-full grid-cols-3 gap-6">
			<!-- Jobls List -->
			<section class="col-span-1 flex h-[calc(100vh-280px)] flex-col">
				<div class="min-h-0 flex-1 overflow-y-auto p-2">
					<div class="space-y-3">
						{#if isLoading}
							<!-- Show skeleton JobCards during loading -->
							{#each Array.from({ length: 8 }, (_, i) => i) as i (i)}
								<div class="transition-all duration-100 rounded-lg">
									<JobCard loading={true} />
								</div>
							{/each}
						{:else if paginatedJobs.length === 0}
							<div class="flex items-center justify-center h-full text-gray-500">
								No jobs match your search or filters.
							</div>
						{:else}
							{#each paginatedJobs as job (job.id)}
								<div
									class={`rounded-lg transition-all duration-100 ${
										selectedJob?.id === job.id
											? 'border-green-200 shadow-md ring-2 ring-green-600 ring-offset-2'
											: ''
									}`}
								>
									<JobCard {job} onclick={() => (selectedJob = job)} />
								</div>
							{/each}

							<!-- Pagination inside scrollable area -->
							{#if !isLoading && totalPages > 1}
								<div class="flex items-center justify-center gap-3 py-3">
									<button
										aria-label="Previous page"
										class="flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
										onclick={() => (currentPage = Math.max(1, currentPage - 1))}
										disabled={currentPage === 1}
									>
										<ChevronLeft class="h-4 w-4" />
									</button>

									<div class="flex items-center gap-2">
										<span class="text-sm text-gray-600">Page</span>
										<span class="text-sm font-medium text-gray-600">
											{currentPage}
										</span>
										<span class="text-sm text-gray-600">of</span>
										<span class="text-sm font-medium text-gray-600">
											{totalPages}
										</span>
									</div>

									<button
										aria-label="Next page"
										class="flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-2 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
										onclick={() => (currentPage = Math.min(totalPages, currentPage + 1))}
										disabled={currentPage === totalPages}
									>
										<ChevronRight class="h-4 w-4" />
									</button>
								</div>
							{/if}
						{/if}
					</div>
				</div>
			</section>

			<!-- Job Detail -->
			<section
				class="col-span-2 flex flex-col h-[calc(100vh-280px)] overflow-hidden"
			>
				{#if isLoading}
					<div class="flex flex-1 items-center justify-center text-gray-500">Loading job details...</div>
				{:else if selectedJob}
					<div class="h-full overflow-y-auto">
						<JobDetailCard
							job={selectedJob}
							{companyInfo}
							onApply={() => selectedJob && applyJob(selectedJob)}
							onBookmark={async () => {
								if (selectedJob?.id) {
									const user = getUserInfo();
									if (user?.userID) {
										await bookmarkService.toggleBookmark(selectedJob.id, user.userID);
									}
								}
							}}
							isBookmarked={selectedJobBookmarked}
						/>
					</div>
				{:else}
					<div class="flex flex-1 items-center justify-center text-gray-500">No job selected.</div>
				{/if}
			</section>
		</div>
	</main>
</div>

<!-- Apply Modal -->
{#if jobToApply}
	<ApplyModal bind:isOpen={showApplyModal} job={jobToApply} />
{/if}
