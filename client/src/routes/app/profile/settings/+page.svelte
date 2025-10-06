<script>
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { mockSeekerData } from '$lib/stores/mockData.js';
	
	let activeTab = $state('user');
	let userData = $state({ ...mockSeekerData });
	
	const tabs = [
		{
			id: 'user',
			label: 'User',
			title: 'User info',
			description: 'Update your photo and personal details here.'
		},
		{
			id: 'personal',
			label: 'Personal Info',
			title: 'Personal Info',
			description: 'Manage your personal information and preferences.'
		},
		{
			id: 'documents',
			label: 'Documents',
			title: 'Documents',
			description: 'Upload and manage your documents and files.'
		}
	];
	
	function handleSave(data) {
		console.log('Saving seeker data:', data);
		alert('Profile saved successfully!');
	}
</script>

<ProfileSettings 
	{tabs}
	bind:activeTab
	bind:userData
	onSave={handleSave}
	userType="seeker"
/>