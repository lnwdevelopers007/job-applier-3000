<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';

	const USER_ACTIONS = [
		{ label: 'View', disabled: false },
		{ label: 'Edit', disabled: false },
		{ label: 'Ban', disabled: false },
		{ label: 'Delete', disabled: false }
	];
	const TABLE_HEADER = ['Name', 'Email', 'Role', 'Verified'];
	const ROW_ATTRIBUTES = TABLE_HEADER.map((attribute: string) => {
		return attribute.toLowerCase();
	});

	let selectedThing = $state<number | null>(null);
	let things = $state<any[]>([]);

	function selectRow(index: number) {
		selectedThing = selectedThing === index ? null : index;
	}

	function handleAction(arg: any, anotherArg: any) {
		return null;
	}

	async function loadUserData() {
		const usersFromDB = await fetchUsers();
		for (let user of usersFromDB) {
			user.actions = USER_ACTIONS;
		}
		things = usersFromDB;
	}

	$effect(() => {
		loadUserData();
	});
</script>

<TableWithAction
	{things}
	table_header={TABLE_HEADER}
	row_attributes={ROW_ATTRIBUTES}
	{selectRow}
	{selectedThing}
	{handleAction}
/>
