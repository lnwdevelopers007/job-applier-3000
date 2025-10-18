<script>
	import { Search, ArrowUpDown, MapPin, Users, Building } from 'lucide-svelte';
	import SafeHTML from '$lib/utils/SafeHTML.svelte';
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

	let jobs = [];
	let filteredJobs = [];
	let selectedJob = null;
	let searchQuery = '';
	let currentPage = 1;
	const pageSize = 4;

	let userInfo = null;
	let companyInfo = null;
	let isLoggedIn = false;
	let appliedJobs = new Set();

	let isApplying = false;

	$: totalPages = Math.ceil(filteredJobs.length / pageSize);
	$: paginatedJobs = filteredJobs.slice((currentPage - 1) * pageSize, currentPage * pageSize);
	$: if (selectedJob?.companyID) {
		fetchCompanyInfo(selectedJob.companyID);
	}

	let activeFilters = {
		type: null,
		posted: null,
		arrangement: null
	};
	let sortBy = null;

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
					posted: job.postOpenDate ? new Date(job.postOpenDate).toLocaleDateString() : 'Unknown',
					closeDate: job.applicationDeadline
						? new Date(job.applicationDeadline).toLocaleDateString()
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

  async function applyJob(job) {
    try {
      isLoggedIn = isAuthenticated();
      userInfo = getUserInfo();
      if (!isLoggedIn || !userInfo.userID) {
        alert("❌ You must be logged in to apply.");
        return;
      }

			if (appliedJobs.has(job.id)) {
				alert('⚠️ You have already applied to this job.');
				return;
			}

			isApplying = true;
			const payload = {
				applicantID: userInfo.userID,
				jobID: job.id,
				companyID: job.companyID,
				status: 'PENDING',
				createdAt: new Date().toISOString()
			};

			const res = await fetch('/apply', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(payload)
			});

			if (!res.ok) {
				if (res.status === 409) {
					alert('⚠️ You already applied to this job.');
				} else {
					throw new Error(`Failed to apply: ${res.status}`);
				}
				return;
			}

			appliedJobs.add(job.id);
			alert(`✅ Successfully applied to ${job.title}`);
		} catch (err) {
			console.error('Error applying to job:', err);
			alert('❌ Failed to apply. Please try again.');
		}
		isApplying = false;
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

		<div class="grid w-full grid-cols-2 gap-6">
			<!-- Jobs List -->
			<section class="col-span-1 flex flex-col space-y-4">
				<div class="w-full flex-1 space-y-3 overflow-y-auto p-2">
					{#if paginatedJobs.length === 0}
						<div class="mt-4 text-center text-gray-500">No jobs match your search or filters.</div>
					{:else}
						{#each paginatedJobs as job (job.id)}
							<button
								onclick={() => (selectedJob = job)}
								class={`flex w-full cursor-pointer flex-col rounded-lg bg-white shadow ring-offset-2 hover:ring-2 hover:ring-green-300 ${selectedJob?.id === job.id ? 'ring-2 ring-green-500' : ''}`}
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
				class="col-span-1 flex min-h-[400px] flex-col space-y-2 rounded-lg bg-white p-6 shadow"
			>
				{#if selectedJob}
					<!-- Job Info -->
					<div class="mb-2">
						<div class="flex items-start gap-3 px-2">
							<img
								src={selectedJob.logo}
								alt={selectedJob.company}
								class="mt-0.5 h-4 w-4 flex-shrink-0 rounded-full object-cover"
							/>
							<p class="text-sm font-semibold text-gray-600">{selectedJob.company}</p>
						</div>

						<h2 class="text-xl font-bold">{selectedJob.title}</h2>

						<div class="my-2 flex items-center gap-2 text-sm text-black">
							<MapPin class="mt-0.5 h-4 w-4" />
							{selectedJob.location}
						</div>

						<div class="my-2 flex gap-4 text-sm text-gray-500">
							<span>Open: {selectedJob.posted}</span>
							<span>Close: {selectedJob.closeDate}</span>
						</div>
					</div>

					<!-- Tags -->
					<div class="flex gap-2">
						{#each selectedJob.tags as tag, i (i)}
							<span class="rounded-full bg-gray-100 px-2 py-1 text-sm">{tag}</span>
						{/each}
					</div>

					<!-- Apply / Bookmark -->
					<div class="mb-4 flex gap-2">
						{#if new Date(selectedJob.closeDate) < new Date()}
							<button
								class="cursor-not-allowed rounded-lg bg-gray-400 px-4 py-2 text-sm text-white"
								disabled
							>
								Closed
							</button>
						{:else if new Date(selectedJob.posted) > new Date()}
							<button
								class="cursor-not-allowed rounded-lg bg-gray-400 px-4 py-2 text-sm text-white"
								disabled
							>
								Not Open Yet
							</button>
						{:else}
							<button
								class={`rounded-lg px-4 py-2 text-sm text-white transition-colors duration-200 ${
									isApplying ? 'cursor-not-allowed bg-gray-500' : 'bg-green-600'
								}`}
								onclick={() => applyJob(selectedJob)}
								disabled={isApplying}
							>
								{isApplying ? 'Applying…' : 'Apply'}
							</button>
						{/if}

						<button class="rounded-lg bg-yellow-400 px-4 py-2 text-sm">Bookmark</button>
					</div>

					<!-- Job Description -->
					<h3 class="text-lg font-semibold">Job Description</h3>
					<div class="space-y-4">
						<SafeHTML html={selectedJob.description} />
					</div>

					<h3 class="text-lg font-semibold">About the company</h3>
					{#if companyInfo}
						<div class="mt-2 space-y-2 rounded-lg bg-gray-50 p-4 shadow">
							<div class="flex items-start gap-4">
								<img
									src={companyInfo.logo ||
										'https://images.unsplash.com/photo-1534237710431-e2fc698436d0?fm=jpg&q=60&w=3000'}
									alt={companyInfo.name}
									class="h-16 w-16 flex-shrink-0 rounded-full object-cover"
								/>
								<div class="flex flex-1 flex-col space-y-1">
									<h3 class="text-lg font-semibold">{companyInfo.name}</h3>

									{#if companyInfo.location}
										<p class="flex gap-1 text-sm text-gray-600">
											<MapPin class="mt-0.5 h-4 w-4" />
											{companyInfo.location}
										</p>
									{/if}
								</div>
							</div>

							{#if companyInfo.size || companyInfo.industry}
								<div class="flex gap-6 text-sm text-gray-700">
									{#if companyInfo.size}
										<span class="flex items-center gap-1">
											<Users class="mt-0.5 h-4 w-4" />
											{companyInfo.size}
										</span>
									{/if}
									{#if companyInfo.industry}
										<span class="flex items-center gap-1">
											<Building class="mt-0.5 h-4 w-4" />
											{companyInfo.industry}
										</span>
									{/if}
								</div>
							{/if}

							{#if companyInfo.aboutUs}
								<p class="mt-1 text-sm text-gray-700">{companyInfo.aboutUs}</p>
							{/if}
						</div>
					{:else}
						<div class="mt-4 text-gray-500">Company information not available.</div>
					{/if}
				{:else}
					<div class="flex flex-1 items-center justify-center text-gray-500">No job selected.</div>
				{/if}
			</section>
		</div>
	</main>
</div>
