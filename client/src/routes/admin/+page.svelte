<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';
	import ConfirmActionWithReason from '$lib/components/modals/ConfirmActionWithReason.svelte';
	import ConfirmDropdownAction from '$lib/components/modals/ConfirmDropdownAction.svelte';

	const USER_ACTIONS = [
		{ label: 'View', disabled: false },
		{ label: 'Edit Permissions', disabled: false },
		{ label: 'Ban', disabled: false },
		{ label: 'Delete', disabled: false }
	];
	const TABLE_HEADER = ['Name', 'Email', 'Role', 'Verified'];
	const ROW_ATTRIBUTES = TABLE_HEADER.map((attribute: string) => {
		return attribute.toLowerCase();
	});

	const VALID_ROLES = ['jobSeeker', 'company', 'faculty'];

	let selectedUserIndex = $state<number | null>(null);
	let users = $state<any[]>([]);

	let NameOfSelectedUser = $state('');

	let showDeleteModal = $state(false);
	let deleteReason = $state('');
	let isDeleting = $state(false);

	let showBanModal = $state(false);
	let banReason = $state('');
	let isBanning = $state(false);

	let showRoleEditModal = $state(false);

	async function deleteUserConfirmed() {
		// TODO: add a real delete operation here
		showDeleteModal = false;
	}

	function handleAction(action: any, user: any) {
		switch (action.label) {
			case 'View':
				showRoleEditModal = true;
				break;
			case 'Edit Permissions':
				console.log('Edit Permissions');
				break;
			case 'Ban':
				NameOfSelectedUser = user.name;
				showBanModal = true;
				break;
			case 'Delete':
				NameOfSelectedUser = user.name;
				showDeleteModal = true;
				break;
			default:
				console.log('Nani ga suki?');
		}
	}

	async function loadUserData() {
		const usersFromDB = await fetchUsers();
		for (let user of usersFromDB) {
			user.actions = USER_ACTIONS;
		}
		users = usersFromDB;
	}

	$effect(() => {
		loadUserData();
	});
</script>

<TableWithAction
	things={users}
	selectedThingIndex={selectedUserIndex}
	table_header={TABLE_HEADER}
	row_attributes={ROW_ATTRIBUTES}
	{handleAction}
/>

<ConfirmActionWithReason
	bind:isVisible={showDeleteModal}
	thing={'User'}
	actionName={'Delete'}
	thingName={NameOfSelectedUser}
	isActionInProgress={isDeleting}
	reasonForAction={deleteReason}
	action={deleteUserConfirmed}
/>

<ConfirmActionWithReason
	bind:isVisible={showBanModal}
	thing={'User'}
	actionName={'Ban'}
	thingName={NameOfSelectedUser}
	isActionInProgress={isBanning}
	reasonForAction={banReason}
	action={deleteUserConfirmed}
/>
