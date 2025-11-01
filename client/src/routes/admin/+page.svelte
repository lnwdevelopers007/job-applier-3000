<script lang="ts">
	import { fetchUsers } from '$lib/utils/fetcher';

	interface UserWithActions {
		id: string;
		name: string;
		email: string;
		role: string;
		verified: boolean;
		actions: Array<any>;
	}

	const USER_ACTIONS = [
		{ label: 'View', disabled: false },
		{ label: 'Edit', disabled: false },
		{ label: 'Delete', disabled: false }
	];

	let selectedUser = $state<number | null>(null);

	function selectRow(index: number) {
		selectedUser = selectedUser === index ? null : index;
	}

	function handleUserAction(arg: any, anotherArg: any) {
		return null;
	}
	let users = $state<UserWithActions[]>([]);

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

<!-- Jobs Table -->
<div class="mt-6 overflow-x-auto rounded-lg bg-white shadow">
	<table class="min-w-full text-left text-sm text-gray-600">
		<thead class="bg-gray-100 text-sm font-semibold text-gray-700">
			<tr>
				<th class="px-4 py-3">Username</th>
				<th class="px-4 py-3">Email</th>
				<th class="px-4 py-3">Role</th>
				<th class="px-4 py-3">Verification</th>
				<th class="px-4 py-3">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each users as user, i (i)}
				<tr
					class="cursor-pointer border-t {selectedUser === i ? 'bg-green-200' : ''}"
					onclick={() => selectRow(i)}
				>
					<td class="px-4 py-3 font-medium text-gray-900">{user.name}</td>
					<td class="px-4 py-3">{user.email}</td>
					<td class="px-4 py-3">{user.role}</td>
					<td class="px-4 py-3">{user.verified}</td>

					<td class="space-x-2 px-4 py-3">
						{#each user.actions as action, j (j)}
							<button
								class="border-1 rounded border-gray-500 bg-gray-100 px-3 py-1 text-sm text-gray-900 enabled:hover:bg-gray-300 disabled:cursor-not-allowed disabled:opacity-50"
								disabled={action.disabled}
								onclick={(e) => {
									e.stopPropagation();
									handleUserAction(action, user);
								}}
							>
								{action.label}
							</button>
						{/each}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
