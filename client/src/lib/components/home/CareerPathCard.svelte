<script lang="ts">
	import type { ComponentType } from 'svelte';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { jobSearchStore } from '$lib/stores/jobSearch';
	import { isAuthenticated, navigateWithAuth } from '$lib/utils/auth';

	type CareerPath = {
		id: string; // Used for keying in each loops
		title: string;
		icon: ComponentType;
		iconStyle: string;
		positions: number;
	};

	// eslint-disable-next-line svelte/no-unused-props
	let { path }: { path: CareerPath } = $props();
	let showAuthModal = $state(false);
	const Icon = path.icon;

	function handlePathClick() {
		if (isAuthenticated()) {
			// User is logged in - search for jobs in this category
			jobSearchStore.setSearchQuery(path.title);
			navigateWithAuth('/app/jobs');
		} else {
			// User not logged in - show auth modal
			showAuthModal = true;
			sessionStorage.setItem('pendingSearch', path.title);
		}
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		sessionStorage.removeItem('pendingSearch');
	}
</script>

<button
	class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-sm transition-all duration-200 cursor-pointer h-full w-full text-left"
	onclick={handlePathClick}
>
	<div class="flex items-center gap-4">
		<div class="w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0 {path.iconStyle}">
			<Icon class="w-6 h-6" />
		</div>
		<div class="flex-1 min-w-0">
			<h3
				class="font-semibold text-gray-900 transition-colors truncate"
				title={path.title}
			>
				{path.title}
			</h3>
		</div>
	</div>
</button>

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title="Sign in to explore career path"
	description="Please log in or sign up to view jobs in {path.title}"
/>