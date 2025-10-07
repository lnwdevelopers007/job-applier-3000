<script lang="ts">
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { SettingsService, SEEKER_TABS, INITIAL_SEEKER_DATA, type SeekerUserData } from '$lib/services/settingsService';
	import { LoaderCircle } from 'lucide-svelte';
	import { onMount } from 'svelte';
	
	// State using Svelte 5 runes
	let userData = $state<SeekerUserData>({ ...INITIAL_SEEKER_DATA });
	let loading = $state(true);
	let error = $state<string | null>(null);
	let activeTab = $state('user');
	
	// Utility functions
	function showError(message: string): void {
		error = message;
	}
	
	function clearMessages(): void {
		error = null;
	}
	
	// Main functions
	async function loadUserData(): Promise<void> {
		try {
			loading = true;
			clearMessages();
			
			userData = await SettingsService.loadUserData<SeekerUserData>('seeker');
			
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to load profile data. Please try again.';
			showError(errorMessage);
		} finally {
			loading = false;
		}
	}
	
	async function handleSave(data: SeekerUserData, changedFields?: string[]): Promise<void> {
		const updatedUser = await SettingsService.saveUserData(data, 'seeker', changedFields);
		SettingsService.updateLocalData(userData, updatedUser, 'seeker');
	}
	
	async function handleDocumentSave(data: SeekerUserData): Promise<void> {
		await SettingsService.saveDocuments(data);
		userData.documents = data.documents;
	}
	
	function setActiveTab(tab: string): void {
		activeTab = tab;
	}
	
	// Lifecycle
	onMount(() => {
		loadUserData();
	});
</script>

<!-- Loading State -->
{#if loading}
	<div class="flex items-center justify-center min-h-[900px]">
		<LoaderCircle class="animate-spin text-gray-300" size="48" strokeWidth="2" />
	</div>
{:else}
	<!-- Error Message -->
	{#if error}
		<div class="mb-6 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
			{error}
		</div>
	{/if}
	
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