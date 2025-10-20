<script>
	import { Search, ArrowUpDown, MapPin } from 'lucide-svelte';
	import ApplyModal from '$lib/components/job/ApplyModal.svelte';
	import JobDetailCard from '$lib/components/job/JobDetailCard.svelte';
	import { onMount } from 'svelte';
	import { jobSearchStore } from '$lib/stores/jobSearch';
	import { get } from 'svelte/store';
	import { isAuthenticated, getUserInfo } from '$lib/utils/auth';
	import {
		fetchCompany,
		fetchCompanyNameLogo,
		DEFAULT_COMPANY_LOGO,
		DEFAULT_COMPANY_NAME
	} from '$lib/utils/fetcher';
	import { formatDateShort } from '$lib/utils/datetime';

	let jobs = $state([]);
	let filteredJobs = $state([]);
	let selectedJob = $state(null);
	let searchQuery = $state('');
	let currentPage = $state(1);
	const pageSize = 4;

	let userInfo = $state(null);
	let companyInfo = $state(null);
	let isLoggedIn = $state(false);
	let appliedJobs = $state(new Set());

	let showApplyModal = $state(false);
	let jobToApply = $state(null);
	let isBookmarked = $state(false);
	let bookmarkedJobs = $state(new Set());

	const totalPages = $derived(Math.ceil(filteredJobs.length / pageSize));
	const paginatedJobs = $derived(filteredJobs.slice((currentPage - 1) * pageSize, currentPage * pageSize));
	
	$effect(() => {
		if (selectedJob?.companyID) {
			fetchCompanyInfo(selectedJob.companyID);
		}
		checkBookmarkStatus();
	});

	let activeFilters = $state({
		type: null,
		posted: null,
		arrangement: null
	});
	let sortBy = $state(null);

	const typeCycle = ['Full-time', 'Part-time', 'Contract', 'Casual'];
	const arrangementCycle = ['On-site', 'Remote', 'Hybrid'];
	const sortOptions = [null, 'dateAsc', 'dateDesc', 'title'];
	const sortLabels = {
		null: 'Sort',
		dateAsc: 'Oldest',
		dateDesc: 'Newest',
		title: 'Title A-Z'
	};

	async function fetchAppliedJobs() {
		if (!isLoggedIn || !userInfo?.userID) return;

		try {
			const res = await fetch(`/apply/query?applicantID=${userInfo.userID}`, {
				credentials: 'include'
			});
			if (!res.ok) throw new Error('Failed to fetch applied jobs');

			const data = await res.json();
			appliedJobs = new Set(data.map((app) => app.jobApplication.jobID));
		} catch (err) {
			console.error('Error fetching applied jobs:', err);
			appliedJobs = new Set();
		}
	}

	async function fetchCompanyInfo(companyID) {
		try {
			const raw = await fetchCompany(companyID);
			const infoArray = raw.userInfo || [];
			const info = Object.fromEntries(infoArray.map((item) => [item.Key, item.Value]));

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

	async function fetchJobs(query = '', filters = {}, sort = null) {
		try {
			const params = new URLSearchParams();
			if (query) params.set('title', query);
			if (filters.type) params.set('workType', filters.type);
			if (filters.posted) params.set('postOpenDate', filters.posted);
			if (filters.arrangement) params.set('workArrangement', filters.arrangement);
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

			const jobPromises = data.map(async (job) => {
				// fetch company name and logo.
				let [companyName, companyLogo] = await fetchCompanyNameLogo(job.companyID || '');

				return {
					id: job.id,
					title: job.title,
					company: companyName,
					companyID: job.companyID,
					location: job.location || 'N/A',
					type: job.workType || 'Full-time',
					tags: job.requiredSkills
						? job.requiredSkills.split(',').map((skill) => skill.trim())
						: [],
					posted: job.postOpenDate ? formatDateShort(job.postOpenDate) : 'Unknown',
					closeDate: job.applicationDeadline
						? formatDateShort(job.applicationDeadline)
						: 'Unknown',
					description: job.jobDescription || 'No description provided.',
					logo: companyLogo
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
		}
	}

	function applyJob(job) {
		jobToApply = job;
		showApplyModal = true;
	}

	function checkBookmarkStatus() {
		if (isAuthenticated() && userInfo?.userID && selectedJob?.id) {
			isBookmarked = bookmarkedJobs.has(selectedJob.id);
		}
	}

	function toggleBookmark() {
		if (!selectedJob?.id) return;

		if (isBookmarked) {
			bookmarkedJobs.delete(selectedJob.id);
			isBookmarked = false;
		} else {
			bookmarkedJobs.add(selectedJob.id);
			isBookmarked = true;
		}
		
		// TODO: Add backend API call when ready
		// await fetch('/bookmarks', { method: isBookmarked ? 'POST' : 'DELETE', ... });
	}


	function toggleCycle(field) {
		if (field === 'type') {
			let idx = typeCycle.indexOf(activeFilters.type);
			activeFilters.type =
				idx === -1 ? typeCycle[0] : idx + 1 < typeCycle.length ? typeCycle[idx + 1] : null;
		} else if (field === 'arrangement') {
			let idx = arrangementCycle.indexOf(activeFilters.arrangement);
			activeFilters.arrangement =
				idx === -1
					? arrangementCycle[0]
					: idx + 1 < arrangementCycle.length
						? arrangementCycle[idx + 1]
						: null;
		}
		activeFilters = { ...activeFilters };
		fetchJobs(searchQuery, activeFilters);
	}

	function toggleFilter(type, value) {
		activeFilters[type] = activeFilters[type] === value ? null : value;
		activeFilters = { ...activeFilters };
		fetchJobs(searchQuery, activeFilters, sortBy);
	}

	function toggleSort() {
		let idx = sortOptions.indexOf(sortBy);
		sortBy = idx === -1 || idx + 1 === sortOptions.length ? sortOptions[0] : sortOptions[idx + 1];
		fetchJobs(searchQuery, activeFilters, sortBy);
	}

	function onSearchInput() {
		fetchJobs(searchQuery, activeFilters);
	}

	onMount(() => {
		isLoggedIn = isAuthenticated();
		if (isLoggedIn) {
			userInfo = getUserInfo();
			fetchAppliedJobs();
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

<div class="flex h-screen">
	<main class="w-full min-w-0 flex-1 space-y-6 p-0.5">
		<!-- Search and Filters -->
		<div class="flex items-center gap-3">
			<div class="flex flex-1 items-center rounded-lg bg-white px-3 py-1 shadow">
				<Search class="h-5 w-5 text-gray-500" />
				<input
					type="text"
					placeholder="Search jobs..."
					class="ml-2 flex-1 border-none outline-none"
					bind:value={searchQuery}
					oninput={onSearchInput}
				/>
			</div>

			<div class="flex gap-2">
				<button
					class={`rounded-full px-3 py-1 text-sm ${activeFilters.type && typeCycle.includes(activeFilters.type) ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
					onclick={() => toggleCycle('type')}
				>
					{activeFilters.type || 'Work Type'}
				</button>

				<button
					class={`rounded-full px-3 py-1 text-sm ${activeFilters.posted === '1d' ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
					onclick={() => toggleFilter('posted', '1d')}
				>
					1 day ago
				</button>

				<button
					class={`rounded-full px-3 py-1 text-sm ${activeFilters.posted === '6w' ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
					onclick={() => toggleFilter('posted', '6w')}
				>
					6 weeks
				</button>

				<button
					class={`rounded-full px-3 py-1 text-sm ${activeFilters.arrangement && arrangementCycle.includes(activeFilters.arrangement) ? 'bg-green-600 text-white' : 'bg-gray-200'}`}
					onclick={() => toggleCycle('arrangement')}
				>
					{activeFilters.arrangement || 'Arrangement'}
				</button>

				<button
					class={`flex items-center rounded-full px-3 py-1 text-sm ${sortBy ? 'bg-green-600 text-white' : 'bg-gray-200 text-black'}`}
					onclick={toggleSort}
				>
					<ArrowUpDown class="mr-1 h-4 w-4" />
					{sortBy ? sortLabels[sortBy] : 'Sort'}
				</button>
			</div>
		</div>

		<div class="grid w-full grid-cols-3 gap-6 h-[calc(100vh-200px)]">
			<!-- Jobs List -->
			<section class="col-span-1 flex flex-col space-y-4 h-full">
				<div class="w-full flex-1 space-y-3 overflow-y-auto p-2">
					{#if paginatedJobs.length === 0}
						<div class="mt-4 text-center text-gray-500">No jobs match your search or filters.</div>
					{:else}
						{#each paginatedJobs as job (job.id)}
							<button
								onclick={() => (selectedJob = job)}
								class={`flex w-full cursor-pointer flex-col rounded-lg bg-white shadow ring-offset-2 hover:ring-2 hover:ring-green-500 ${selectedJob?.id === job.id ? 'ring-2 ring-green-600' : ''}`}
							>
								<div class="flex items-start gap-3 p-2">
									<img
										src={job.logo}
										alt={job.company}
										class="mt-2 h-12 w-12 flex-shrink-0 rounded-full object-cover"
									/>
									<div class="flex flex-1 flex-col text-left">
										<div class="font-semibold">{job.title}</div>
										<div class="text-sm text-gray-600">{job.company}</div>
										<div class="mt-1 flex items-center gap-2 text-sm text-black">
											<MapPin class="mt-0.5 h-4 w-4" />
											{job.location}
										</div>
									</div>
								</div>

								<div class="mt-2 flex flex-wrap gap-2 p-2">
									{#each job.tags as tag, i (i)}
										<span class="rounded-full bg-gray-100 px-2 py-1 text-sm">{tag}</span>
									{/each}
								</div>

								<div class="mx-3 mt-1 pb-2 text-left text-xs text-gray-500">{job.posted}</div>
							</button>
						{/each}
					{/if}
				</div>

				<!-- Pagination controls -->
				{#if totalPages > 1}
					<div class="my-4 flex justify-center gap-2">
						<button
							class="rounded bg-gray-200 px-3 py-1 disabled:opacity-50"
							onclick={() => (currentPage = Math.max(1, currentPage - 1))}
							disabled={currentPage === 1}
						>
							Prev
						</button>

						<span class="px-3 py-1 text-sm">
							Page {currentPage} of {totalPages}
						</span>

						<button
							class="rounded bg-gray-200 px-3 py-1 disabled:opacity-50"
							onclick={() => (currentPage = Math.min(totalPages, currentPage + 1))}
							disabled={currentPage === totalPages}
						>
							Next
						</button>
					</div>
				{/if}
			</section>

			<!-- Job Detail -->
			<section
				class="col-span-2 flex flex-col rounded-lg bg-white border border-gray-200 h-full overflow-hidden"
			>
				{#if selectedJob}
					<div class="h-full overflow-y-auto">
						<JobDetailCard 
							job={selectedJob}
							{companyInfo}
							onApply={applyJob}
							onBookmark={toggleBookmark}
							{isBookmarked}
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
	<ApplyModal 
		bind:isOpen={showApplyModal}
		job={jobToApply}
	/>
{/if}
