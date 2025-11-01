<script lang="ts">
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { SettingsService, FACULTY_TABS, INITIAL_FACULTY_DATA, type FacultyUserData } from '$lib/services/settingsService';
	import { authStore } from '$lib/stores/auth.svelte';
	import { LoaderCircle } from 'lucide-svelte';
	import { toast } from 'svelte-french-toast';
	
	let userData = $state<FacultyUserData>({ ...INITIAL_FACULTY_DATA });
	let loading = $state(true);
	let activeTab = $state('user');
	
	async function loadUserData(): Promise<void> {
		try {
			loading = true;
			
			userData = await SettingsService.loadUserData<FacultyUserData>('faculty');
			
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to load profile data. Please try again.';
			toast.error(errorMessage);
		} finally {
			loading = false;
		}
	}
	
	async function handleSave(data: FacultyUserData, changedFields?: string[]): Promise<void> {
		try {
			const updatedUser = await SettingsService.saveUserData(data, 'faculty', changedFields);
			SettingsService.updateLocalData(userData, updatedUser, 'faculty');
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to save profile. Please try again.';
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
		tabs={FACULTY_TABS}
		bind:activeTab
		bind:userData
		onSave={handleSave}
		userType="faculty"
		title="Faculty Settings"
		onTabChange={setActiveTab}
	/>
{/if}