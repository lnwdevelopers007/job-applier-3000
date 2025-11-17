<script lang="ts">
	import Modal from './Modal.svelte';
	import { Loader2 } from 'lucide-svelte';
	import type { UserDisplay } from '../tables/columns/userColumns';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onConfirm: (role: string, verified: boolean) => Promise<void>;
		user: UserDisplay | null;
		isUpdating?: boolean;
	}

	let { 
		isOpen = $bindable(false), 
		onClose, 
		onConfirm, 
		user,
		isUpdating = false
	}: Props = $props();

	let selectedRole = $state('');
	let selectedVerified = $state(false);

	const roleOptions = [
		{ value: 'jobSeeker', label: 'Job Seeker' },
		{ value: 'company', label: 'Company' },
		{ value: 'faculty', label: 'Faculty' },
		{ value: 'admin', label: 'Admin' }
	];

	$effect(() => {
		if (user && isOpen) {
			selectedRole = user.role;
			selectedVerified = user.verified;
		}
	});

	async function handleConfirm() {
		if (!user) return;
		await onConfirm(selectedRole, selectedVerified);
	}

	function handleClose() {
		if (!isUpdating) {
			onClose();
		}
	}

	const hasChanges = $derived(
		user && (selectedRole !== user.role || selectedVerified !== user.verified)
	);
</script>

<Modal bind:isOpen={isOpen} size="lg" closeOnBackdrop={false} onClose={handleClose}>
	{#if user}
		<div class="p-6">
			<!-- Header -->
			<div class="flex items-start space-x-4 mb-6">
				<div class="flex-1 min-w-0">
					<h3 class="text-lg font-medium text-gray-900 mb-2">
						Edit Permissions
					</h3>
					<p class="text-sm text-gray-600">
						Editing permissions for <span class="font-medium text-gray-900">"{user.name}"</span>. Changes will be applied immediately.
					</p>
				</div>
			</div>

			<!-- Form -->
			<div class="space-y-6 mb-6">
				<div>
					<label class="block text-sm font-medium text-gray-700 mb-3">
						User Role <span class="text-red-500">*</span>
					</label>
					<select
						bind:value={selectedRole}
						disabled={isUpdating}
						class="w-full border border-gray-300 rounded-lg text-sm px-3 py-3 text-gray-800 focus:outline-none focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all disabled:bg-gray-100 disabled:cursor-not-allowed"
					>
						{#each roleOptions as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 mb-3">
						Verification Status <span class="text-red-500">*</span>
					</label>
					<div class="flex items-center space-x-6">
						<label class="flex items-center">
							<input
								type="radio"
								bind:group={selectedVerified}
								value={true}
								disabled={isUpdating}
								class="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300"
							/>
							<span class="ml-2 text-sm text-gray-700">Verified</span>
						</label>
						<label class="flex items-center">
							<input
								type="radio"
								bind:group={selectedVerified}
								value={false}
								disabled={isUpdating}
								class="h-4 w-4 text-green-600 focus:ring-green-500 border-gray-300"
							/>
							<span class="ml-2 text-sm text-gray-700">Unverified</span>
						</label>
					</div>
				</div>
			</div>

			<!-- Actions -->
			<div class="flex justify-end space-x-3">
				<button
					class="px-4 py-2 bg-white border border-gray-300 text-gray-700 text-sm rounded-lg hover:bg-gray-50 font-medium transition-colors"
					onclick={handleClose}
					disabled={isUpdating}
				>
					Cancel
				</button>
				<button
					class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 text-sm disabled:opacity-50 disabled:cursor-not-allowed font-medium transition-colors flex items-center"
					disabled={!hasChanges || isUpdating}
					onclick={handleConfirm}
				>
					{#if isUpdating}
						<Loader2 class="animate-spin mr-2 h-4 w-4" />
						Updating...
					{:else}
						Update Permissions
					{/if}
				</button>
			</div>
		</div>
	{/if}
</Modal>