<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { jwtDecode } from 'jwt-decode';
	import { jobSearchStore } from '$lib/stores/jobSearch';
	import { AlertTriangle } from 'lucide-svelte';

	type TokenPayload = {
		email: string;
		name: string;
		exp: number;
		role?: string;
		avatarURL?: string;
		userID?: string;
		provider?: string;
	};

	let error: string | null = null;
	let isProcessing = true;

	onMount(async () => {
		const params = new URLSearchParams(window.location.search);
		const token = params.get('token');
		const errorParam = params.get('error');

		if (errorParam) {
			error = errorParam;
			isProcessing = false;
			setTimeout(() => {
				goto('/login');
			}, 3000);
			return;
		}

		if (token) {
			try {
				const decoded = jwtDecode<TokenPayload>(token);
				
				if (decoded.exp && decoded.exp * 1000 < Date.now()) {
					throw new Error('Token has expired');
				}

				localStorage.setItem('access_token', token);
				localStorage.setItem('user', JSON.stringify({
					email: decoded.email,
					name: decoded.name,
					role: decoded.role,
					avatarURL: decoded.avatarURL,
					userID: decoded.userID
				}));

				// Check for pending actions after login
				const pendingSearch = sessionStorage.getItem('pendingSearch');
				const pendingNavigation = sessionStorage.getItem('pendingNavigation');
				const pendingJobApplication = sessionStorage.getItem('pendingJobApplication');

				if (pendingSearch) {
					sessionStorage.removeItem('pendingSearch');
					// Set the search in the store and navigate to jobs page
					jobSearchStore.setSearchQuery(pendingSearch);
					await goto('/app/jobs');
					return;
				}

				if (pendingJobApplication) {
					sessionStorage.removeItem('pendingJobApplication');
					// Navigate to job application page
					await goto(`/app/jobs/${pendingJobApplication}/apply`);
					return;
				}

				if (pendingNavigation) {
					sessionStorage.removeItem('pendingNavigation');
					// Navigate to the stored URL
					await goto(pendingNavigation);
					return;
				}

				// Otherwise, redirect based on user role
				if (decoded.role === 'company') {
					await goto('/company/dashboard');
				} else if (decoded.role === 'jobSeeker') {
					await goto('/app/jobs');
				} else {
					await goto('/app/unverified');
				}
			} catch (err) {
				console.error('Failed to process token:', err);
				error = 'Authentication failed. Please try again.';
				isProcessing = false;
				
				setTimeout(() => {
					goto('/login');
				}, 3000);
			}
		} else {
			// No token provided
			error = 'No authentication token received.';
			isProcessing = false;
			setTimeout(() => {
				goto('/login');
			}, 2000);
		}
	});
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
	<div class="text-center">
		{#if isProcessing}
			<div class="space-y-4">
				<div class="inline-flex items-center justify-center w-10 h-10 border-3 border-gray-600 border-t-transparent rounded-full animate-spin"></div>
			</div>
		{:else if error}
			<div class="space-y-1">
				<div class="inline-flex items-center justify-center w-16 h-16 mb-4 bg-red-100 rounded-full">
					<AlertTriangle class="w-8 h-8 text-red-800" />
				</div>
				<p class="text-lg text-gray-600">{error}</p>
				<p class="text-sm text-gray-500">Redirecting to login...</p>
			</div>
		{/if}
	</div>
</div>
