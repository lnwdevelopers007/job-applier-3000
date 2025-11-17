<script lang="ts">
	import { MoreVertical } from 'lucide-svelte';
	import Dropdown from '$lib/components/ui/Dropdown.svelte';

	interface CustomAction {
		label: string;
		action: () => void;
		variant: 'default' | 'danger' | 'warning' | 'success';
	}

	interface Props {
		onView?: () => void;
		onEdit?: () => void;
		onDelete?: () => void;
		customActions?: CustomAction[];
	}

	let { onView, onEdit, onDelete, customActions = [] }: Props = $props();
	let showMenu = $state(false);
	let triggerElement = $state<HTMLElement>();

	const dropdownItems = [
		onView && {
			label: 'View',
			action: onView,
			variant: 'default' as const
		},
		onEdit && {
			label: 'Edit',
			action: onEdit,
			variant: 'default' as const
		},
		...customActions,
		onDelete && {
			label: 'Delete',
			action: onDelete,
			variant: 'danger' as const
		}
	].filter((item): item is NonNullable<typeof item> => Boolean(item));
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