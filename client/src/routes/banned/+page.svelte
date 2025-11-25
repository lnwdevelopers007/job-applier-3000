<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { ShieldAlert, Mail, LogOut } from 'lucide-svelte';

	async function handleLogout() {
		try {
			await fetch(`${import.meta.env.VITE_LOCAL}/auth/google/logout`, {
				method: 'POST',
				credentials: 'include'
			});
		} catch (error) {
			console.error('Logout failed:', error);
		}

		goto('/');
	}

	onMount(() => {
		// Clear any pending actions in session storage
		sessionStorage.clear();
	});
</script>

<svelte:head>
	<title>Account Banned - Job Applier 3000</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-br from-red-50 via-orange-50 to-yellow-50 p-4"
>
	<div class="w-full max-w-md">
		<!-- Card -->
		<div class="rounded-2xl bg-white p-8 text-center shadow-xl">
			<!-- Icon -->
			<div class="mb-6 flex justify-center">
				<div class="flex h-20 w-20 items-center justify-center rounded-full bg-red-100">
					<ShieldAlert class="h-10 w-10 text-red-600" />
				</div>
			</div>

			<!-- Title -->
			<h1 class="mb-3 text-2xl font-bold text-gray-900">Account Suspended</h1>

			<!-- Message -->
			<p class="mb-6 leading-relaxed text-gray-600">
				Your account has been suspended and you can no longer access Job Applier 3000 services.
			</p>

			<!-- Info Box -->
			<div class="mb-6 rounded-lg border border-red-200 bg-red-50 p-4 text-left">
				<h2 class="mb-2 text-sm font-semibold text-red-900">Why was my account suspended?</h2>
				<p class="mb-3 text-sm text-red-800">
					Your account may have been suspended for violating our Terms of Service or Community
					Guidelines.
				</p>
				<h2 class="mb-2 text-sm font-semibold text-red-900">What can I do?</h2>
				<p class="text-sm text-red-800">
					If you believe this is a mistake, please contact our support team for assistance.
				</p>
			</div>

			<!-- Actions -->
			<div class="space-y-3">
				<a
					href="mailto:support@jobapplier3000.com?subject=Account%20Ban%20Appeal"
					class="flex w-full items-center justify-center gap-2 rounded-lg bg-red-600 px-4 py-3 text-sm font-medium text-white transition-colors hover:bg-red-700"
				>
					<Mail class="h-4 w-4" />
					Contact Support
				</a>

				<button
					onclick={handleLogout}
					class="flex w-full items-center justify-center gap-2 rounded-lg bg-gray-100 px-4 py-3 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-200"
				>
					<LogOut class="h-4 w-4" />
					Return to Home
				</button>
			</div>

			<!-- Footer note -->
			<p class="mt-6 text-xs text-gray-500">
				For urgent matters, you can reach us at support@jobapplier3000.com
			</p>
		</div>

		<!-- Additional info -->
		<div class="mt-6 text-center">
			<p class="text-sm text-gray-600">
				Need help? Visit our <a href="/" class="font-medium text-green-600 hover:text-green-700"
					>Help Center</a
				>
			</p>
		</div>
	</div>
</div>
