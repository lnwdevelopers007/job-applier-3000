<script>
	import ProfileSettings from '$lib/components/profile/ProfileSettings.svelte';
	import { mockCompanyData } from '$lib/stores/mockData.js';
	
	let activeTab = $state('company');
	let userData = $state({ ...mockCompanyData });
	
	const tabs = [
		{
			id: 'user',
			label: 'User Info',
			title: 'User info',
			description: 'Manage your personal account settings and profile'
		},
		{
			id: 'company',
			label: 'Company Info',
			title: 'Company Information',
			description: 'Update your company details and public information'
		},
		{
			id: 'benefits',
			label: 'Benefits',
			title: 'Benefits',
			description: 'Manage employee benefits and perks'
		}
	];
	
	function handleSave(data) {
		console.log('Saving company data:', data);
		alert('Company profile saved successfully!');
	}
</script>

<ProfileSettings 
	{tabs}
	bind:activeTab
	bind:userData
	onSave={handleSave}
	userType="company"
/>