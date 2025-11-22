<script lang="ts">
	import { onMount } from 'svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	onMount(async () => {
		// Initialize auth store with user data from server
		if (data.userInfo) {
			authStore.initFromPageData({
				email: data.userInfo.email || '',
				name: data.userInfo.name || '',
				avatarURL: data.userInfo.avatarURL,
				userID: data.userInfo.userID || '',
				role: data.userInfo.role as any || 'jobSeeker',
				verified: data.userInfo.verified || false,
				isAuthenticated: true
			});
		}

		// Clean up any pending actions
		sessionStorage.removeItem('pendingSearch');
		sessionStorage.removeItem('pendingNavigation');
		sessionStorage.removeItem('pendingJobApplication');

		// Wait a moment to ensure cookie is set, then redirect
		if (data.redirectPath) {
			// Wait for cookie to be available, then navigate
			setTimeout(() => {
				
				// Use window.location instead of goto() to preserve cookies
				window.location.href = data.redirectPath;
			}, 300);
		} else {
			// Fallback redirect
			setTimeout(() => {
				window.location.href = '/';
			}, 300);
		}
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
	<div class="text-center space-y-4">
		<div class="inline-flex items-center justify-center w-10 h-10 border-3 border-gray-600 border-t-transparent rounded-full animate-spin"></div>
		<p class="text-lg text-gray-600">Processing authentication...</p>
	</div>
</div>