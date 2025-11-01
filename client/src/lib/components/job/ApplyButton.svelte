<script lang="ts">
	import type { UserRole } from '$lib/stores/auth.svelte';
	
	interface Props {
		isApplied?: boolean;
		closeDateRaw?: string | null;
		postedDate?: string | null;
		closeDate?: string;
		posted?: string;
		onClick: () => void;
		class?: string;
		size?: 'sm' | 'md' | 'lg';
		fullWidth?: boolean;
		userRole?: UserRole | null;
	}
	
	let {
		isApplied = false,
		closeDateRaw = null,
		postedDate = null,
		closeDate = '',
		posted = '',
		onClick,
		class: className = '',
		size = 'md',
		fullWidth = false,
		userRole = null
	}: Props = $props();
	
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
	
	const isJobSeeker = $derived(userRole === 'jobSeeker');
</script>

{#if isJobSeeker}
{#if isApplied}
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
{/if}