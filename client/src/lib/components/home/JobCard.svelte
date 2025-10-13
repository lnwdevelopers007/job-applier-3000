<script lang="ts">
	import { MapPin, Clock } from 'lucide-svelte';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { isAuthenticated, navigateWithAuth } from '$lib/utils/auth';

	type Job = {
		id: string;
		company: string;
		companyLogo?: string; // Will be populated from backend in the future
		title: string;
		location: string;
		locationType: 'on-site' | 'remote' | 'hybrid';
		minSalary: number;
		maxSalary: number;
		currency: string;
		type: string;
		tags: string[];
		postedAt: string;
		badge?: {
			text: string;
			type: 'new' | 'remote' | 'internship';
		};
		logoStyle?: string; // Placeholder styling - remove when companyLogo is implemented
	};

	let { job, skeleton = false }: { job?: Job; skeleton?: boolean } = $props();
	let showAuthModal = $state(false);

	const badgeColors = {
		new: 'bg-green-600 text-white',
		remote: 'bg-blue-600 text-white',
		internship: 'bg-purple-600 text-white'
	};

	function getCompanyInitial(name: string): string {
		return name.charAt(0).toUpperCase();
	}

	function handleApply() {
		if (isAuthenticated()) {
			// User is logged in - proceed to job application
			navigateWithAuth('/app/jobs');
		} else {
			// User not logged in - show auth modal
			showAuthModal = true;
			// Store job ID for after login
			sessionStorage.setItem('pendingJobApplication', job.id);
		}
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		// Clear pending application if user cancels
		sessionStorage.removeItem('pendingJobApplication');
	}

	function formatSalary(minSalary: number, maxSalary: number, currency: string): string {
		const formatNumber = (num: number) => {
			if (num >= 1000000) {
				return (num / 1000000).toFixed(1) + 'M';
			}
			if (num >= 1000) {
				return (num / 1000).toFixed(0) + 'k';
			}
			return num.toString();
		};

		const min = formatNumber(minSalary);
		const max = formatNumber(maxSalary);
		const symbol = currency === 'THB' ? '฿' : currency === 'USD' ? '$' : currency;
		
		return `${symbol}${min} - ${symbol}${max}`;
	}
</script>

{#if skeleton}
<div class="relative bg-white border border-gray-200 rounded-xl p-6 hover:shadow-sm transition-all duration-200 group flex flex-col h-full animate-pulse">
	<!-- Company logo and info skeleton -->
	<div class="flex items-start gap-4 mb-5">
		<div class="w-12 h-12 rounded-lg bg-gray-200 flex-shrink-0"></div>
		<div class="flex-1 min-w-0">
			<div class="h-4 bg-gray-200 rounded mb-2 w-24"></div>
			<div class="h-6 bg-gray-200 rounded w-48"></div>
		</div>
	</div>
	
	<!-- Location skeleton -->
	<div class="flex items-center gap-2 mb-3">
		<div class="w-4 h-4 bg-gray-200 rounded"></div>
		<div class="h-4 bg-gray-200 rounded w-32"></div>
	</div>
	
	<!-- Salary and type skeleton -->
	<div class="flex items-center gap-2 mb-5">
		<div class="h-4 bg-gray-200 rounded w-24"></div>
		<div class="h-6 bg-gray-200 rounded w-16"></div>
		<div class="h-6 bg-gray-200 rounded w-20"></div>
	</div>
	
	<!-- Tags skeleton -->
	<div class="flex flex-wrap gap-2 mb-4">
		<div class="h-6 bg-gray-200 rounded w-16"></div>
		<div class="h-6 bg-gray-200 rounded w-20"></div>
		<div class="h-6 bg-gray-200 rounded w-14"></div>
	</div>
	
	<!-- Spacer -->
	<div class="flex-grow"></div>
	
	<!-- Footer skeleton -->
	<div class="flex items-center justify-between pt-4 border-t border-gray-100 mt-auto">
		<div class="h-3 bg-gray-200 rounded w-24"></div>
		<div class="h-9 bg-gray-200 rounded w-20"></div>
	</div>
</div>
{:else}
<div
	class="relative bg-white border border-gray-200 rounded-xl p-6 hover:shadow-sm transition-all duration-200 group flex flex-col h-full"
>
	{#if job?.badge}
		<span
			class="absolute top-4 right-4 px-2.5 py-1 text-xs font-bold uppercase tracking-wider rounded {badgeColors[
				job.badge.type
			]}"
		>
			{job.badge.text}
		</span>
	{/if}

	<div class="flex items-start gap-4 mb-4">
		{#if job?.companyLogo}
			<div
				class="w-12 h-12 rounded-lg flex items-center justify-center text-xl font-semibold flex-shrink-0 {job.logoStyle ||
					'bg-gray-100'}"
			>
				{getCompanyInitial(job.company)}
			</div>
		{:else}
			<div
				class="w-12 h-12 rounded-lg bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center text-gray-600 font-semibold flex-shrink-0"
			>
				{job ? getCompanyInitial(job.company) : ''}
			</div>
		{/if}
		<div class="flex-1 min-w-0">
			<p class="text-sm font-semibold text-gray-600 mb-1">{job?.company}</p>
			<h3 class="text-lg font-semibold text-gray-900 transition-colors">
				{job?.title}
			</h3>
		</div>
	</div>

	<div class="flex items-center gap-2 mb-3 text-sm text-gray-600">
		{#if job?.locationType === 'remote'}
			<div class="flex items-center gap-1">
				<MapPin class="w-4 h-4" />
				<span>Remote (Worldwide)</span>
			</div>
		{:else}
			<div class="flex items-center gap-1">
				<MapPin class="w-4 h-4" />
				<span>{job?.location}</span>
			</div>
		{/if}
	</div>

	<div class="flex items-center gap-2 mb-4 text-sm">
		<span class="font-semibold text-green-600 flex items-center gap-1 mr-2">
			{job ? formatSalary(job.minSalary, job.maxSalary, job.currency) : ''}
		</span>
		{#if job?.type}
			<span class="px-2 py-1 bg-purple-100 text-purple-700 text-xs rounded font-medium">{job.type}</span>
		{/if}
		{#if job?.locationType}
			<span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded font-medium capitalize">{job.locationType}</span>
		{/if}
	</div>

	<div class="flex flex-wrap gap-2 mb-4">
		{#each job?.tags || [] as tag, index (tag + index)}
			<span class="px-3 py-1 bg-gray-100 text-gray-700 text-xs rounded-md font-medium">
				{tag}
			</span>
		{/each}
	</div>

	<!-- Spacer to push footer to bottom -->
	<div class="flex-grow"></div>

	<div class="flex items-center justify-between pt-4 border-t border-gray-100 mt-auto">
		<span class="text-xs text-gray-500 flex items-center gap-1">
			<Clock class="w-3 h-3" />
			{job?.postedAt}
		</span>
		<button
			onclick={handleApply}
			class="px-5 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 cursor-pointer transition-colors"
		>
			Apply
		</button>
	</div>
</div>
{/if}

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title="Sign in to apply"
	description="Please log in or sign up to apply for this job position"
/>