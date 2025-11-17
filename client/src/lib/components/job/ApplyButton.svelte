<script lang="ts">
	import { JobApplicationService } from '$lib/services/jobApplicationService';
	import { getUserInfo } from '$lib/utils/auth';
	import { onMount } from 'svelte';

	interface Props {
		isApplied?: boolean;
		jobId?: string;
		closeDateRaw?: string | null;
		postedDate?: string | null;
		closeDate?: string;
		posted?: string;
		onClick: () => void;
		class?: string;
		size?: 'sm' | 'md' | 'lg';
		fullWidth?: boolean;
	}
	
	let {
		isApplied = false,
		jobId,
		closeDateRaw = null,
		postedDate = null,
		closeDate = '',
		posted = '',
		onClick,
		class: className = '',
		size = 'md',
		fullWidth = false
	}: Props = $props();
	
	let actualIsApplied = $state(isApplied);
	let isCheckingStatus = $state(false);
	
	// If jobId is provided, check the applied status directly
	async function checkApplicationStatus() {
		if (!jobId) {
			actualIsApplied = isApplied;
			return;
		}
		
		try {
			isCheckingStatus = true;
			const user = getUserInfo();
			if (!user?.userID) {
				actualIsApplied = false;
			} else {
				actualIsApplied = await JobApplicationService.hasUserAppliedToJob(user.userID, jobId);
			}
		} catch (error) {
			console.error('Error checking application status:', error);
			actualIsApplied = isApplied; // fallback to prop value
		} finally {
			isCheckingStatus = false;
		}
	}
	
	// Check on mount and when jobId changes
	onMount(() => {
		checkApplicationStatus();
	});
	
	$effect(() => {
		checkApplicationStatus();
	});
	
	const isClosed = $derived(() => {
		if (closeDateRaw) return new Date(closeDateRaw) < new Date();
		if (closeDate) return new Date(closeDate) < new Date();
		return false;
	});
	
	const isNotOpen = $derived(() => {
		if (postedDate) return new Date(postedDate) > new Date();
		if (posted) return new Date(posted) > new Date();
		return false;
	});
	
	const sizeClasses = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-6 py-2 text-sm',
		lg: 'px-6 py-2 text-sm'
	};
	
	const baseClasses = $derived(
		`${sizeClasses[size]} font-medium rounded-md transition-all duration-150 ${fullWidth ? 'w-full' : ''} ${className}`
	);
</script>

{#if actualIsApplied}
	<button
		disabled
		class="{baseClasses} bg-green-700 text-white cursor-not-allowed"
	>
		Applied
	</button>
{:else if isClosed()}
	<button
		disabled
		class="{baseClasses} bg-gray-400 text-white cursor-not-allowed"
	>
		Closed
	</button>
{:else if isNotOpen()}
	<button
		disabled
		class="{baseClasses} bg-gray-400 text-white cursor-not-allowed"
	>
		Not Open Yet
	</button>
{:else}
	<button
		onclick={onClick}
		class="{baseClasses} bg-green-600 text-white hover:bg-green-700 active:scale-[0.98]"
	>
		{size === 'lg' ? 'Apply now' : 'Apply'}
	</button>
{/if}