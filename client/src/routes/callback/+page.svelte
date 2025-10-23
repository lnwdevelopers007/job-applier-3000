<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(async () => {
		// Get the access token from URL query parameter
		const urlParams = new URLSearchParams(window.location.search);
		const token = urlParams.get('token');

		if (token) {
			// Store the access token in localStorage
			localStorage.setItem('access_token', token);
			console.log('Access token saved successfully');
		} else {
			console.error('No token found in URL');
		}

		// Clean up any pending actions
		sessionStorage.removeItem('pendingSearch');
		sessionStorage.removeItem('pendingNavigation');
		sessionStorage.removeItem('pendingJobApplication');

		// Redirect to home or dashboard after a short delay
		setTimeout(() => {
			goto('/');
		}, 1000);
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
	<div class="text-center space-y-4">
		<div class="inline-flex items-center justify-center w-10 h-10 border-3 border-gray-600 border-t-transparent rounded-full animate-spin"></div>
		<p class="text-lg text-gray-600">Processing authentication...</p>
	</div>
</div>