<script lang="ts">
	import { MapPin, Banknote, Share2, ExternalLink, Bookmark } from 'lucide-svelte';
	import ApplyButton from './ApplyButton.svelte';
	import SafeHTML from '$lib/utils/SafeHTML.svelte';
	import Badge from './Badge.svelte';
	import WorkInfoBadge from './WorkInfoBadge.svelte';
	import SkillTag from './SkillTag.svelte';
	import CompanyCard from './CompanyCard.svelte';
	import { formatRelativeTime } from '$lib/utils/datetime';
	
	interface Job {
		id: string;
		title: string;
		company: string;
		logo: string;
		location: string;
		workType?: string;
		workArrangement?: string;
		salary?: string;
		posted: string;
		closeDate: string;
		description?: string;
		tags?: string[];
	}
	
	interface CompanyInfo {
		name?: string;
		logo?: string;
		location?: string;
		size?: string;
		industry?: string;
		aboutUs?: string;
	}
	
	let {
		job,
		companyInfo = null,
		onApply,
		onBookmark = null,
		isBookmarked = false,
		isApplied = false
	}: {
		job: Job;
		companyInfo?: CompanyInfo | null;
		onApply: (job: Job) => void;
		onBookmark?: () => void;
		isBookmarked?: boolean;
		isApplied?: boolean;
	} = $props();
	
	
	const displayDescription = $derived(job.description || '');
	
	let showAllSkills = $state(false);
	
	
	// Transform companyInfo to match CompanyCard interface
	const company = $derived(companyInfo ? {
		name: companyInfo.name || job.company,
		logo: companyInfo.logo || job.logo,
		industry: companyInfo.industry || 'Software Development',
		employees: companyInfo.size || '1000+ employees',
		description: companyInfo.aboutUs || "Company information not available."
	} : null);


	function handleShare() {
		if (navigator.share) {
			navigator.share({ 
				title: job.title || 'Job Opportunity',
				text: `Check out this job: ${job.title} at ${job.company}`,
				url: `${window.location.origin}/app/jobs/${job.id}`
			});
		} else {
			navigator.clipboard.writeText(`${window.location.origin}/app/jobs/${job.id}`);
			// Could add toast notification here if available
		}
	}

</script>

<div class="p-8">
	<!-- Header Section -->
	<div class="mb-12">
		<div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-3">
			<div class="flex-1">
				<div class="flex items-center gap-3 mb-3">
					<img 
						src={job.logo} 
						alt={job.company} 
						class="w-10 h-10 rounded-lg object-cover"
					/>
					<div>
						<span class="font-medium text-lg text-gray-900">{job.company}</span>
					</div>
				</div>
				<h1 class="text-xl font-semibold tracking-tight text-gray-900">{job.title}</h1>
			</div>

			<!-- Top Right Buttons -->
			<div class="flex gap-2">
				<a 
					href={`/app/jobs/${job.id}`}
					class="px-2 py-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center gap-2 hover:cursor-pointer"
				>
					<ExternalLink class="w-4 h-4" />
				</a>
				<button 
					onclick={handleShare}
					class="px-2 py-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center hover:cursor-pointer"
				>
					<Share2 class="w-4 h-4" />
				</button>
			</div>
		</div>

		<div class="flex flex-col gap-2 text-sm text-gray-600 mb-4">
			<div class="flex items-center gap-2">
				<MapPin class="w-4 h-4 text-gray-400" />
				<span>{job.location}{job.workType === 'remote' ? ' (Remote)' : job.workType === 'hybrid' ? ' (Hybrid)' : ''}</span>
			</div>
			{#if job.salary}
				<div class="flex items-center gap-2">
					<Banknote class="w-4 h-4 text-gray-400" />
					<span>{job.salary}</span>
				</div>
			{/if}
			<div class="flex items-center gap-2 mt-1">
				<span>Posted {formatRelativeTime(job.posted)}</span>
			</div>
		</div>

		<!-- Job Type and Work Arrangement Badges -->
		<div class="flex flex-wrap gap-2 mb-4">
			<WorkInfoBadge type="workArrangement" value={job.workArrangement || 'on-site'} />
			<WorkInfoBadge type="workType" value={job.workType || 'full-time'} />
		</div>

		<!-- Action Buttons -->
		<div class="flex gap-2">
			<ApplyButton
				{isApplied}
				closeDate={job.closeDate}
				posted={job.posted}
				onClick={() => onApply(job)}
			/>
			{#if onBookmark}
				<button 
					onclick={onBookmark}
					class="p-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center {isBookmarked ? 'text-green-700' : 'text-gray-700'} hover:cursor-pointer"
				>
					<Bookmark class="w-4 h-4 {isBookmarked ? 'fill-current' : ''}" />
				</button>
			{/if}
		</div>
	</div>

	<!-- Job Description -->
	<div class="mb-6">
		<h2 class="text-lg font-semibold text-gray-900 mb-3">Job Description</h2>
		<div class="text-sm text-gray-700 leading-6 line-clamp-15">
			<SafeHTML html={displayDescription} />
		</div>
		{#if displayDescription.length > 1250}
			<a 
				href="/app/jobs/{job.id}"
				class="text-sm text-green-600 hover:text-green-700 transition-colors font-medium mt-2 inline-block"
			>
				Show more
			</a>
		{/if}
	</div>

	<!-- Skills/Tags Section -->
	{#if job.tags && job.tags.length > 0}
		<div class="mb-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-3">Skills</h3>
			<div class="flex flex-wrap gap-2">
				{#each (showAllSkills ? job.tags : job.tags) as tag, index (index)}
					<SkillTag skill={tag} />
				{/each}
			</div>
		</div>
	{/if}

	<!-- About the Company Section -->
	{#if company}
		<div class="pt-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">About the Company</h3>
			<CompanyCard {company} />
		</div>
	{/if}
</div>
