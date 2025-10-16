<script lang="ts">
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { SettingsService, SEEKER_TABS, INITIAL_SEEKER_DATA, type SeekerUserData } from '$lib/services/settingsService';
	import { authStore } from '$lib/stores/auth.svelte';
	import { LoaderCircle } from 'lucide-svelte';
	import { toast } from 'svelte-french-toast';
	
	let userData = $state<SeekerUserData>({ ...INITIAL_SEEKER_DATA });
	let loading = $state(true);
	let activeTab = $state('user');
	
	async function loadUserData(): Promise<void> {
		try {
			loading = true;
			
			userData = await SettingsService.loadUserData<SeekerUserData>('seeker');
			
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to load profile data. Please try again.';
			toast.error(errorMessage);
		} finally {
			loading = false;
		}
	}
	
	async function handleSave(data: SeekerUserData, changedFields?: string[]): Promise<void> {
		try {
			const updatedUser = await SettingsService.saveUserData(data, 'seeker', changedFields);
			SettingsService.updateLocalData(userData, updatedUser, 'seeker');
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to save profile. Please try again.';
			toast.error(errorMessage);
			throw err; // Re-throw to let ProfileSettings handle it
		}
	}
	
	async function handleDocumentSave(data: SeekerUserData): Promise<void> {
		try {
			await SettingsService.saveDocuments(data);
			userData.documents = data.documents;
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to save documents. Please try again.';
			toast.error(errorMessage);
			throw err; // Re-throw to let ProfileSettings handle it
		}
	}
	
	function setActiveTab(tab: string): void {
		activeTab = tab;
	}
	
	// Lifecycle - wait for auth to be ready
	$effect(() => {
		if (authStore.isAuthenticated && authStore.user) {
			loadUserData();
		}
	});
</script>

<!-- Loading State -->
{#if loading}
	<div class="flex items-center justify-center min-h-[900px]">
		<LoaderCircle class="animate-spin text-gray-300" size="48" strokeWidth="2" />
	</div>
{:else}
	<!-- Profile Settings Component -->
	<ProfileSettings 
		tabs={SEEKER_TABS}
		bind:activeTab
		bind:userData
		onSave={activeTab === 'documents' ? handleDocumentSave : handleSave}
		userType="seeker"
		title="Profile Settings"
		onTabChange={setActiveTab}
	/>
{/if}