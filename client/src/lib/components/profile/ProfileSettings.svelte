<script lang="ts">
	import UserTab from './tabs/UserTab.svelte';
	import PersonalInfoTab from './tabs/PersonalInfoTab.svelte';
	import CompanyTab from './tabs/CompanyTab.svelte';
	import DocumentsTab from './tabs/DocumentsTab.svelte';
	import SeekerProfilePreviewDrawer from './SeekerProfilePreviewDrawer.svelte';
	import { userService } from '$lib/services/userService';
	import { Plus, LoaderCircle } from 'lucide-svelte';
	import { toast } from 'svelte-french-toast';
	
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
		googleConnected?: boolean;
		// Profile data (seeker)
		fullName?: string;
		location?: string;
		phone?: string;
		linkedin?: string;
		desiredRole?: string;
		aboutMe?: string;
		dateOfBirth?: string;
		portfolio?: string;
		github?: string;
		skills?: string[];
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
		documents?: unknown[];
		// File uploads
		avatarFile?: File;
		companyLogoFile?: File;
	}
	
	type SaveHandler = (data: UserData, changedFields?: string[]) => Promise<void>;
	
	let {
		tabs = [] as Tab[],
		activeTab = $bindable(''),
		userData = $bindable({} as UserData),
		onSave,
		title = 'Settings',
		userType = 'seeker' as 'seeker' | 'company',
		onTabChange
	}: {
		tabs: Tab[];
		activeTab: string;
		userData: UserData;
		onSave: SaveHandler;
		title: string;
		userType: 'seeker' | 'company';
		onTabChange?: (tabId: string) => void;
	} = $props();
	
	// State management
	let hasChanges = $state(false);
	let isSaving = $state(false);
	let initialData = $state({} as UserData);
	let isInitialized = $state(false);
	let changedFields = $state(new Set<string>());
	let isPreviewOpen = $state(false);
	
	// Constants
	const TRACKABLE_FIELDS = [
		'name', 'email', 'avatar', 'role', 'verified',
		'fullName', 'location', 'phone', 'linkedin', 'desiredRole', 'aboutMe', 
		'dateOfBirth', 'portfolio', 'github', 'skills',
		'companyName', 'aboutCompany', 'industry', 'companySize', 'companyWebsite',
		'companyLogo', 'foundedYear', 'headquarters', 'companyLinkedin'
	] as const;
	
	const EXCLUDED_FIELDS = ['currentPassword', 'newPassword', 'confirmPassword'] as const;
	
	// Utility functions
	function createCleanCopy(data: UserData): UserData {
		const copy = { ...data };
		EXCLUDED_FIELDS.forEach(field => delete (copy as Record<string, unknown>)[field]);
		return copy;
	}
	
	function resetChanges(): void {
		userData = JSON.parse(JSON.stringify(initialData));
		hasChanges = false;
		changedFields.clear();
	}
	
	function updateInitialData(): void {
		hasChanges = false;
		changedFields.clear();
		initialData = JSON.parse(JSON.stringify(userData));
	}
	
	function buildUserInfoPayload(): Record<string, string> {
		if (userType === 'seeker') {
			return {
				fullName: userData.fullName || '',
				location: userData.location || '',
				phone: userData.phone || '',
				linkedIn: userData.linkedin || '',
				desiredRole: userData.desiredRole || '',
				aboutMe: userData.aboutMe || '',
				dateOfBirth: userData.dateOfBirth || '',
				portfolio: userData.portfolio || '',
				github: userData.github || '',
				skills: userData.skills ? userData.skills.join(', ') : ''
			};
		}
		
		if (userType === 'company') {
			return {
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
		
		return {};
	}
	
	// Save handlers
	async function saveUserTab(): Promise<void> {
		if (!userData.id) throw new Error('User ID not found');
		
		const payload = {
			name: userData.name || '',
			email: userData.email || '',
			avatarURL: userData.avatar || '',
			provider: userData.provider || '',
			userID: userData.userID || '',
			role: userData.role || 'jobSeeker',
			verified: userData.verified || false,
			userInfo: buildUserInfoPayload()
		};
		
		await userService.updateUser(userData.id, payload);
	}
	
	async function saveProfileTab(): Promise<void> {
		await onSave(userData, Array.from(changedFields));
	}
	
	async function saveDocumentsTab(): Promise<void> {
		await onSave(userData);
	}
	
	// Main save handler
	async function handleSave(): Promise<void> {
		if (isSaving) return;
		
		const savePromise = (async () => {
			isSaving = true;
			try {
				switch (activeTab) {
					case 'user':
						await saveUserTab();
						break;
					case 'personal':
					case 'company':
						await saveProfileTab();
						break;
					case 'documents':
						await saveDocumentsTab();
						break;
					default:
						throw new Error(`Unknown tab: ${activeTab}`);
				}
				updateInitialData();
			} finally {
				isSaving = false;
			}
		})();

		toast.promise(savePromise, {
			loading: 'Saving changes...',
			success: 'Profile saved successfully!',
			error: (err) => {
				// Extract meaningful error message
				if (err instanceof Error) {
					return err.message;
				}
				return 'Failed to save changes. Please try again.';
			}
		});
	}
	
	// Direct save handler for operations like OAuth disconnect
	function handleDirectSave(): void {
		updateInitialData();
	}
	
	// Tab change handler
	function handleTabChange(tabId: string): void {
		activeTab = tabId;
		onTabChange?.(tabId);
	}
	
	// Preview handler
	function handlePreview(): void {
		isPreviewOpen = true;
	}
	
	// Initialize fields to prevent undefined binding errors
	function initializeFields(): void {
		const fieldsToInit = [
			'documents', 'fullName', 'location', 'desiredRole', 'aboutMe', 'phone',
			'dateOfBirth', 'linkedin', 'portfolio', 'github', 'name', 'email', 'skills'
		];
		
		fieldsToInit.forEach(field => {
			if (userData[field as keyof UserData] === undefined) {
				if (field === 'documents' || field === 'skills') {
					(userData as Record<string, unknown>)[field] = [];
				} else {
					(userData as Record<string, unknown>)[field] = '';
				}
			}
		});
		
		// Set derived fields
		userData.googleConnected = userData.provider === 'google';
	}
	
	// Track field changes
	function trackChanges(): void {
		if (!isInitialized) return;
		
		const currentData = createCleanCopy(userData);
		const initData = createCleanCopy(initialData);
		
		const newChangedFields = new Set<string>();
		
		for (const field of TRACKABLE_FIELDS) {
			const currentValue = (currentData as Record<string, unknown>)[field];
			const initValue = (initData as Record<string, unknown>)[field];
			
			// Special handling for arrays (like skills)
			if (field === 'skills') {
				const currentSkills = Array.isArray(currentValue) ? currentValue : [];
				const initSkills = Array.isArray(initValue) ? initValue : [];
				
				// Compare array contents
				if (currentSkills.length !== initSkills.length || 
					!currentSkills.every((skill, index) => skill === initSkills[index])) {
					newChangedFields.add(field);
				}
			} else {
				// Regular comparison for other fields
				if (currentValue !== initValue) {
					newChangedFields.add(field);
				}
			}
		}
		
		changedFields = newChangedFields;
		hasChanges = newChangedFields.size > 0;
		
		// Debug logging in development
		if (import.meta.env.DEV && newChangedFields.size > 0) {
			console.log('Changed fields:', Array.from(newChangedFields));
		}
	}
	
	// Effects
	$effect(() => {
		initializeFields();
		
		if (!isInitialized) {
			initialData = JSON.parse(JSON.stringify(userData));
			isInitialized = true;
		}
	});
	
	$effect(() => {
		trackChanges();
	});
	
	// Derived values
	const activeTabData = $derived(tabs.find((tab: Tab) => tab.id === activeTab));
</script>

<div class="">
	<!-- Header -->
	<div class="p-2">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-2xl font-semibold text-gray-900">{title}</h1>
			{#if userType === 'seeker'}
				<button
					onclick={handlePreview}
					class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors flex items-center gap-2 hover:cursor-pointer"
				>
					Preview Profile
				</button>
			{/if}
		</div>
		
		<!-- Tabs -->
		<nav class="flex border-b border-gray-200 -mb-px">
			{#each tabs as tab (tab.id)}
				<button
					onclick={() => handleTabChange(tab.id)}
					class="px-4 py-3 text-sm font-medium border-b-2 transition-all cursor-pointer {activeTab === tab.id ? 'border-green-600 text-green-700' : 'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300'}"
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>
	
	<!-- Tab header with actions -->
	<div class="flex justify-between items-start border-b border-gray-200 px-6 py-4">
		<div>
			<h2 class="text-lg font-medium text-gray-900">{activeTabData?.title || ''}</h2>
			<p class="text-sm text-gray-500 mt-1">{activeTabData?.description || ''}</p>
		</div>
		
		<div class="flex items-center gap-3">
			{#if activeTab === 'documents'}
				<button
					onclick={() => {
						const event = new CustomEvent('openUploadModal');
						window.dispatchEvent(event);
					}}
					class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors flex items-center gap-2"
				>
					<Plus class="w-4 h-4" />
					Add Documents
				</button>
			{:else if hasChanges}
				<button
					onclick={handleSave}
					disabled={isSaving}
					class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 disabled:opacity-50 transition-colors flex items-center gap-2 hover:cursor-pointer"
				>
					{#if isSaving}
						<LoaderCircle class="w-4 h-4 animate-spin" />
						Saving...
					{:else}
						Save Changes
					{/if}
				</button>
				<button
					onclick={resetChanges}
					disabled={isSaving}
					class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50 disabled:opacity-50 transition-colors flex items-center gap-2 hover:cursor-pointer"
				>
					Cancel
				</button>
			{/if}
		</div>
	</div>
	
	<!-- Tab content -->
	<div>
		{#if activeTab === 'user'}
			<UserTab bind:userData onDirectSave={handleDirectSave} />
		{:else if activeTab === 'personal' && userType === 'seeker'}
			<PersonalInfoTab bind:profileData={userData} />
		{:else if activeTab === 'company' && userType === 'company'}
			<CompanyTab bind:companyData={userData} />
		{:else if activeTab === 'documents'}
			<DocumentsTab bind:documents={userData.documents} />
		{/if}
	</div>
</div>

<!-- Preview Drawer for Seeker Profile -->
{#if userType === 'seeker'}
	<SeekerProfilePreviewDrawer bind:isOpen={isPreviewOpen} {userData} />
{/if}