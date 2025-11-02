<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';
	import ConfirmActionWithReason from '$lib/components/modals/ConfirmActionWithReason.svelte';
	import ConfirmDropdownAction from '$lib/components/modals/ConfirmDropdownAction.svelte';

	const USER_ACTIONS = [
		// { label: 'View', disabled: false },
		{ label: 'Edit Permissions', disabled: false },
		{ label: 'Ban', disabled: false },
		{ label: 'Delete', disabled: false }
	];
	const TABLE_HEADER = ['Name', 'Email', 'Role', 'Verified'];
	const ROW_ATTRIBUTES = TABLE_HEADER.map((attribute: string) => {
		return attribute.toLowerCase();
	});

	const VALID_ROLES = ['jobSeeker', 'company', 'faculty'];
	const VALID_VERIFICATION_OPTIONS = [true, false];

	let users = $state<any[]>([]);
	let selectedUser = $state<any>(null);

	let showDeleteModal = $state(false);
	let deleteReason = $state('');
	let isDeleting = $state(false);

	let showBanModal = $state(false);
	let banReason = $state('');
	let isBanning = $state(false);

	let showPermissionEditModal = $state(false);
	let showPermissionEditConfirmButton = $state(false);

	let dropdowns = $derived([
		{
			name: 'Roles',
			values: VALID_ROLES,
			defaultVal: selectedUser === null ? '' : selectedUser.role
		},
		{
			name: 'Verified',
			values: VALID_VERIFICATION_OPTIONS,
			defaultVal: selectedUser === null ? '' : selectedUser.verified
		}
	]);

	async function onDeleteUser() {
		isDeleting = true;
		// TODO: add a real delete operation here
		showDeleteModal = false;
		isDeleting = false;
	}

	async function onBanUser() {
		isBanning = true;
		// TODO: add a real ban operation here
		showBanModal = false;
		isBanning = false;
	}

	function onConfirmPermissions(selectedVals: Record<string, any>) {
		if (!selectedUser) return;

		// 1. update selectedUser fields locally
		Object.entries(selectedVals).forEach(([key, val]) => {
			if (key === 'Roles') selectedUser.role = val;
			if (key === 'Verified') selectedUser.verified = val;
		});

		// 2. update the entire users array (immutably)
		users = users.map((u) => (u.id === selectedUser.id ? { ...u, ...selectedUser } : u));

		//TODO: Send updated data to server

		// 4. close modal
		showPermissionEditModal = false;
	}

	function handleAction(action: any, user: any) {
		selectedUser = user;
		switch (action.label) {
			// case 'View':
			// 	console.log('Skibidi');
			// 	break;
			case 'Edit Permissions':
				showPermissionEditModal = true;
				console.log(showPermissionEditModal);
				break;
			case 'Ban':
				showBanModal = true;
				break;
			case 'Delete':
				showDeleteModal = true;
				break;
			default:
				null;
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
	tableHeader={TABLE_HEADER}
	rowAttributes={ROW_ATTRIBUTES}
	{handleAction}
/>

<ConfirmActionWithReason
	bind:isVisible={showDeleteModal}
	actionName={'Delete'}
	actOnKind={'User'}
	actOnIndividual={selectedUser === null ? '' : selectedUser.name}
	bind:isActionInProgress={isDeleting}
	reasonForAction={deleteReason}
	action={onDeleteUser}
/>

<ConfirmActionWithReason
	bind:isVisible={showBanModal}
	actionName={'Ban'}
	actOnKind={'User'}
	actOnIndividual={selectedUser === null ? '' : selectedUser.name}
	bind:isActionInProgress={isBanning}
	reasonForAction={banReason}
	action={onBanUser}
/>

<ConfirmDropdownAction
	bind:isVisible={showPermissionEditModal}
	actionName={`Editing ${selectedUser === null ? '' : selectedUser.name} Permission`}
	bind:showConfirmButton={showPermissionEditConfirmButton}
	bind:dropdowns
	action={onConfirmPermissions}
	linkData={selectedUser
		? {
				title: 'View Credentials of this User',
				link: 'https://www.google.com'
			}
		: null}
/>
