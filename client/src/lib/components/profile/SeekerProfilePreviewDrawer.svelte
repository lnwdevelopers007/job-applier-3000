<script lang="ts">
	import SideDrawer from '$lib/components/ui/SideDrawer.svelte';
	import SeekerProfileCard from './SeekerProfileCard.svelte';

	let { isOpen = $bindable(false), userData = {} } = $props();
</script>

<SideDrawer
	bind:isOpen
	title="Profile Preview"
	subtitle="This is how your profile appears to employers"
	width="820px"
>
	<div class="border border-gray-200 rounded-xl">
		<SeekerProfileCard 
			{userData}
			isPreviewMode={true}
			showApplicationInfo={false}
		/>
	</div>
</SideDrawer>
