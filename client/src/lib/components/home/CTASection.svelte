<script lang="ts">
	import { goto } from "$app/navigation";
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { isAuthenticated, navigateWithAuth } from '$lib/utils/auth';

	let showAuthModal = $state(false);

	function handleExplore() {
		if (isAuthenticated()) {
			// User is logged in - go to jobs page
			navigateWithAuth('/app/jobs');
		} else {
			// User not logged in - show auth modal
			showAuthModal = true;
			sessionStorage.setItem('pendingNavigation', '/app/jobs');
		}
	}

	function handleGetStarted() {
		goto('/signup');
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		sessionStorage.removeItem('pendingNavigation');
	}
</script>

<section class="bg-gradient-to-l from-green-100/50 to-gray-50 py-20 px-4 sm:px-6 lg:px-8">
	<div class="container mx-auto max-w-4xl text-center">
		<h2 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-6">
			Ready to Launch Your Tech Career?
		</h2>
		<p class="text-lg sm:text-xl text-gray-600 mb-10 max-w-2xl mx-auto">
			Join thousands of KU Computer Engineering students who have successfully launched their careers
			with top tech companies
		</p>
		<div class="flex flex-col sm:flex-row gap-4 justify-center">
			<button
				onclick={handleExplore}
				class="px-6 py-3 bg-white text-green-700 font-semibold text-lg rounded-lg hover:bg-gray-50 hover:cursor-pointer transition-colors"
			>
				Explore Opportunities
			</button>
			<button
				onclick={handleGetStarted}
				class="px-6 py-3 bg-green-700 text-white font-semibold text-lg rounded-lg hover:bg-green-800 hover:cursor-pointer transition-colors"
			>
				Get Started Today
			</button>
		</div>
	</div>
</section>

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title="Sign in to explore opportunities"
	description="Please log in or sign up to view and explore job opportunities"
/>