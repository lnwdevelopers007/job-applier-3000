<script lang="ts">
	import UserTab from './tabs/UserTab.svelte';
	import PersonalInfoTab from './tabs/PersonalInfoTab.svelte';
	import CompanyTab from './tabs/CompanyTab.svelte';
	import DocumentsTab from './tabs/DocumentsTab.svelte';
	import { userService } from '$lib/services/userService';
	
	interface Tab {
		id: string;
		label: string;
		title: string;
		description: string;
	}
	
	interface UserData {
		// User account data
		id?: string;
		name?: string;
		email?: string;
		avatar?: string;
		provider?: string;
		userID?: string;
		role?: string;
		verified?: boolean;
		currentPassword?: string;
		newPassword?: string;
		confirmPassword?: string;
		googleConnected?: boolean;
		// Profile data (seeker)
		fullName?: string;
		location?: string;
		phone?: string;
		linkedin?: string;
		desiredRole?: string;
		aboutMe?: string;
		dateOfBirth?: string;
		gender?: string;
		portfolio?: string;
		github?: string;
		// Company data
		companyName?: string;
		aboutCompany?: string;
		industry?: string;
		companySize?: string;
		companyWebsite?: string;
		companyLogo?: string;
		foundedYear?: string;
		headquarters?: string;
		companyLinkedin?: string;
		// Documents
		documents?: any[];
		// File uploads
		avatarFile?: File;
		companyLogoFile?: File;
	}
	
	let {
		tabs = [] as Tab[],
		activeTab = $bindable(''),
		userData = $bindable({} as UserData),
		onSave = (data: UserData, changedFields?: string[]) => Promise.resolve(),
		title = 'Settings',
		userType = 'seeker' as 'seeker' | 'company'
	} = $props();
	
	// Track if tab has unsaved changes
	let hasChanges = $state(false);
	let isSaving = $state(false);
	let initialData = $state({} as UserData);
	let isInitialized = $state(false);
	let changedFields = $state(new Set<string>());
	
	// Handle tab-level save - different logic for each tab
	async function handleTabSave() {
		isSaving = true;
		try {
			if (activeTab === 'user') {
				// User tab - save account-level fields
				await handleUserTabSave();
			} else if (activeTab === 'personal' || activeTab === 'company') {
				// Profile tabs - save through transform function with changed fields
				await onSave(userData, Array.from(changedFields));
			} else if (activeTab === 'documents') {
				// Documents tab - handle document operations
				await onSave(userData);
			}
			
			hasChanges = false;
			// Clear changed fields after successful save
			changedFields.clear();
			// Update initial data after successful save
			initialData = JSON.parse(JSON.stringify(userData));
		} catch (error) {
			console.error('Save failed:', error);
		} finally {
			isSaving = false;
		}
	}
	
	// Handle user tab save - account level fields
	async function handleUserTabSave() {
		// Send complete user object to avoid losing any data
		const payload: any = {
			name: userData.name || '',
			email: userData.email || '',
			avatarURL: userData.avatar || '',
			provider: userData.provider || '',
			userID: userData.userID || '',
			role: userData.role || 'jobSeeker',
			verified: userData.verified || false,
			userInfo: {}
		};
		
		// Include complete userInfo to preserve it
		if (userType === 'seeker') {
			payload.userInfo = {
				fullName: userData.fullName || '',
				location: userData.location || '',
				phone: userData.phone || '',
				linkedIn: userData.linkedin || '',
				desiredRole: userData.desiredRole || '',
				aboutMe: userData.aboutMe || '',
				dateOfBirth: userData.dateOfBirth || '',
				gender: userData.gender || '',
				portfolio: userData.portfolio || '',
				github: userData.github || ''
			};
		} else if (userType === 'company') {
			payload.userInfo = {
				name: userData.companyName || '',
				aboutUs: userData.aboutCompany || '',
				industry: userData.industry || '',
				size: userData.companySize || '',
				website: userData.companyWebsite || '',
				logo: userData.companyLogo || '',
				foundedYear: userData.foundedYear || '',
				headquarters: userData.headquarters || '',
				linkedIn: userData.companyLinkedin || ''
			};
		}
		
		// Call userService directly for user account updates
		await userService.updateUser(userData.id, payload);
	}
	
	
	// Handle direct saves (like disconnect) that bypass the transform function
	function handleDirectSave() {
		hasChanges = false;
		// Clear changed fields
		changedFields.clear();
		// Update initial data to reflect current userData
		initialData = JSON.parse(JSON.stringify(userData));
	}
	
	// Track changes to userData (only after initialization, excluding password fields)
	$effect(() => {
		if (isInitialized) {
			// Create copies without password fields for comparison
			const currentData = { ...userData };
			const initData = { ...initialData };
			
			// Remove password fields from comparison
			const currentAny = currentData as any;
			const initAny = initData as any;
			if (currentAny.currentPassword !== undefined) delete currentAny.currentPassword;
			if (currentAny.newPassword !== undefined) delete currentAny.newPassword; 
			if (currentAny.confirmPassword !== undefined) delete currentAny.confirmPassword;
			if (initAny.currentPassword !== undefined) delete initAny.currentPassword;
			if (initAny.newPassword !== undefined) delete initAny.newPassword;
			if (initAny.confirmPassword !== undefined) delete initAny.confirmPassword;
			
			// Track specific field changes
			const newChangedFields = new Set<string>();
			
			// List of all trackable fields
			const trackableFields = [
				'name', 'email', 'avatar', 'role', 'verified',
				'fullName', 'location', 'phone', 'linkedin', 'desiredRole', 'aboutMe', 
				'dateOfBirth', 'gender', 'portfolio', 'github',
				'companyName', 'aboutCompany', 'industry', 'companySize', 'companyWebsite',
				'companyLogo', 'foundedYear', 'headquarters', 'companyLinkedin'
			];
			
			for (const field of trackableFields) {
				if (currentAny[field] !== initAny[field]) {
					newChangedFields.add(field);
				}
			}
			
			changedFields = newChangedFields;
			hasChanges = newChangedFields.size > 0;
			
			// Debug: Log changed fields (only in development)
			if (import.meta.env.DEV && newChangedFields.size > 0) {
				console.log('Changed fields:', Array.from(newChangedFields));
			}
		}
	});
	
	// Initialize nested objects and arrays if they don't exist
	$effect(() => {
		if (!userData.documents) userData.documents = [];
		
		// Initialize profile fields to prevent undefined binding errors
		if (userData.fullName === undefined) userData.fullName = '';
		if (userData.location === undefined) userData.location = '';
		if (userData.desiredRole === undefined) userData.desiredRole = '';
		if (userData.aboutMe === undefined) userData.aboutMe = '';
		if (userData.phone === undefined) userData.phone = '';
		if (userData.dateOfBirth === undefined) userData.dateOfBirth = '';
		if (userData.gender === undefined) userData.gender = '';
		if (userData.linkedin === undefined) userData.linkedin = '';
		if (userData.portfolio === undefined) userData.portfolio = '';
		if (userData.github === undefined) userData.github = '';
		
		// Initialize user fields
		if (userData.name === undefined) userData.name = '';
		if (userData.email === undefined) userData.email = '';
		
		// Set googleConnected based on provider field
		userData.googleConnected = userData.provider === 'google';
		
		// After initialization, set initial data and mark as initialized
		if (!isInitialized) {
			initialData = JSON.parse(JSON.stringify(userData));
			isInitialized = true;
		}
	});
	
	function handleTabChange(tabId: string) {
		activeTab = tabId;
	}
	
	const activeTabData = $derived(tabs.find((tab: Tab) => tab.id === activeTab));
</script>

<div class="bg-white rounded-lg border border-gray-200 p-6 min-h-[1000px]">
	<h1 class="text-2xl font-semibold text-gray-900 p-3">{title}</h1>
	
	<div class="border-b border-gray-200">
		<nav class="flex -mb-px">
			{#each tabs as tab}
				<button
					onclick={() => handleTabChange(tab.id)}
					class="px-4 py-3 text-sm font-medium border-b-2 transition-all cursor-pointer {activeTab === tab.id ? 'border-green-600 text-green-700' : 'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300'}"
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>
	
	<div class="flex justify-between items-start px-8 py-6 border-b border-gray-200">
		<div>
			<h2 class="text-lg font-medium text-gray-900">{activeTabData?.title || ''}</h2>
			<p class="text-sm text-gray-500 mt-1">{activeTabData?.description || ''}</p>
		</div>
		{#if activeTab === 'documents'}
			<button
				onclick={() => document.getElementById('file-upload')?.click()}
				class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors hover:cursor-pointer"
			>
				Add Documents
			</button>
		{:else if hasChanges}
			<div class="flex items-center gap-3">
				<button
					onclick={handleTabSave}
					disabled={isSaving}
					class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:opacity-50 transition-colors hover:cursor-pointer"
				>
					{isSaving ? 'Saving...' : 'Save Changes'}
				</button>
				<button
					onclick={() => {
						// Revert userData back to initial state
						userData = JSON.parse(JSON.stringify(initialData));
						hasChanges = false;
						// Clear changed fields
						changedFields.clear();
					}}
					disabled={isSaving}
					class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors hover:cursor-pointer"
				>
					Cancel
				</button>
			</div>
		{/if}
	</div>
	
	<div class="tab-content">
		{#if activeTab === 'user'}
			<UserTab bind:userData={userData} {onSave} onDirectSave={handleDirectSave} />
		{:else if activeTab === 'personal' && userType === 'seeker'}
			<PersonalInfoTab bind:profileData={userData} />
		{:else if activeTab === 'company' && userType === 'company'}
			<CompanyTab bind:companyData={userData} />
		{:else if activeTab === 'documents'}
			<DocumentsTab bind:documents={userData.documents} />
		{/if}
	</div>
</div>