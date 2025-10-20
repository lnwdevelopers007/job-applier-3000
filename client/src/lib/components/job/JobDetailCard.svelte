<script lang="ts">
	import { MapPin, Clock, Calendar, Banknote, Share2, ExternalLink, Bookmark } from 'lucide-svelte';
	import SafeHTML from '$lib/utils/SafeHTML.svelte';
	import Badge from './Badge.svelte';
	import SkillTag from './SkillTag.svelte';
	import CompanyCard from './CompanyCard.svelte';
	
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
		isBookmarked = false
	}: {
		job: Job;
		companyInfo?: CompanyInfo | null;
		onApply: (job: Job) => void;
		onBookmark?: () => void;
		isBookmarked?: boolean;
	} = $props();
	
	// Utility functions
	function getRelativeTime(dateString: string): string {
		if (!dateString || dateString === 'Unknown') return 'Unknown';
		
		const date = new Date(dateString);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - date.getTime());
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		
		if (date > now) {
			if (diffDays === 0) return 'Today';
			if (diffDays === 1) return 'Tomorrow';
			if (diffDays < 7) return `In ${diffDays} days`;
			if (diffDays < 30) return `In ${Math.floor(diffDays / 7)} weeks`;
			return date.toLocaleDateString();
		} else {
			if (diffDays === 0) return 'Today';
			if (diffDays === 1) return 'Yesterday';
			if (diffDays < 7) return `${diffDays} days ago`;
			if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
			if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
			return date.toLocaleDateString();
		}
	}
	
	const isClosed = $derived(new Date(job.closeDate) < new Date());
	const isNotOpen = $derived(new Date(job.posted) > new Date());
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
				<span>Posted {getRelativeTime(job.posted)}</span>
			</div>
		</div>

		<!-- Job Type and Work Arrangement Badges -->
		<div class="flex flex-wrap gap-2 mb-4">
			<Badge variant="success">{job.workArrangement || 'FULL-TIME'}</Badge>
			<Badge variant={job.workType === 'remote' ? 'info' : job.workType === 'hybrid' ? 'secondary' : 'warning'}>
				{job.workType?.toUpperCase() || 'ON-SITE'}
			</Badge>
		</div>

		<!-- Action Buttons -->
		<div class="flex gap-2">
			{#if isClosed}
				<button
					class="px-4 py-2 text-sm font-medium bg-gray-400 text-white rounded-md cursor-not-allowed"
					disabled
				>
					Closed
				</button>
			{:else if isNotOpen}
				<button
					class="px-4 py-2 text-sm font-medium bg-gray-400 text-white rounded-md cursor-not-allowed"
					disabled
				>
					Not Open Yet
				</button>
			{:else}
				<button
					class="px-4 py-2 text-sm font-medium bg-green-600 text-white rounded-md hover:bg-green-700 active:scale-[0.98] transition-all duration-150 hover:cursor-pointer"
					onclick={() => onApply(job)}
				>
					Apply now
				</button>
				{#if onBookmark}
					<button 
						onclick={onBookmark}
						class="px-3 py-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center {isBookmarked ? 'text-green-700 bg-green-50 border-green-200' : 'text-gray-700'} hover:cursor-pointer"
					>
						<Bookmark class="w-4 h-4 {isBookmarked ? 'fill-current' : ''}" />
					</button>
				{/if}
			{/if}
		</div>
	</div>

	<!-- Job Description -->
	<div class="mb-6">
		<h2 class="text-lg font-semibold text-gray-900 mb-3">Job Description</h2>
		<div class="text-sm text-gray-700 leading-6">
			<SafeHTML html={displayDescription} />
		</div>
	</div>

	<!-- Skills/Tags Section -->
	{#if job.tags && job.tags.length > 0}
		<div class="mb-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-3">Skills Required</h3>
			<div class="flex flex-wrap gap-2">
				{#each (showAllSkills ? job.tags : job.tags.slice(0, 6)) as tag}
					<SkillTag skill={tag} />
				{/each}
			</div>
			{#if job.tags.length > 6}
				<button 
					onclick={() => showAllSkills = !showAllSkills}
					class="mt-2 text-sm text-green-600 hover:text-green-700 transition-colors font-medium"
				>
					{showAllSkills ? 'Show less' : `Show ${job.tags.length - 6} more skills`}
				</button>
			{/if}
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


<style>
	.job-description :global(h2) {
		font-size: 1.125rem;
		font-weight: 600;
		color: rgb(31, 41, 55);
		margin-bottom: 0.5rem;
	}
	
	.job-description :global(p) {
		margin-bottom: 0.75rem;
		color: rgb(55, 65, 81);
	}
	
	.job-description :global(ul) {
		list-style-type: disc;
		padding-left: 1.25rem;
		margin-bottom: 0.75rem;
		color: rgb(55, 65, 81);
	}
	
	.job-description :global(li) {
		margin-bottom: 0.25rem;
	}
	
	.job-description :global(strong) {
		font-weight: 600;
		color: rgb(17, 24, 39);
	}
</style>