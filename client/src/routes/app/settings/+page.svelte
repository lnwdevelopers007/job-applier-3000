<script lang="ts">
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { userService } from '$lib/services/userService';
	import { onMount } from 'svelte';
	
	let activeTab = $state('user');
	// Combined user and profile data for ProfileSettings component
	let userData = $state({
		// User account data
		id: '',
		name: '',
		email: '',
		avatar: '',
		provider: '',
		userID: '',
		role: '',
		verified: false,
		currentPassword: '',
		newPassword: '',
		confirmPassword: '',
		googleConnected: false,
		// Profile data
		fullName: '',
		location: '',
		phone: '',
		linkedin: '',
		desiredRole: '',
		aboutMe: '',
		dateOfBirth: '',
		gender: '',
		portfolio: '',
		github: '',
		// Documents
		documents: [] as any[]
	});
	let loading = $state(true);
	let error = $state(null);
	let successMessage = $state('');
	
	const tabs = [
		{
			id: 'user',
			label: 'User',
			title: 'User info',
			description: 'Update your photo and user details here.'
		},
		{
			id: 'personal',
			label: 'Personal Info',
			title: 'Personal Info',
			description: 'Manage your personal information that will be shared with employers.'
		},
		{
			id: 'documents',
			label: 'Documents',
			title: 'Documents',
			description: 'Upload and manage your documents and files.'
		}
	];
	
	onMount(async () => {
		await loadUserData();
	});
	
	async function loadUserData() {
		try {
			loading = true;
			error = null;
			
			const user = await userService.getCurrentUser();
			
			// Combine all data into userData
			userData = {
				// User account data
				id: user.id || user._id,
				name: user.name || '',
				email: user.email || '',
				avatar: user.avatarURL || '',
				provider: user.provider || '',
				userID: user.userID || '',
				role: user.role || 'jobSeeker',
				verified: user.verified || false,
				currentPassword: '',
				newPassword: '',
				confirmPassword: '',
				googleConnected: user.provider === 'google',
				// Profile data - default empty strings
				fullName: '',
				location: '',
				phone: '',
				linkedin: '',
				desiredRole: '',
				aboutMe: '',
				dateOfBirth: '',
				gender: '',
				portfolio: '',
				github: '',
				documents: []
			};
			
			// Merge userInfo if it exists
			if (user.userInfo) {
				const info = user.userInfo as any;
				userData.fullName = info.fullName || '';
				userData.location = info.location || '';
				userData.phone = info.phone || '';
				userData.linkedin = info.linkedIn || '';
				userData.desiredRole = info.desiredRole || '';
				userData.aboutMe = info.aboutMe || '';
				userData.dateOfBirth = info.dateOfBirth || '';
				userData.gender = info.gender || '';
				userData.portfolio = info.portfolio || '';
				userData.github = info.github || '';
			}
		} catch (err: any) {
			console.error('Failed to load user data:', err);
			error = 'Failed to load profile data. Please try again.';
			
			// Fallback to mock data for development
			if (import.meta.env.DEV) {
				const { mockSeekerData } = await import('$lib/stores/mockData.js');
				const mockData = { ...mockSeekerData };
				userData.name = mockData.fullName || '';
				userData.email = mockData.email || '';
				userData.avatar = mockData.avatar || '';
				userData.googleConnected = mockData.googleConnected || false;
				userData.fullName = mockData.fullName || '';
				userData.location = mockData.location || '';
				userData.phone = mockData.phone || '';
				userData.documents = mockData.documents || [];
			}
		} finally {
			loading = false;
		}
	}
	
	async function handleSave(data: any, changedFields?: string[]) {
		try {
			error = null;
			successMessage = '';
			
			// Get user ID from stored data
			const userId = userData.id;
			if (!userId) {
				throw new Error('User ID not found');
			}
			
			// Transform data to backend format with optional changed fields
			const payload = userService.transformToBackendFormat(data, 'seeker', changedFields);
			
			// Handle password update separately
			if (data.newPassword && data.currentPassword) {
				await userService.updatePassword(data.currentPassword, data.newPassword);
				// Clear password fields
				userData.currentPassword = '';
				userData.newPassword = '';
				userData.confirmPassword = '';
			}
			
			// Update user profile
			const updatedUser = await userService.updateUser(userId, payload);
			
			// Update userData with the response
			if (updatedUser.name) userData.name = updatedUser.name;
			if (updatedUser.avatarURL) userData.avatar = updatedUser.avatarURL;
			
			if (updatedUser.userInfo) {
				const info = updatedUser.userInfo as any;
				if (info.fullName) userData.fullName = info.fullName;
				if (info.location) userData.location = info.location;
				if (info.phone) userData.phone = info.phone;
				if (info.linkedIn) userData.linkedin = info.linkedIn;
				// Update other profile fields as needed
			}
			
			successMessage = 'Profile saved successfully!';
			setTimeout(() => successMessage = '', 3000);
		} catch (err: any) {
			console.error('Failed to save profile:', err);
			error = err.message || 'Failed to save profile. Please try again.';
		}
	}
	
	async function handleDocumentSave(data: any) {
		try {
			error = null;
			successMessage = '';
			
			// Upload new documents with files
			const documentsToUpload = data.documents?.filter((doc: any) => doc.file) || [];
			for (const doc of documentsToUpload) {
				try {
					const uploadedDoc = await userService.uploadDocument(doc.file);
					// Update document with server response
					doc.id = uploadedDoc.id;
					doc.url = uploadedDoc.url;
					delete doc.file;
				} catch (uploadErr) {
					console.error('Failed to upload document:', uploadErr);
				}
			}
			
			// Update user with new documents list
			const userId = userData.id;
			if (userId) {
				const payload = { documents: data.documents };
				await userService.updateUser(userId, payload);
			}
			
			userData.documents = data.documents;
			successMessage = 'Documents saved successfully!';
			setTimeout(() => successMessage = '', 3000);
		} catch (err: any) {
			console.error('Failed to save documents:', err);
			error = err.message || 'Failed to save documents. Please try again.';
		}
	}
</script>

{#if loading}
	<div class="flex items-center justify-center min-h-[400px]">
		<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
	</div>
{:else}
	{#if error}
		<div class="mb-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
			{error}
		</div>
	{/if}
	
	{#if successMessage}
		<div class="mb-4 p-4 bg-green-50 border border-green-200 text-green-700 rounded-lg">
			{successMessage}
		</div>
	{/if}
	
	<ProfileSettings 
		{tabs}
		bind:activeTab
		bind:userData={userData}
		onSave={activeTab === 'documents' ? handleDocumentSave : handleSave}
		userType="seeker"
		title="Profile Settings"
	/>
{/if}