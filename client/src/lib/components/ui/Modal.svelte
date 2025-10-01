<script lang="ts">
	import { X } from 'lucide-svelte';
	import { fade, scale } from 'svelte/transition';
	import { onMount } from 'svelte';

	type ModalSize = 'sm' | 'md' | 'lg' | 'xl' | 'full';

	let {
		isOpen = $bindable(false),
		onClose,
		size = 'md',
		showCloseButton = true,
		closeOnBackdrop = false,
		closeOnEscape = true,
		class: className = '',
		children
	}: {
		isOpen: boolean;
		onClose?: () => void;
		size?: ModalSize;
		showCloseButton?: boolean;
		closeOnBackdrop?: boolean;
		closeOnEscape?: boolean;
		class?: string;
		children?: import('svelte').Snippet;
	} = $props();

	const sizeClasses: Record<ModalSize, string> = {
		sm: 'max-w-sm',     // 384px
		md: 'max-w-md',     // 448px
		lg: 'max-w-lg',     // 512px
		xl: 'max-w-xl',     // 576px
		full: 'max-w-4xl'   // 896px
	};

	function handleClose() {
		if (onClose) {
			onClose();
		} else {
			isOpen = false;
		}
	}

	function handleBackdropClick(event: MouseEvent) {
		// Only close if clicking directly on backdrop
		if (closeOnBackdrop && event.target === event.currentTarget) {
			handleClose();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (closeOnEscape && event.key === 'Escape') {
			handleClose();
			return;
		}
		
		// Trap focus within modal
		if (event.key === 'Tab') {
			trapFocus(event);
		}
	}

	function trapFocus(event: KeyboardEvent) {
		const modal = document.querySelector('[role="dialog"]') as HTMLElement;
		if (!modal) return;

		const focusableElements = modal.querySelectorAll(
			'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
		) as NodeListOf<HTMLElement>;

		const firstElement = focusableElements[0];
		const lastElement = focusableElements[focusableElements.length - 1];

		if (event.shiftKey) {
			if (document.activeElement === firstElement) {
				lastElement?.focus();
				event.preventDefault();
			}
		} else {
			if (document.activeElement === lastElement) {
				firstElement?.focus();
				event.preventDefault();
			}
		}
	}

	onMount(() => {
		// Always add keydown listener when mounted
		window.addEventListener('keydown', handleKeydown);
		
		return () => {
			window.removeEventListener('keydown', handleKeydown);
			// Restore body scroll on cleanup
			document.body.style.overflow = '';
			document.body.style.position = '';
			document.body.style.top = '';
			document.body.style.width = '';
		};
	});

	// Handle body scroll and interaction prevention when modal state changes
	$effect(() => {
		if (isOpen) {
			// Save current scroll position
			const scrollY = window.scrollY;
			
			// Prevent scroll and interaction
			document.body.style.position = 'fixed';
			document.body.style.top = `-${scrollY}px`;
			document.body.style.width = '100%';
			document.body.style.overflow = 'hidden';
			
			// Focus first focusable element in modal
			setTimeout(() => {
				const modal = document.querySelector('[role="dialog"]') as HTMLElement;
				const firstFocusable = modal?.querySelector(
					'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
				) as HTMLElement;
				firstFocusable?.focus();
			}, 100);
		} else {
			// Restore body and scroll position
			const scrollY = document.body.style.top;
			document.body.style.position = '';
			document.body.style.top = '';
			document.body.style.width = '';
			document.body.style.overflow = '';
			
			// Restore scroll position
			if (scrollY) {
				window.scrollTo(0, parseInt(scrollY.replace('-', '').replace('px', '')));
			}
		}
	});
</script>

{#if isOpen}
	<!-- Full-screen overlay to block interaction -->
	<div
		class="fixed inset-0 z-50 overflow-y-auto"
		transition:fade={{ duration: 200 }}
	>
		<!-- Backdrop with click handler -->
		<button
			class="fixed inset-0 bg-black/50 cursor-default"
			onclick={handleBackdropClick}
			onkeydown={(e) => e.key === 'Enter' && handleClose()}
			aria-label="Close modal"
			tabindex="-1"
		></button>

		<!-- Modal Container -->
		<div class="flex min-h-full items-center justify-center p-4">
			<!-- Modal Content -->
			<div
				class="relative bg-white rounded-lg shadow-xl w-full {sizeClasses[size]} {className}"
				transition:scale={{ duration: 200, start: 0.95 }}
				role="dialog"
				aria-modal="true"
			>
				{#if showCloseButton}
					<button
						onclick={handleClose}
						class="absolute top-4 right-4 z-10 p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 cursor-pointer transition-colors"
						aria-label="Close"
					>
						<X class="w-4 h-4" />
					</button>
				{/if}

				{@render children()}
			</div>
		</div>
	</div>
{/if}