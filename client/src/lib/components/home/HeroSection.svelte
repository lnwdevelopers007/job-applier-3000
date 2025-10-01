<script lang="ts">
	import { Search } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import { jobSearchStore } from '$lib/stores/jobSearch';
	import { isAuthenticated } from '$lib/utils/auth';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';

	let searchQuery = $state('');
	let showAuthModal = $state(false);

	function handleSearch() {
		// Only proceed if there's a search query
		if (searchQuery.trim()) {
			if (isAuthenticated()) {
				// User is logged in - proceed with search
				jobSearchStore.setSearchQuery(searchQuery.trim());
				goto('/app/jobs');
			} else {
				// User not logged in - show modal
				showAuthModal = true;
				// Store the search query for after login
				sessionStorage.setItem('pendingSearch', searchQuery.trim());
			}
		}
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		// Clear the pending search if user cancels
		sessionStorage.removeItem('pendingSearch');
	}

	// TODO: Replace with actual API fetch from backend	
	const stats = [
		{ number: '500+', label: 'Active Jobs' },
		{ number: '120+', label: 'Top Companies' },
		{ number: '95%', label: 'Success Rate' }
	];
</script>

<section class="relative py-24 px-4 sm:px-6 lg:px-8 overflow-hidden bg-white">
	<div class="absolute inset-0" style="background: radial-gradient(circle at center, #FDFFFE 0%, #EEFFF4 100%)"></div>
	<div class="relative z-10 container mx-auto max-w-7xl">
		<h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 text-center leading-tight mb-6">
			<div class="flex items-center justify-center space-x-3">
				<span class="font-semibold text-black">Job Applier</span>
				<span class="font-semibold text-green-700">3000</span>
			</div>
		</h1>
		<p class="text-lg text-gray-600 text-center max-w-2xl mx-auto mb-10">
			Connect with leading tech companies and discover opportunities tailored for KU Computer
			Engineering students
		</p>

		<!-- Search Bar -->
		<div class="max-w-3xl mx-auto mb-12">
			<div class="relative">
				<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
					<Search class="h-5 w-5 text-gray-400" />
				</div>
				<input
					bind:value={searchQuery}
					type="text"
					placeholder="Search for jobs, companies, or skills..."
					class="block w-full pl-12 pr-32 py-4 text-md border border-gray-200 focus:shadow focus:shadow-gray-300 rounded-xl focus:ring-0 outline-none transition-all"
					onkeydown={(e) => e.key === 'Enter' && handleSearch()}
				/>
			</div>
		</div>

		<!-- Stats -->
		<div class="flex flex-wrap justify-center gap-8 sm:gap-12">
			{#each stats as stat, i (stat.label)}
				<div class="text-center px-8 {i < stats.length - 1 ? 'border-r border-gray-300' : ''}">
					<div class="text-3xl sm:text-4xl font-bold text-gray-900">{stat.number}</div>
					<div class="text-sm text-gray-500 font-medium mt-1">{stat.label}</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title="Sign in to search jobs"
	description="Please log in or sign up to search and view job opportunities"
/>