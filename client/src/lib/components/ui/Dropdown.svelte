<script lang="ts">
	import { scale } from 'svelte/transition';
	import { tick } from 'svelte';
	
	interface DropdownItem {
		label: string;
		action: () => void;
		variant?: 'default' | 'danger' | 'success' | 'warning';
	}
	
	interface Props {
		items: DropdownItem[];
		isOpen: boolean;
		onClose: () => void;
		triggerElement?: HTMLElement;
	}
	
	let { items, isOpen, onClose, triggerElement }: Props = $props();
	
	let dropdownRef = $state<HTMLDivElement>();
	let position = $state({ top: 0, left: 0 });
	
	function handleClickOutside(event: MouseEvent) {
		if (dropdownRef && !dropdownRef.contains(event.target as Node)) {
			onClose();
		}
	}
	
	async function updatePosition() {
		if (!triggerElement || !dropdownRef) return;
		
		await tick();
		const rect = triggerElement.getBoundingClientRect();
		const dropdownRect = dropdownRef.getBoundingClientRect();
		
		// Calculate position relative to viewport
		let top = rect.bottom + 4; // 4px gap
		let left = rect.right - dropdownRect.width;
		
		// Adjust if dropdown goes off screen
		if (top + dropdownRect.height > window.innerHeight) {
			top = rect.top - dropdownRect.height - 4;
		}
		
		if (left < 0) {
			left = rect.left;
		}
		
		position = { top, left };
	}
	
	$effect(() => {
		if (isOpen) {
			updatePosition();
			// Delay adding the event listener to avoid immediately closing
			setTimeout(() => {
				document.addEventListener('click', handleClickOutside);
				window.addEventListener('scroll', updatePosition, true);
				window.addEventListener('resize', updatePosition);
			}, 0);
		}
		
		return () => {
			document.removeEventListener('click', handleClickOutside);
			window.removeEventListener('scroll', updatePosition, true);
			window.removeEventListener('resize', updatePosition);
		};
	});
</script>

{#if isOpen}
	<div 
		bind:this={dropdownRef}
		class="fixed min-w-32 bg-white border border-gray-200 rounded-lg shadow-lg z-[100]"
		style="top: {position.top}px; left: {position.left}px;"
		in:scale={{ duration: 150, start: 0.95 }}
		out:scale={{ duration: 100, start: 0.95 }}
	>
		<div class="py-1 flex flex-col">
			{#each items as item, index (index)}
				<button
					onclick={() => {
						item.action();
						onClose();
					}}
					class={`w-full text-left px-4 py-2 text-sm transition-colors hover:cursor-pointer block ${
						item.variant === 'danger'
							? 'text-red-600 hover:bg-red-50'
							: item.variant === 'success'
								? 'text-green-600 hover:bg-green-50'
								: item.variant === 'warning'
									? 'text-orange-600 hover:bg-orange-50'
									: 'text-gray-700 hover:bg-gray-50'
					}`}
				>
					{item.label}
				</button>
			{/each}
		</div>
	</div>
{/if}