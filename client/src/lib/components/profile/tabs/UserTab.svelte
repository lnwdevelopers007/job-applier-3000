<script>
	import TextInput from '$lib/components/forms/TextInput.svelte';
	import FileUpload from '$lib/components/forms/FileUpload.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
	import { userService } from '$lib/services/userService';
	import { Mail } from 'lucide-svelte';
	
	let {
		userData = $bindable({}),
		onSave = () => {},
		onDirectSave = () => {} // Callback for when we save directly without going through transform
	} = $props();
	
	// Separate password change handling
	let passwordData = $state({
		currentPassword: '',
		newPassword: '',
		confirmPassword: ''
	});
	let isPasswordSaving = $state(false);
	let passwordError = $state('');
	let showPasswordForm = $state(false);
	
	// Google disconnect modal
	let showDisconnectModal = $state(false);
	let isDisconnecting = $state(false);
	
	async function handlePasswordChange() {
		passwordError = '';
		
		if (!passwordData.currentPassword || !passwordData.newPassword || !passwordData.confirmPassword) {
			passwordError = 'All password fields are required';
			return;
		}
		
		if (passwordData.newPassword !== passwordData.confirmPassword) {
			passwordError = 'New passwords do not match';
			return;
		}
		
		if (passwordData.newPassword.length < 8) {
			passwordError = 'New password must be at least 8 characters';
			return;
		}
		
		isPasswordSaving = true;
		try {
			// Send only password data
			await onSave({
				currentPassword: passwordData.currentPassword,
				newPassword: passwordData.newPassword
			});
			
			// Clear password fields on success and hide form
			passwordData = {
				currentPassword: '',
				newPassword: '',
				confirmPassword: ''
			};
			showPasswordForm = false;
		} catch (error) {
			passwordError = error.message || 'Failed to update password';
		} finally {
			isPasswordSaving = false;
		}
	}
	
	// Check if user is actually connected to Google based on provider field
	const isGoogleConnected = $derived(userData.provider === 'google');
	
	// Google disconnect handling
	async function handleGoogleDisconnect() {
		isDisconnecting = true;
		try {
			// Call API directly to disconnect Google OAuth account
			await userService.updateUser(userData.id, { provider: "" });
			
			// Update userData after successful API call
			userData.provider = '';
			userData.userID = '';
			
			// Notify parent that we saved directly (to update initialData)
			onDirectSave();
			
			showDisconnectModal = false;
		} catch (error) {
			console.error('Failed to disconnect Google account:', error);
		} finally {
			isDisconnecting = false;
		}
	}
</script>

<div class="divide-y divide-gray-200">
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Username</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={userData.name}
					placeholder="Display name"
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Email address</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={userData.email}
					type="email"
					readonly={true}
					iconComponent={Mail}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Profile Picture</label>
				<p class="text-xs text-gray-500 mt-1">This will be displayed on your profile.</p>
			</div>
			<div class="col-span-2">
				<FileUpload 
					currentImage={userData.avatar}
					onFileSelect={(file) => {
						userData.avatarFile = file;
					}}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Connected accounts</label>
				<p class="text-xs text-gray-500 mt-1">Manage your connected social accounts</p>
			</div>
			<div class="col-span-1">
				<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 p-4 bg-white border border-gray-200 rounded-lg">
					<div class="flex items-center gap-3 min-w-0">
						<div class="w-10 h-10 bg-white border border-gray-200 rounded-full flex items-center justify-center flex-shrink-0">
							<svg class="w-5 h-5" viewBox="0 0 24 24">
								<path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
								<path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
								<path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
								<path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
							</svg>
						</div>
						<div class="min-w-0">
							<p class="text-sm font-medium text-gray-900">Google</p>
							<p class="text-xs text-gray-500 truncate">{userData.email}</p>
						</div>
					</div>
					<button 
						onclick={() => {
							if (isGoogleConnected) {
								showDisconnectModal = true;
							} else {
								window.location.href = `${import.meta.env.VITE_BACKEND}/auth/google?role=login`;
							}
						}}
						class="px-3 py-1.5 text-xs font-medium border rounded-md transition-colors cursor-pointer flex-shrink-0 {isGoogleConnected ? 'text-red-600 border-red-200 hover:bg-red-50' : 'text-green-600 border-green-200 hover:bg-green-50'}"
					>
						{isGoogleConnected ? 'Disconnect' : 'Connect'}
					</button>
				</div>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Password</label>
				<p class="text-xs text-gray-500 mt-1">Update your login password</p>
			</div>
			<div class="col-span-2">
				{#if !showPasswordForm}
					<button
						onclick={() => showPasswordForm = true}
						class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 transition-colors hover:cursor-pointer"
					>
						Change Password
					</button>
				{:else}
					<div class="max-w-md space-y-4">
						<PasswordInput
							id="current-password"
							label="Current password"
							placeholder="Enter current password"
							bind:value={passwordData.currentPassword}
						/>
						
						<div>
							<PasswordInput
								id="new-password"
								label="New password"
								placeholder="Enter new password"
								bind:value={passwordData.newPassword}
							/>
							<p class="text-xs text-gray-500 mt-1">Must be at least 8 characters.</p>
						</div>
						
						<PasswordInput
							id="confirm-password"
							label="Confirm password"
							placeholder="Confirm new password"
							bind:value={passwordData.confirmPassword}
						/>
						
						{#if passwordError}
							<p class="text-xs text-red-600">{passwordError}</p>
						{/if}
						
						<div class="flex items-center gap-3 pt-2">
							<div class="relative group">
								<button
								onclick={handlePasswordChange}
								disabled={true}
								class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md disabled:opacity-50 transition-colors hover:cursor-not-allowed"
								>
								{isPasswordSaving ? 'Changing...' : 'Change Password'}
								</button>
								<div
								class="absolute bottom-full mb-2 left-1/2 -translate-x-1/2 hidden group-hover:block
										bg-gray-800 text-white text-xs rounded py-1 px-2 whitespace-nowrap z-10"
								>
								Login with password is not implemented yet
								</div>
							</div>

							<button
								onclick={() => {
								showPasswordForm = false;
								passwordData = { currentPassword: '', newPassword: '', confirmPassword: '' };
								passwordError = '';
								}}
								disabled={isPasswordSaving}
								class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors hover:cursor-pointer"
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

<Modal bind:isOpen={showDisconnectModal} size="md" closeOnEscape={!isDisconnecting}>
	<div class="p-6">
		<h3 class="text-lg font-medium text-gray-900 mb-4">Disconnect Google Account</h3>
		<p class="text-sm text-gray-600 mb-6">
			Are you sure you want to disconnect your Google account? You'll need to set up a password to continue accessing your account.
		</p>
		<div class="flex items-center justify-end gap-3">
			<button
				onclick={() => showDisconnectModal = false}
				disabled={isDisconnecting}
				class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors hover:cursor-pointer"
			>
				Cancel
			</button>
			<button
				onclick={handleGoogleDisconnect}
				disabled={isDisconnecting}
				class="px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-md hover:bg-red-700 disabled:opacity-50 transition-colors hover:cursor-pointer"
			>
				{isDisconnecting ? 'Disconnecting...' : 'Disconnect'}
			</button>
		</div>
	</div>
</Modal>