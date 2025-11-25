<script lang="ts">
	import { authStore } from '$lib/stores/auth.svelte';
	import { fetchUsers } from '$lib/utils/fetcher';
	import { Search } from 'lucide-svelte';
	import DeleteModal from '$lib/components/ui/DeleteModal.svelte';
	import PermissionEditModal from '$lib/components/ui/PermissionEditModal.svelte';
	import DataTable from '$lib/components/tables/DataTable.svelte';
	import { createUserColumns, type UserDisplay } from '$lib/components/tables/columns/userColumns';

	let users = $state<UserDisplay[]>([]);
	let loading = $state(true);
	let searchQuery = $state('');
	let roleFilter = $state('all');
	let verificationFilter = $state('all');

	let showDeleteModal = $state(false);
	let userToDelete = $state<UserDisplay | null>(null);
	let deleting = $state(false);

	let showBanModal = $state(false);
	let userToBan = $state<UserDisplay | null>(null);
	let banning = $state(false);

	let showPermissionModal = $state(false);
	let userToEditPermissions = $state<UserDisplay | null>(null);
	let updatingPermissions = $state(false);

	const roleOptions = [
		{ value: 'all', label: 'All Roles' },
		{ value: 'jobSeeker', label: 'Job Seeker' },
		{ value: 'company', label: 'Company' },
		{ value: 'faculty', label: 'Faculty' },
		{ value: 'admin', label: 'Admin' }
	];

	const verificationOptions = [
		{ value: 'all', label: 'All Status' },
		{ value: 'verified', label: 'Verified' },
		{ value: 'unverified', label: 'Unverified' }
	];

	function handleEdit(user: UserDisplay) {
		userToEditPermissions = user;
		showPermissionModal = true;
	}

	function handleBan(user: UserDisplay) {
		userToBan = user;
		showBanModal = true;
	}

	function handleDelete(user: UserDisplay) {
		userToDelete = user;
		showDeleteModal = true;
	}

	async function confirmDelete() {
		if (!userToDelete) return;
		deleting = true;
		try {
			const res = await fetch(`/users/${userToDelete.id}`, {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
				// body: JSON.stringify({ reason: deleteReason })
			});

			if (!res.ok) {
				throw new Error('Failed to delete user');
			}

			users = users.filter((u) => u.id !== userToDelete?.id);
			showDeleteModal = false;
			userToDelete = null;
		} catch (err) {
			console.error('Error deleting user:', err);
			throw err;
		} finally {
			deleting = false;
		}
	}

	async function confirmBan() {
		if (!userToBan) return;
		banning = true;
		try {
			const newBanStatus = !userToBan.banned;
			const res = await fetch(`/users/${userToBan.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({ banned: newBanStatus })
			});

			if (!res.ok) {
				throw new Error('Failed to update ban status');
			}

			// Update user in the list
			users = users.map((u) => (u.id === userToBan!.id ? { ...u, banned: newBanStatus } : u));

			showBanModal = false;
			userToBan = null;
		} catch (err) {
			console.error('Error updating ban status:', err);
			throw err;
		} finally {
			banning = false;
		}
	}

	async function confirmPermissionUpdate(role: string, verified: boolean) {
		if (!userToEditPermissions) return;
		updatingPermissions = true;
		try {
			// Update user role if changed
			if (role !== userToEditPermissions.role) {
				const resRole = await fetch(`/users/${userToEditPermissions.id}/role`, {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					credentials: 'include',
					body: JSON.stringify({ role })
				});
				if (!resRole.ok) throw new Error("Can't update role");
			}

			// Update verified status if changed
			if (verified !== userToEditPermissions.verified) {
				const resVerified = await fetch(`/users/${userToEditPermissions.id}/verify`, {
					method: 'PATCH',
					headers: { 'Content-Type': 'application/json' },
					credentials: 'include',
					body: JSON.stringify({ verified })
				});
				if (!resVerified.ok) throw new Error("Can't update verified status");
			}

			// Update user in the list
			users = users.map((u) => (u.id === userToEditPermissions!.id ? { ...u, role, verified } : u));

			showPermissionModal = false;
			userToEditPermissions = null;
		} catch (err) {
			console.error('Error updating permissions:', err);
			throw err;
		} finally {
			updatingPermissions = false;
		}
	}

	const filteredUsers = $derived(
		users.filter((user) => {
			const matchesSearch =
				!searchQuery ||
				user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				user.email.toLowerCase().includes(searchQuery.toLowerCase());

			const matchesRole = roleFilter === 'all' || user.role === roleFilter;

			const matchesVerification =
				verificationFilter === 'all' ||
				(verificationFilter === 'verified' && user.verified) ||
				(verificationFilter === 'unverified' && !user.verified);

			return matchesSearch && matchesRole && matchesVerification;
		})
	);

	const columns = $derived(
		createUserColumns({
			onEdit: handleEdit,
			onBan: handleBan,
			onDelete: handleDelete,
			getCurrentUserId: () => authStore.user?.userID
		})
	);

	$effect(() => {
		if (authStore.isAuthenticated && authStore.user) {
			loadAllUsers();
		}
	});

	async function loadAllUsers() {
		try {
			loading = true;
			const usersData = await fetchUsers();

			// Fetch files for each user
			const usersWithFiles = await Promise.allSettled(
				usersData.map(async (user: any) => {
					let userFiles: string[] = [];
					try {
						const filesResponse = await fetch(`/files/user/${user.id}`, {
							credentials: 'include'
						});

						if (filesResponse.ok) {
							const filesData = await filesResponse.json();
							// Handle both old format (array) and new format (object with files array)
							const files = Array.isArray(filesData) ? filesData : filesData.files || [];
							userFiles = files.map((file: any) => ({
								id: file.id,
								filename: file.filename
							}));
						}
					} catch (fileErr) {
						console.warn(`Failed to load files for user ${user.id}:`, fileErr);
						userFiles = [];
					}

					return {
						...user,
						banned: user.banned || false,
						files: userFiles
					};
				})
			);

			// Filter successful results and map to user format
			users = usersWithFiles
				.filter((result): result is PromiseFulfilledResult<any> => result.status === 'fulfilled')
				.map((result) => result.value);
		} catch (err) {
			console.error('Error loading users:', err);
			users = [];
		} finally {
			loading = false;
		}
	}

	function handleSearch(event: Event) {
		searchQuery = (event.target as HTMLInputElement).value;
	}

	function handleRoleFilter(event: Event) {
		roleFilter = (event.target as HTMLSelectElement).value;
	}

	function handleVerificationFilter(event: Event) {
		verificationFilter = (event.target as HTMLSelectElement).value;
	}
</script>

<div>
	<div>
		<div class="mb-8">
			<h1 class="mb-1 text-2xl font-semibold text-gray-900">User Management</h1>
			<p class="mb-6 text-base text-gray-600">Manage and monitor all users across your platform</p>
		</div>

		<div class="mb-6">
			<div class="flex gap-3">
				<div class="relative">
					<Search
						class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 transform text-gray-400"
					/>
					<input
						type="text"
						placeholder="Search users by name or email..."
						class="min-w-lg w-full rounded-lg border border-gray-200 py-2.5 pl-10 pr-4 text-sm transition-all placeholder:text-gray-500 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400"
						bind:value={searchQuery}
						oninput={handleSearch}
					/>
				</div>

				<select
					class="cursor-pointer appearance-none rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 pr-10 text-sm font-medium text-gray-700 transition-all hover:bg-gray-100 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400"
					bind:value={roleFilter}
					onchange={handleRoleFilter}
				>
					{#each roleOptions as option (option.value)}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>

				<select
					class="cursor-pointer appearance-none rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 pr-10 text-sm font-medium text-gray-700 transition-all hover:bg-gray-100 focus:bg-white focus:outline-none focus:ring-1 focus:ring-gray-400"
					bind:value={verificationFilter}
					onchange={handleVerificationFilter}
				>
					{#each verificationOptions as option (option.value)}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
			</div>
		</div>

		<DataTable data={filteredUsers} {columns} pageSize={15} {loading} />
	</div>
</div>

<DeleteModal
	bind:isOpen={showDeleteModal}
	onClose={() => {
		showDeleteModal = false;
		userToDelete = null;
	}}
	onConfirm={confirmDelete}
	title="Delete User"
	itemName={userToDelete?.name || ''}
	description="You're about to delete this user account. This action cannot be undone and will remove the user along with all associated data."
	reasonLabel=""
	confirmButtonText="Delete User"
	isDeleting={deleting}
/>

<DeleteModal
	bind:isOpen={showBanModal}
	onClose={() => {
		showBanModal = false;
		userToBan = null;
	}}
	onConfirm={confirmBan}
	title="{userToBan?.banned ? 'Unban' : 'Ban'} User"
	itemName={userToBan?.name || ''}
	description="You're about to {userToBan?.banned
		? 'unban'
		: 'ban'} this user account. {userToBan?.banned
		? 'This will restore their access to the platform.'
		: 'This will prevent them from accessing the platform.'}"
	reasonLabel=""
	confirmButtonText="{userToBan?.banned ? 'Unban' : 'Ban'} User"
	isDeleting={banning}
/>

<PermissionEditModal
	bind:isOpen={showPermissionModal}
	onClose={() => {
		showPermissionModal = false;
		userToEditPermissions = null;
	}}
	onConfirm={confirmPermissionUpdate}
	user={userToEditPermissions}
	isUpdating={updatingPermissions}
/>
