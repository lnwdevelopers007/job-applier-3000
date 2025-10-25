<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(async () => {
		// Token is now stored as HttpOnly cookie by the backend
		// No need to store it in localStorage (more secure!)
		console.log('Authentication completed - token stored in secure cookie');

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