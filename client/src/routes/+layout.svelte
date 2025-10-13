<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from 'svelte-french-toast';
	import { authStore } from '$lib/stores/auth.svelte';
	import type { LayoutData } from './$types';

	let { children, data }: { children: unknown; data: LayoutData } = $props();

	$effect(() => {
		authStore.initFromPageData(data.user);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap" rel="stylesheet">
</svelte:head>

{@render children?.()}

<!-- Global Toast Notifications -->
<Toaster position="bottom-right" />
