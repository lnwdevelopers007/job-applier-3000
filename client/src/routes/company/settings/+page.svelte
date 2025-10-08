<script lang="ts">
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { SettingsService, COMPANY_TABS, INITIAL_COMPANY_DATA, type CompanyUserData } from '$lib/services/settingsService';
	import { onMount } from 'svelte';
	import { LoaderCircle } from 'lucide-svelte';
	
	// State using Svelte 5 runes
	let userData = $state<CompanyUserData>({ ...INITIAL_COMPANY_DATA });
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
			
			userData = await SettingsService.loadUserData<CompanyUserData>('company');
			
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to load profile data. Please try again.';
			showError(errorMessage);
		} finally {
			loading = false;
		}
	}
	
	async function handleSave(data: CompanyUserData, changedFields?: string[]): Promise<void> {
		const updatedUser = await SettingsService.saveUserData(data, 'company', changedFields);
		SettingsService.updateLocalData(userData, updatedUser, 'company');
	}
	
	async function handleDocumentSave(data: CompanyUserData): Promise<void> {
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
		tabs={COMPANY_TABS}
		bind:activeTab
		bind:userData
		onSave={activeTab === 'documents' ? handleDocumentSave : handleSave}
		userType="company"
		title="Company Settings"
		onTabChange={setActiveTab}
	/>
{/if}