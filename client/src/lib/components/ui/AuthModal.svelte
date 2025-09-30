<script lang="ts">
	import Modal from './Modal.svelte';
	import { goto } from '$app/navigation';
	import { LogIn, UserPlus } from 'lucide-svelte';

	let {
		isOpen = $bindable(false),
		onClose,
		title = 'Sign in to continue',
		description = 'Please log in or sign up to access this feature',
		pendingAction
	}: {
		isOpen: boolean;
		onClose?: () => void;
		title?: string;
		description?: string;
		pendingAction?: () => void;
	} = $props();

	function handleLogin() {
		// Store pending action if provided
		if (pendingAction) {
			sessionStorage.setItem('pendingAction', 'true');
		}
		isOpen = false;
		goto('/login');
	}

	function handleSignup() {
		// Store pending action if provided
		if (pendingAction) {
			sessionStorage.setItem('pendingAction', 'true');
		}
		isOpen = false;
		goto('/signup');
	}

	function handleClose() {
		if (onClose) {
			onClose();
		} else {
			isOpen = false;
		}
	}
</script>

<Modal bind:isOpen {onClose} size="lg" showCloseButton={true}>
	{#snippet children()}
		<div class="px-12 py-16">
			<!-- Content -->
			<h2 class="text-2xl font-semibold text-gray-900 text-center mb-3">
				{title}
			</h2>
			<p class="text-gray-600 text-md text-center mb-8">
				{description}
			</p>

			<!-- Actions -->
			<div class="space-y-3">
				<button
					onclick={handleLogin}
					class="w-full px-4 py-2 bg-green-600 text-white text-md font-medium rounded-lg hover:bg-green-700 cursor-pointer transition-colors flex items-center justify-center gap-2"
				>
					Log In
				</button>

				<button
					onclick={handleSignup}
					class="w-full px-4 py-2 bg-white border border-gray-200 text-gray-900 font-medium rounded-lg hover:bg-gray-50 cursor-pointer transition-colors flex items-center justify-center gap-2"
				>
					Sign Up
				</button>
			</div>
		</div>
	{/snippet}
</Modal>