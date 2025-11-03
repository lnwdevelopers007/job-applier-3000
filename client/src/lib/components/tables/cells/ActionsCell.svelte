<script lang="ts">
	import { MoreVertical } from 'lucide-svelte';
	import Dropdown from '$lib/components/ui/Dropdown.svelte';

	interface Props {
		onView: () => void;
		onEdit: () => void;
		onDelete: () => void;
	}

	let { onView, onEdit, onDelete }: Props = $props();
	let showMenu = $state(false);
	let triggerElement = $state<HTMLElement>();

	const dropdownItems = [
		{
			label: 'View',
			action: onView,
			variant: 'default' as const
		},
		{
			label: 'Edit',
			action: onEdit,
			variant: 'default' as const
		},
		{
			label: 'Delete',
			action: onDelete,
			variant: 'danger' as const
		}
	];
</script>

<div class="relative">
	<button
		bind:this={triggerElement}
		class="p-1 rounded hover:bg-gray-100"
		onclick={(e) => {
			e.stopPropagation();
			showMenu = !showMenu;
		}}
	>
		<MoreVertical class="w-4 h-4" />
	</button>

	<Dropdown
		items={dropdownItems}
		isOpen={showMenu}
		onClose={() => (showMenu = false)}
		{triggerElement}
	/>
</div>