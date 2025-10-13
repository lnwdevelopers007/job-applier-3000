<script lang="ts">
	import { X } from 'lucide-svelte';
	import { fade, fly } from 'svelte/transition';
	
	let { 
		isOpen = $bindable(false),
		title = 'Preview',
		subtitle = '',
		width = '820px',
		children
	} = $props();
	
	function closeDrawer() {
		isOpen = false;
	}
</script>

{#if isOpen}
	<!-- Backdrop -->
	<div 
		class="fixed inset-0 bg-black/50 z-60"
		onclick={closeDrawer}
		onkeydown={(e) => e.key === 'Escape' && closeDrawer()}
		role="button"
		tabindex="0"
		transition:fade={{ duration: 300 }}
	></div>
	
	<!-- Drawer -->
	<div 
		class="fixed right-0 top-0 h-full bg-white shadow-xl z-70 overflow-y-auto"
		style="width: {width}"
		transition:fly={{ x: parseInt(width), duration: 300, opacity: 1 }}
	>
		<!-- Header -->
		<div class="sticky top-0 bg-white px-6 py-4 flex items-center justify-between border-b border-gray-200">
			<div>
				<h2 class="text-lg font-semibold text-gray-900">{title}</h2>
				{#if subtitle}
					<p class="text-sm text-gray-500 mt-1">{subtitle}</p>
				{/if}
			</div>
			<button
				onclick={closeDrawer}
				class="p-1 hover:bg-gray-100 rounded-md transition-colors"
			>
				<X class="w-5 h-5 text-gray-500" />
			</button>
		</div>
		
		<!-- Content -->
		<div class="p-6">
			{@render children?.()}
		</div>
	</div>
{/if}