<script lang="ts">
	import { scale } from 'svelte/transition';
	
	interface DropdownItem {
		label: string;
		action: () => void;
		variant?: 'default' | 'danger';
	}
	
	interface Props {
		items: DropdownItem[];
		isOpen: boolean;
		onClose: () => void;
	}
	
	let { items, isOpen, onClose }: Props = $props();
	
	let dropdownRef = $state<HTMLDivElement>();
	
	function handleClickOutside(event: MouseEvent) {
		if (dropdownRef && !dropdownRef.contains(event.target as Node)) {
			onClose();
		}
	}
	
	$effect(() => {
		if (isOpen) {
			// Delay adding the event listener to avoid immediately closing
			setTimeout(() => {
				document.addEventListener('click', handleClickOutside);
			}, 0);
		}
		
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

{#if isOpen}
	<div 
		bind:this={dropdownRef}
		class="absolute right-0 top-full mt-1 min-w-32 bg-white border border-gray-200 rounded-lg shadow-lg z-10"
		in:scale={{ duration: 150, start: 0.95 }}
		out:scale={{ duration: 100, start: 0.95 }}
	>
		<div class="py-1">
			{#each items as item}
				<button
					onclick={() => {
						item.action();
						onClose();
					}}
					class={`w-full text-left px-4 py-2 text-sm transition-colors hover:cursor-pointer ${
						item.variant === 'danger'
							? 'text-red-600 hover:bg-red-50'
							: 'text-gray-700 hover:bg-gray-50'
					}`}
				>
					{item.label}
				</button>
			{/each}
		</div>
	</div>
{/if}