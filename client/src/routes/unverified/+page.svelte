<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	function handleGoHome() {
		goto('/');
	}

	function handleLogin() {
		goto('/login');
	}
</script>

<svelte:head>
	<title>Account Verification Pending</title>
</svelte:head>

<div class="min-h-screen bg-white flex items-center justify-center">
	<div class="container mx-auto max-w-md text-left">
		<div class="">
			<h1 class="text-xl font-semibold text-gray-900 mb-2">
				Account Under Review
			</h1>
				
			<p class="text-sm text-gray-600 mb-4">
				{#if data.userName}
					Hi, {data.userName}. Your account is being verified.
				{:else}
					Your account is being verified.
				{/if}
				You'll receive an email once approved.
			</p>

			<p class="text-xs text-gray-500 mb-6">
				If you're not verified in 48 hours, please 
				<a href="#" class="text-blue-600 hover:text-blue-800 underline">contact us</a>.
			</p>

			<div class="space-y-2">
				<button
					onclick={handleGoHome}
					class="w-full px-3 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors hover:cursor-pointer"
				>
					Back to Home
				</button>
				
				<button
					onclick={handleLogin}
					class="w-full px-3 py-2 bg-white border border-gray-200 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors hover:cursor-pointer"
				>
					Go to Login
				</button>
			</div>
		</div>
	</div>
</div>	