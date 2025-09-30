<script lang="ts">
	import type { ComponentType } from 'svelte';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { goto } from '$app/navigation';

	type CareerPath = {
		id: string;
		title: string;
		icon: ComponentType;
		iconStyle: string;
		positions: number;
	};

	let { path }: { path: CareerPath } = $props();
	let showAuthModal = $state(false);
	const Icon = path.icon;

	function handlePathClick() {
		const token = localStorage.getItem('access_token');
		
		if (token) {
			// User is logged in - go to jobs page with career filter
			goto(`/app/jobs?category=${encodeURIComponent(path.title)}`);
		} else {
			// User not logged in - show auth modal
			showAuthModal = true;
			sessionStorage.setItem('pendingNavigation', `/app/jobs?category=${encodeURIComponent(path.title)}`);
		}
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		sessionStorage.removeItem('pendingNavigation');
	}
</script>

<div
	class="bg-white border border-gray-200 rounded-xl p-6 hover:border-green-600 hover:shadow-sm transition-all duration-200 cursor-pointer h-full"
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
			<p class="text-sm text-gray-600 font-medium">{path.positions} open positions</p>
		</div>
	</div>
</div>

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title="Sign in to explore career path"
	description="Please log in or sign up to view jobs in {path.title}"
/>