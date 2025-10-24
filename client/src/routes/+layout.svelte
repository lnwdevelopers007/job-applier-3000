<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from 'svelte-french-toast';
	import { authStore } from '$lib/stores/auth.svelte';
	import LayoutHeader from '$lib/components/LayoutHeader.svelte';
	import { page } from '$app/stores';
	import type { LayoutData } from './$types';

	let { children, data }: { children: any; data: LayoutData } = $props();

	$effect(() => {
		authStore.initFromPageData(data.user || undefined);
	});

	// Check if current page should hide header
	const shouldHideHeader = $derived($page.url.pathname.startsWith('/login') || $page.url.pathname.startsWith('/signup'));
	
	// Get navigation items based on user role
	const navItems = $derived(
		!authStore.isAuthenticated ? [] :
		authStore.user?.role === 'company' ? [
			{ href: '/app/jobs', label: 'Jobs' },
			{ href: '/company/dashboard', label: 'Dashboard' },
			{ href: '/company/applicants', label: 'Applicants' },
			{ href: '/company/analytics', label: 'Analytics' },
			{ href: '/company/create', label: 'Post Job' },
			{ href: '/company/settings', label: 'Settings' }
		] : [
			{ href: '/', label: 'Home' },
			{ href: '/app/jobs', label: 'Jobs' },
			{ href: '/app/applications', label: 'Applications' },
			{ href: '/app/settings', label: 'Settings' }
		]
	);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap" rel="stylesheet">
</svelte:head>

<div class="bg-slate-50/25">
	{#if !shouldHideHeader}
	<LayoutHeader navItems={navItems} />
	{/if}

	{@render children()}
</div>

<!-- Global Toast Notifications -->
<Toaster 
  position="bottom-right" 
  toastOptions={{
    style: 'max-width: 500px;'
  }}
/>
