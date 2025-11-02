<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';
	import ConfirmActionWithReason from '$lib/components/modals/ConfirmActionWithReason.svelte';
	import ConfirmDropdownAction from '$lib/components/modals/ConfirmDropdownAction.svelte';
	import FilterPill from '$lib/components/forms/FilterPill.svelte';

	const USER_ACTIONS = [
		// { label: 'View', disabled: false },
		{ label: 'Edit Permissions', disabled: false },
		{ label: 'Ban', disabled: false },
		{ label: 'Delete', disabled: false }
	];
	const TABLE_HEADER = ['Name', 'Email', 'Role', 'Verified'];

	const VALID_ROLES = ['jobSeeker', 'company', 'faculty', 'admin'];
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

<div
	class="-mt-5 pb-8 pt-8"
	style="margin-left: calc(-50vw + 50%); margin-right: calc(-50vw + 50%); padding-left: calc(50vw - 50%); padding-right: calc(50vw - 50%); background: radial-gradient(circle,rgba(245, 255, 252, 1) 0%, rgba(248, 255, 249, 1) 25%, rgba(243, 255, 245, 1) 50%, rgba(237, 254, 244, 1) 75%, rgba(232, 254, 240, 1) 100%);"
>
	<div class="mx-auto max-w-7xl space-y-4">
		<!-- Search Bar -->
		<!-- <div class="relative mx-auto max-w-3xl">
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
				<Search class="h-5 w-5 text-gray-400" />
			</div>
			<input
				type="text"
				placeholder="Search for jobs..."
				class="w-full rounded-full border border-gray-200 bg-white py-3 pl-12 pr-4 text-gray-900 placeholder-gray-500 shadow-sm outline-none transition-all duration-200 focus:border-transparent focus:ring-1 focus:ring-gray-300"
				bind:value={searchQuery}
				oninput={refreshJobs}
			/>
		</div> -->

		<!-- Filter Pills -->
		<div class="flex flex-wrap items-center justify-center gap-3">
			<div class="flex items-center gap-2">
				<span class="text-sm font-medium text-gray-600">Filters:</span>

				<FilterPill
					label="Role"
					options={VALID_ROLES.map((role) => ({ value: role, label: role }))}
					bind:selectedValue={VALID_ROLES[0]}
					onSelectionChange={(value: string) => {}}
				/>

				<FilterPill
					label="Verified"
					options={VALID_VERIFICATION_OPTIONS.map((opt: boolean) => ({
						value: opt.toString(),
						label: opt.toString()
					}))}
					bind:selectedValue={VALID_ROLES[0]}
					onSelectionChange={(value: string) => {}}
				/>

				<!-- sort -->

				<!-- <div class="h-4 w-px bg-gray-300"></div> -->

				<!-- <FilterPill
					label="Sort"
					options={sortOptions}
					bind:selectedValue={sortBy}
					onSelectionChange={refreshJobs}
					type="sort"
				/> -->
			</div>
		</div>
	</div>
</div>

<TableWithAction
	things={users}
	tableHeader={TABLE_HEADER}
	rowAttributes={TABLE_HEADER.map((attribute: string) => {
		return attribute.toLowerCase();
	})}
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
