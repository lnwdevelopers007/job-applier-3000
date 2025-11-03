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
		fetchUser,
		fetchCompanyNameLogo,
		DEFAULT_COMPANY_LOGO,
		DEFAULT_COMPANY_NAME
	} from '$lib/utils/fetcher';
	import { formatDateShort } from '$lib/utils/datetime';
	import { bookmarkService } from '$lib/services/bookmarkService';
	import type { JobUI, UserInfo, JobCompanyInfo } from '$lib/types';

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
	let appliedJobs = $state(new Set<string>());

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

	async function fetchAppliedJobs() {
		if (!isLoggedIn || !userInfo?.userID) return;

		try {
			const res = await fetch(`/apply/query?applicantID=${userInfo.userID}`, {
				credentials: 'include'
			});
			if (!res.ok) throw new Error('Failed to fetch applied jobs');

			const data = await res.json();
			appliedJobs = new Set(data.map((app: any) => app.jobApplication.jobID));
		} catch (err) {
			console.error('Error fetching applied jobs:', err);
			appliedJobs = new Set();
		}
	}

	async function fetchCompanyInfo(companyID: string) {
		try {
			const raw = await fetchUser(companyID);
			const infoArray = raw.userInfo || [];
			const info = Object.fromEntries(infoArray.map((item: any) => [item.Key, item.Value]));

			companyInfo = {
				name: info.name || raw.name || DEFAULT_COMPANY_NAME,
				logo: info.logo || raw.avatarURL || DEFAULT_COMPANY_LOGO,
				location: info.headquarters || 'N/A',
				website: info.website || '',
				aboutUs: info.aboutUs || '',
				industry: info.industry || '',
				size: info.size || '',
				foundedYear: info.foundedYear || '',
				linkedIn: info.linkedIn || ''
			};
		} catch (err) {
			console.error('Error fetching company info:', err);
			companyInfo = null;
		}
	}

	async function fetchJobs(query = '', filters: typeof activeFilters = activeFilters, sort = '') {
		try {
			isLoading = true;
			const params = new URLSearchParams();
			if (query) params.set('title', query);

			if (filters.workType) {
				params.set('workType', filters.workType);
			}

			if (filters.postTime) {
				params.set('postOpenDate', filters.postTime);
			}

			if (filters.arrangement) {
				params.set('workArrangement', filters.arrangement);
			}

			if (sort) params.set('sort', sort);
			const res = await fetch(`/jobs/query?${params.toString()}`);

			if (res.status === 404) {
				jobs = [];
				filteredJobs = [];
				selectedJob = null;
				return;
			}

			if (!res.ok) throw new Error(`Failed to load jobs: ${res.status}`);
			const data = await res.json();

			const jobPromises = data.map(async (job: any) => {
				// fetch company name and logo.
				let [companyName, companyLogo] = await fetchCompanyNameLogo(job.companyID || '');

				return {
					id: job.id,
					title: job.title,
					company: companyName,
					companyID: job.companyID,
					location: job.location || 'N/A',
					workType: job.workType || 'Full-time',
					workArrangement: job.workArrangement || 'On-site',
					tags: job.requiredSkills
						? job.requiredSkills.split(',').map((skill: string) => skill.trim())
						: [],
					posted: job.postOpenDate ? formatDateShort(job.postOpenDate) : 'Unknown',
					closeDate: job.applicationDeadline ? formatDateShort(job.applicationDeadline) : 'Unknown',
					description: job.jobDescription || 'No description provided.',
					logo: companyLogo,
					salary:
						job.minSalary && job.maxSalary
							? `${job.currency || 'THB'} ${job.minSalary.toLocaleString()}-${job.maxSalary.toLocaleString()}`
							: undefined
				};
			});

			jobs = await Promise.all(jobPromises);
			filteredJobs = jobs;
			selectedJob = jobs[0] || null;
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
		isLoggedIn = isAuthenticated();
		if (isLoggedIn) {
			userInfo = getUserInfo();
			fetchAppliedJobs();
			// Initialize bookmark service with user ID
			if (userInfo?.userID) {
				bookmarkService.initializeBookmarks(userInfo.userID);
			}
		}

		const unsubscribe = jobSearchStore.subscribe((state) => {
			if (state.shouldFetch) {
				searchQuery = state.query;
				fetchJobs(state.query, activeFilters).then(() => {
					jobSearchStore.clearFetchFlag();
				});
			}
		});

		const currentState = get(jobSearchStore);
		if (!currentState.shouldFetch) {
			fetchJobs();
		}

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
			<!-- Jobs List -->
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
							isApplied={appliedJobs.has(selectedJob.id)}
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
