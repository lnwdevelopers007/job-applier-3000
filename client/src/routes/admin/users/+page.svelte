<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';
	import TableWithAction from '$lib/components/table/TableWithAction.svelte';
	import ConfirmActionWithReason from '$lib/components/modals/ConfirmActionWithReason.svelte';
	import ConfirmDropdownAction from '$lib/components/modals/ConfirmDropdownAction.svelte';
	import FilterPill from '$lib/components/forms/FilterPill.svelte';
	import { Search } from 'lucide-svelte';

	const USER_ACTIONS = [
		// { label: 'View', disabled: false },
		{ label: 'Edit Permissions', disabled: false },
		{ label: 'Ban', disabled: false },
		{ label: 'Delete', disabled: false }
	];
	const TABLE_HEADER = ['Name', 'Email', 'Role', 'Verified'];

	const VALID_ROLES = ['jobSeeker', 'company', 'faculty', 'admin'];
	const VALID_VERIFICATION_OPTIONS = [true, false];
	let currentFilteredBannedStatus = $state('');

	const SORT_OPTIONS = [
		{ value: 'name', label: 'Name' },
		{ value: 'email', label: 'Email' },
		{ value: 'role', label: 'Role' },
		{ value: 'verified', label: 'Verification Status' }
	];

	let users = $state<any[]>([]);
	let originalUsers = $state<any[]>([]);
	let selectedUser = $state<any>(null);

	let showDeleteModal = $state(false);
	let deleteReason = $state('');
	let isDeleting = $state(false);

	let showBanModal = $state(false);
	let banReason = $state('');
	let isBanning = $state(false);

	let showUnbanModal = $state(false);
	let unbanReason = $state('');
	let isUnbanning = $state(false);

	let showPermissionEditModal = $state(false);
	let showPermissionEditConfirmButton = $state(false);

	let currentFilteredRole = $state('');
	let currentFilteredVerificationStatus = $state('');

	let searchQuery = $state('');
	let sortBy = $state('');

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

		originalUsers = originalUsers.filter((u) => u.id !== selectedUser.id);
		users = users.filter((u) => u.id !== selectedUser.id);

		const res = await fetch(`/users/${selectedUser.id}`, {
			method: 'DELETE',
			credentials: 'include'
		});

		if (!res.ok) {
			console.log('Error: Deletion failed');
		}

		showDeleteModal = false;
		isDeleting = false;
	}

	async function onBanUser() {
		isBanning = true;

		try {
			const newBanStatus = !selectedUser.ban;
			const res = await fetch(`/users/${selectedUser.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ ban: newBanStatus })
			});

			if (!res.ok) {
				console.error('Error: Failed to update ban status');
				isBanning = false;
				return;
			}

			selectedUser.ban = newBanStatus;

			originalUsers = originalUsers.map((u) =>
				u.id === selectedUser.id
					? {
							...u,
							ban: newBanStatus,
							actions: u.actions.map((a: { label: string; }) =>
								a.label === 'Ban' || a.label === 'Unban'
									? {
											...a,
											label: newBanStatus ? 'Unban' : 'Ban'
									  }
									: a
							)
					  }
					: u
			);

			users = [...originalUsers];
			console.log(
				`User ${selectedUser.name} has been ${newBanStatus ? 'banned' : 'unbanned'} successfully ✅`
			);
		} catch (err) {
			console.error('Network error banning user:', err);
		} finally {
			showBanModal = false;
			isBanning = false;
			banReason = '';
		}
	}

	async function onConfirmPermissions(selectedVals: Record<string, any>) {
		if (!selectedUser) return;

		const role = selectedVals['Roles'];
		const verified = selectedVals['Verified'];

		// Update user role
		if (role !== undefined && role !== selectedUser.role) {
			const resRole = await fetch(`/users/${selectedUser.id}/role`, {
				method: 'PATCH',
        headers: {'Content-Type': 'application/json' },
				credentials: "include",
				body: JSON.stringify({ role })
			});
			if (!resRole.ok) console.error("Error: Can't update role");
			else selectedUser.role = role;
		}

		// Update verified status
		if (verified !== undefined && verified !== selectedUser.verified) {
			const resVerified = await fetch(`/users/${selectedUser.id}/verify`, {
				method: 'PATCH',
				headers: {'Content-Type': 'application/json' },
				credentials: "include",
				body: JSON.stringify({ verified })
			});
			if (!resVerified.ok) console.error("Error: Can't update verified status");
			else selectedUser.verified = verified;
		}

		// Update field locally
		Object.entries(selectedVals).forEach(([key, val]) => {
			if (key === 'Roles') selectedUser.role = val;
			if (key === 'Verified') selectedUser.verified = val;
		});

		// Update the user array
		originalUsers = originalUsers.map((u) =>
			u.id === selectedUser.id ? { ...u, ...selectedUser } : u
		);
		users = users.map((u) => (u.id === selectedUser.id ? { ...u, ...selectedUser } : u));

		// 4. close modal
		showPermissionEditModal = false;
	}

	function onFilteringRole(role: string) {
		users = role === '' ? originalUsers : originalUsers.filter((user: any) => user.role === role);
	}

	function onFilteringVerificationStatus(status: string) {
		let strToBool = (str: string) => (str === 'true' ? true : false);

		users =
			status === ''
				? originalUsers
				: originalUsers.filter((user: any) => user.verified === strToBool(status));
	}

	function onFilteringBannedStatus(status: string) {
		let strToBool = (str: string) => str === 'true';
		users =
			status === ''
				? originalUsers
				: originalUsers.filter((user: any) => user.ban === strToBool(status));
	}

	function onUserSearch() {
		users =
			searchQuery === ''
				? originalUsers
				: originalUsers.filter((user: any) =>
						user.name.toLowerCase().includes(searchQuery.toLowerCase())
					);
	}

	function onUserSort() {
		// Reset or apply sort
		users =
			sortBy === ''
				? originalUsers
				: [...users].sort((a: any, b: any) => {
						if (a[sortBy] < b[sortBy]) return -1;
						if (a[sortBy] > b[sortBy]) return 1;
						return 0;
					});
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
			case 'Unban':
				showUnbanModal = true;
				break;
			case 'Delete':
				showDeleteModal = true;
				break;
			default:
				console.log('hello there');
		}
	}

	async function loadUserData() {
		const usersFromDB = await fetchUsers();
		for (let user of usersFromDB) {
			user.actions = USER_ACTIONS.map((a) =>
				a.label === 'Ban'
					? { ...a, label: user.ban ? 'Unban' : 'Ban' }
					: a
			);
		}
		users = originalUsers = usersFromDB;
	}

	async function onUnbanUser() {
		isUnbanning = true;
		try {
			const res = await fetch(`/users/${selectedUser.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ban: false })
			});
			if (!res.ok) {
				console.error('Error: Failed to unban user');
				isUnbanning = false;
				return;
			}
			selectedUser.ban = false;
			originalUsers = originalUsers.map((u) =>
				u.id === selectedUser.id
					? {
							...u,
							ban: false,
							actions: u.actions.map((a: { label: string }) =>
								a.label === 'Ban' || a.label === 'Unban'
									? {
											...a,
											label: 'Ban'
										}
									: a
							)
					  }
					: u
			);
			users = [...originalUsers];
			console.log(`User ${selectedUser.name} has been unbanned successfully ✅`);
		} catch (err) {
			console.error('Network error unbanning user:', err);
		} finally {
			showUnbanModal = false;
			isUnbanning = false;
			unbanReason = '';
		}
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
		<div class="relative mx-auto max-w-3xl">
			<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
				<Search class="h-5 w-5 text-gray-400" />
			</div>
			<input
				type="text"
				placeholder="Search Username"
				class="w-full rounded-full border border-gray-200 bg-white py-3 pl-12 pr-4 text-gray-900 placeholder-gray-500 shadow-sm outline-none transition-all duration-200 focus:border-transparent focus:ring-1 focus:ring-gray-300"
				bind:value={searchQuery}
				oninput={onUserSearch}
			/>
		</div>

		<!-- Filter Pills -->
		<div class="flex flex-wrap items-center justify-center gap-3">
			<div class="flex items-center gap-2">
				<span class="text-sm font-medium text-gray-600">Filters:</span>

				<FilterPill
					label="Role"
					options={VALID_ROLES.map((role) => ({ value: role, label: role }))}
					bind:selectedValue={currentFilteredRole}
					onSelectionChange={onFilteringRole}
				/>

				<FilterPill
					label="Verified"
					options={VALID_VERIFICATION_OPTIONS.map((opt: boolean) => ({
						value: opt.toString(),
						label: opt.toString()
					}))}
					bind:selectedValue={currentFilteredVerificationStatus}
					onSelectionChange={onFilteringVerificationStatus}
				/>

				<FilterPill
					label="Banned"
					options={[
						{ value: 'true', label: 'true' },
						{ value: 'false', label: 'false' }
					]}
					bind:selectedValue={currentFilteredBannedStatus}
					onSelectionChange={onFilteringBannedStatus}
				/>

				<!-- sort -->

				<!-- <div class="h-4 w-px bg-gray-300"></div> -->

				<FilterPill
					label="Sort"
					options={SORT_OPTIONS}
					bind:selectedValue={sortBy}
					onSelectionChange={onUserSort}
					type="sort"
				/>
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
	actionName="Delete"
	actOnKind="User"
	actOnIndividual={selectedUser === null ? '' : selectedUser.name}
	bind:isActionInProgress={isDeleting}
	bind:reasonForAction={deleteReason}
	action={onDeleteUser}
/>

<ConfirmActionWithReason
	bind:isVisible={showBanModal}
	actionName="Ban"
	actOnKind="User"
	actOnIndividual={selectedUser === null ? '' : selectedUser.name}
	bind:isActionInProgress={isBanning}
	bind:reasonForAction={banReason}
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

<ConfirmActionWithReason
	bind:isVisible={showUnbanModal}
	actionName="Unban"
	actOnKind="User"
	actOnIndividual={selectedUser === null ? '' : selectedUser.name}
	bind:isActionInProgress={isUnbanning}
	bind:reasonForAction={unbanReason}
	action={onUnbanUser}
/>
