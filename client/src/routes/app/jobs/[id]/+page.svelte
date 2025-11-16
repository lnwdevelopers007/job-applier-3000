<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { MapPin, Bookmark, Share2, Clock, Banknote } from 'lucide-svelte';
	import SafeHTML from '$lib/utils/SafeHTML.svelte';
	import SkillTag from '$lib/components/job/SkillTag.svelte';
	import ApplyButton from '$lib/components/job/ApplyButton.svelte';
	import CompanyCard from '$lib/components/job/CompanyCard.svelte';
	import JobCard from '$lib/components/job/JobCard.svelte';
	import JobDetailSkeleton from '$lib/components/job/JobDetailSkeleton.svelte';
	import ApplyModal from '$lib/components/job/ApplyModal.svelte';
	import FloatingJobHeader from '$lib/components/job/FloatingJobHeader.svelte';
	import {
		fetchJob,
		DEFAULT_COMPANY_LOGO
	} from '$lib/utils/fetcher';
	import { isAuthenticated, getUserInfo } from '$lib/utils/auth';
	import { formatDateDMY, formatRelativeTime, formatDateShort } from '$lib/utils/datetime';
	import { bookmarkService } from '$lib/services/bookmarkService';
	import { JobService } from '$lib/services/jobService';
	import { UserService } from '$lib/services/userService';

	let { data }: { data: { jobId: string } } = $props();

	let job = $state<any>(null);
	let company = $state<any>(null);
	let similarJobs = $state<any[]>([]);
	let otherCompanyJobs = $state<any[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let isBookmarked = $state(false);
	let userInfo: any = null;
	let showFloatingCard = $state(false);
	let applyButtonRef = $state<HTMLElement | undefined>();
	let showApplyModal = $state(false);
	let showAllSkills = $state(false);

	// Subscribe to bookmark changes
	$effect(() => {
		const unsubscribe = bookmarkService.subscribe((jobs) => {
			if (job?.id) {
				isBookmarked = jobs.has(job.id);
			}
		});

		return unsubscribe;
	});

	// Initialize bookmarks when job loads
	$effect(() => {
		if (job?.id) {
			const user = getUserInfo();
			if (user?.userID) {
				bookmarkService.initializeBookmarks(user.userID).then(() => {
					isBookmarked = bookmarkService.isBookmarked(job.id);
				});
			}
		}
	});


	function formatSalary(min: number | null, max: number | null, currency: string): string {
		if (!min && !max) return '';

		const formatNumber = (num: number) => {
			return new Intl.NumberFormat('en-US').format(Math.round(num));
		};

		if (min && max) {
			return `${formatNumber(min)} - ${formatNumber(max)} ${currency}/month`;
		} else if (min) {
			return `From ${formatNumber(min)} ${currency}/month`;
		} else if (max) {
			return `Up to ${formatNumber(max)} ${currency}/month`;
		}
		return '';
	}


	async function loadJobData() {
		try {
			loading = true;
			error = null;

			// Check authentication
			if (isAuthenticated()) {
				userInfo = getUserInfo();
			}

			// Use JobService to get the job data
			const jobData = await JobService.getJobById(data.jobId);
			const displayJob = await JobService.transformJobForDisplay(jobData);

			job = {
				id: jobData.id || '',
				title: jobData.title || '',
				company: displayJob.company || 'Unknown Company',
				companyId: jobData.companyID || '',
				location: jobData.location || 'N/A',
				workType: jobData.workType || 'onsite',
				workArrangement: jobData.workArrangement || 'full-time',
				salary: formatSalary(jobData.minSalary, jobData.maxSalary, jobData.currency || 'THB'),
				posted: jobData.postOpenDate ? formatRelativeTime(jobData.postOpenDate) : 'Unknown',
				postedDate: jobData.postOpenDate || null,
				closeDate: jobData.applicationDeadline
					? formatDateDMY(jobData.applicationDeadline)
					: 'Open until filled',
				closeDateRaw: jobData.applicationDeadline || null,
				description: jobData.jobDescription || 'No description provided.',
				skills: jobData.requiredSkills
					? jobData.requiredSkills.split(',').map((skill: string) => skill.trim())
					: [],
				logo: displayJob.companyLogo || DEFAULT_COMPANY_LOGO
			};

			if (jobData.companyID) {
				try {
					// Use UserService to get company details
					const user = await UserService.getUserById(jobData.companyID);
					const companyInfo = user.userInfo as any;

					company = {
						name: displayJob.company || 'Unknown Company',
						logo: displayJob.companyLogo || DEFAULT_COMPANY_LOGO,
						industry: companyInfo?.industry || 'Software Development',
						employees: companyInfo?.size || '1000+ employees',
						description: companyInfo?.aboutUs || 'Company information not available.'
					};

					await loadSimilarJobs();
					await loadOtherCompanyJobs(jobData.companyID);
				} catch (companyErr) {
					console.warn('Failed to fetch company info:', companyErr);
					// Set basic company info from display job
					company = {
						name: displayJob.company || 'Unknown Company',
						logo: displayJob.companyLogo || DEFAULT_COMPANY_LOGO,
						industry: 'Software Development',
						employees: '1000+ employees',
						description: 'Company information not available.'
					};
				}
			}
		} catch (err) {
			console.error('Error loading job:', err);
			error = 'Failed to load job details. Please try again.';
		} finally {
			loading = false;
		}
	}

	async function loadSimilarJobs() {
		try {
			// Use JobService to get all jobs
			const allJobs = await JobService.getAllJobs();
			if (!allJobs || allJobs.length === 0) return;

			// Filter out current job and calculate similarity scores
			const otherJobs = allJobs.filter((j: any) => j.id !== job?.id);
			const jobsWithScores = otherJobs.map((jobData: any) => {
				let score = 0;

				// Score based on title similarity
				if (job?.title && jobData.title) {
					const currentTitle = job.title.toLowerCase();
					const otherTitle = jobData.title.toLowerCase();
					const commonWords = ['the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for', 'of', 'with', 'by'];
					const currentKeywords = currentTitle.split(/\s+/).filter((word: string) => word.length > 2 && !commonWords.includes(word));
					const otherKeywords = otherTitle.split(/\s+/).filter((word: string) => word.length > 2 && !commonWords.includes(word));

					const matchingKeywords = currentKeywords.filter((keyword: string) =>
						otherKeywords.some((otherKeyword: string) => otherKeyword.includes(keyword) || keyword.includes(otherKeyword))
					);
					score += (matchingKeywords.length / Math.max(currentKeywords.length, 1)) * 3;
				}

				// Score based on skills similarity
				if (job?.skills && jobData.requiredSkills) {
					const currentSkills = job.skills.map((s: string) => s.toLowerCase().trim());
					const otherSkills = jobData.requiredSkills.split(',').map((s: string) => s.toLowerCase().trim());

					const matchingSkills = currentSkills.filter((skill: string) =>
						otherSkills.some((otherSkill: string) => otherSkill.includes(skill) || skill.includes(otherSkill))
					);
					score += (matchingSkills.length / Math.max(currentSkills.length, 1)) * 2;
				}

				// Minor bonus for same work type
				if (job?.workType && jobData.workType === job.workType) {
					score += 0.5;
				}

				return { ...jobData, similarityScore: score };
			});

			// Sort by similarity score and take top 3
			const filteredJobs = jobsWithScores
				.sort((a: any, b: any) => b.similarityScore - a.similarityScore)
				.slice(0, 3);

			// Transform jobs to display format using JobService
			similarJobs = await Promise.all(
				filteredJobs.map(async (jobData: any) => {
					try {
						const displayJob = await JobService.transformJobForDisplay(jobData);
						return {
							id: jobData.id,
							title: jobData.title,
							company: displayJob.company,
							location: jobData.location || 'N/A',
							logo: displayJob.companyLogo,
							workType: jobData.workType || 'FULL-TIME',
							workArrangement: jobData.workArrangement || 'ON-SITE',
							posted: jobData.postOpenDate ? formatDateShort(jobData.postOpenDate) : 'Unknown'
						};
					} catch (err) {
						console.warn(`Failed to transform similar job ${jobData.id}:`, err);
						return {
							id: jobData.id,
							title: jobData.title,
							company: 'Unknown Company',
							location: jobData.location || 'N/A',
							logo: DEFAULT_COMPANY_LOGO,
							workType: jobData.workType || 'FULL-TIME',
							workArrangement: jobData.workArrangement || 'ON-SITE',
							posted: jobData.postOpenDate ? formatDateShort(jobData.postOpenDate) : 'Unknown'
						};
					}
				})
			);
		} catch (err) {
			console.error('Error loading similar jobs:', err);
		}
	}

	async function loadOtherCompanyJobs(companyId: string) {
		try {
			// Use JobService to query jobs with company filter
			const companyJobs = await JobService.queryJobs({ companyID: companyId });
			if (!companyJobs || companyJobs.length === 0) return;

			const filteredJobs = companyJobs.filter((j: any) => j.id !== job?.id).slice(0, 2);

			otherCompanyJobs = await Promise.all(
				filteredJobs.map(async (jobData: any) => {
					try {
						const displayJob = await JobService.transformJobForDisplay(jobData);
						return {
							id: jobData.id,
							title: jobData.title,
							company: displayJob.company,
							location: jobData.location || 'N/A',
							logo: displayJob.companyLogo,
							workType: jobData.workType || 'FULL-TIME',
							workArrangement: jobData.workArrangement || 'ON-SITE',
							posted: jobData.postOpenDate ? formatDateShort(jobData.postOpenDate) : 'Unknown'
						};
					} catch (err) {
						console.warn(`Failed to transform other company job ${jobData.id}:`, err);
						return {
							id: jobData.id,
							title: jobData.title,
							company: 'Unknown Company',
							location: jobData.location || 'N/A',
							logo: DEFAULT_COMPANY_LOGO,
							workType: jobData.workType || 'FULL-TIME',
							workArrangement: jobData.workArrangement || 'ON-SITE',
							posted: jobData.postOpenDate ? formatDateShort(jobData.postOpenDate) : 'Unknown'
						};
					}
				})
			);
		} catch (err) {
			console.error('Error loading other company jobs:', err);
		}
	}

	async function toggleBookmark() {
		if (!job?.id) return;

		const user = getUserInfo();
		if (user?.userID) {
			const newState = await bookmarkService.toggleBookmark(job.id, user.userID);
			isBookmarked = newState;
		} else {
			console.warn('User must be logged in to bookmark jobs');
		}
	}


	function handleShare() {
		if (navigator.share) {
			navigator.share({
				title: job?.title || 'Job Opportunity',
				text: `Check out this job: ${job?.title} at ${job?.company}`,
				url: window.location.href
			});
		} else {
			navigator.clipboard.writeText(window.location.href);
			alert('Job link copied to clipboard!');
		}
	}

	function navigateToJob(jobId: string) {
		goto(`/app/jobs/${jobId}`);
	}

	function handleApplyClick() {
		showApplyModal = true;
	}

	// React to job ID changes
	$effect(() => {
		if (data.jobId) {
			loadJobData();
		}
	});

	onMount(() => {
		// Handle scroll detection for floating card
		const handleScroll = () => {
			if (applyButtonRef) {
				const rect = applyButtonRef.getBoundingClientRect();
				showFloatingCard = rect.bottom < 0;
			}
		};

		window.addEventListener('scroll', handleScroll);
		return () => {
			window.removeEventListener('scroll', handleScroll);
		};
	});
</script>

{#if loading}
	<JobDetailSkeleton />
{:else if error}
	<div class="flex min-h-[400px] items-center justify-center">
		<div class="text-center">
			<p class="mb-4 text-red-600">{error}</p>
			<button
				onclick={loadJobData}
				class="rounded-md bg-green-600 px-4 py-2 text-white transition-colors hover:bg-green-700"
			>
				Try Again
			</button>
		</div>
	</div>
{:else if job}
	<div class="mx-auto pb-10">
		<div class="flex gap-6">
			<!-- Main Content Container -->
			<div class="flex-1">
				<!-- Content Grid -->
				<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
					<!-- Left Column - Combined Job Header and Details -->
					<div class="lg:col-span-2">
						<div class="rounded-xl p-8">
							<!-- Job Header Section -->
							<div
								class="mb-5 flex flex-col gap-4 border-b border-gray-100 pb-5 sm:flex-row sm:items-start sm:justify-between"
							>
								<div class="flex-1">
									<div class="mb-3 flex items-center gap-3">
										<img src={job.logo} alt={job.company} class="h-10 w-10 rounded-lg" />
										<div>
											<div class="flex items-center gap-2 text-sm text-gray-600">
												<span class="text-lg font-medium">{job.company}</span>
											</div>
										</div>
									</div>
									<h1 class="mb-3 text-2xl font-semibold tracking-tight text-gray-900">
										{job.title}
									</h1>

									<div class="mt-4 flex flex-col gap-2.5 text-sm text-gray-600">
										<div class="flex items-center gap-2">
											<MapPin class="h-4 w-4 text-gray-400" />
											<span
												>{job.location}{job.workType === 'remote'
													? ' (Remote)'
													: job.workType === 'hybrid'
														? ' (Hybrid)'
														: ''}</span
											>
										</div>
										<div class="flex items-center gap-2">
											<Clock class="h-4 w-4 text-gray-400" />
											<span
												>{job.workArrangement
													.replace(/-/g, ' ')
													.toLowerCase()
													.replace(/^\w/, (c: string) => c.toUpperCase())}</span
											>
										</div>
										{#if job.salary}
											<div class="flex items-center gap-2">
												<Banknote class="h-4 w-4 text-gray-400" />
												<span>{job.salary}</span>
											</div>
										{/if}
										<div class="mt-2 flex items-center gap-2">
											<span>Posted {job.posted}</span>
										</div>
									</div>
								</div>

								<div class="flex gap-2">
									<button
										onclick={toggleBookmark}
										class="flex items-center rounded-md border border-gray-200 px-2 py-2 text-sm font-medium transition-colors hover:cursor-pointer hover:bg-gray-50 {isBookmarked
											? 'text-green-700'
											: 'text-gray-700'}"
									>
										<Bookmark class="h-4 w-4 {isBookmarked ? 'fill-current' : ''}" />
									</button>
									<button
										onclick={handleShare}
										class="rounded-md border border-gray-200 px-2 py-2 text-sm font-medium text-gray-700 transition-colors hover:cursor-pointer hover:bg-gray-50"
									>
										<Share2 class="h-4 w-4" />
									</button>
								</div>
							</div>

							<!-- Job Description Section -->
							<h2 class="mb-4 text-lg font-semibold text-gray-900">Job Description</h2>
							<div class="space-y-5 text-sm leading-7 text-gray-700">
								<SafeHTML html={job.description} />
							</div>

							<!-- Skills Section -->
							{#if job.skills && job.skills.length > 0}
								<div class="mt-6 pt-6">
									<h2 class="mb-4 text-lg font-semibold text-gray-900">Skills</h2>
									<div class="flex flex-wrap gap-2">
										{#each showAllSkills ? job.skills : job.skills as skill, index (index)}
											<SkillTag {skill} />
										{/each}
									</div>
								</div>
							{/if}

							<!-- About the Company Section -->
							{#if company}
								<div class="mt-6 pt-6">
									<h2 class="mb-4 text-lg font-semibold text-gray-900">About the Company</h2>
									<CompanyCard {company} />
								</div>
							{/if}
						</div>
					</div>

					<!-- Right Column -->
					<div class="lg:col-span-1">
						<div bind:this={applyButtonRef} class="mb-6 overflow-hidden rounded-xl bg-gray-100 p-6">
							<h2 class="mb-2 text-lg font-medium text-gray-900">Ready to apply?</h2>
							<p class="mb-4 text-sm text-gray-600">
								Take the next step in your career. Your application will be sent directly to the
								employer for review.
							</p>
							<div class="mb-3 text-xs text-gray-500">
								Application deadline: {job.closeDate || 'Open until filled'}
							</div>
							<ApplyButton
								jobId={job.id}
								closeDateRaw={job.closeDateRaw}
								postedDate={job.postedDate}
								onClick={handleApplyClick}
								size="lg"
								fullWidth={true}
							/>
						</div>
						<!-- Similar Jobs -->
						{#if similarJobs.length > 0}
							<div class="mb-6 overflow-hidden rounded-xl border border-gray-200 bg-white p-4">
								<h2 class="mb-3 text-lg font-medium text-gray-900">Similar Jobs</h2>
								<div class="space-y-3">
									{#each similarJobs as similarJob (similarJob.id)}
										<JobCard
											job={similarJob}
											onclick={() => navigateToJob(similarJob.id)}
										/>
									{/each}
								</div>
							</div>
						{/if}

						<!-- Other Jobs from Company -->
						{#if otherCompanyJobs.length > 0}
							<div class="overflow-hidden rounded-xl border border-gray-200 bg-white p-4">
								<h2 class="mb-3 text-lg font-medium text-gray-900">
									Other jobs from {job.company}
								</h2>
								<div class="space-y-3">
									{#each otherCompanyJobs as otherJob (otherJob.id)}
										<JobCard
											job={otherJob}
											onclick={() => navigateToJob(otherJob.id)}
										/>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="flex min-h-[400px] items-center justify-center">
		<p class="text-gray-600">Job not found.</p>
	</div>
{/if}

<!-- Floating Job Header -->
<FloatingJobHeader
	show={showFloatingCard}
	{job}
	onApply={handleApplyClick}
	onBookmark={toggleBookmark}
	onShare={handleShare}
	{isBookmarked}
/>

<!-- Apply Modal -->
{#if job}
	<ApplyModal bind:isOpen={showApplyModal} {job} />
{/if}
