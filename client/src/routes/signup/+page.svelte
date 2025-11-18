<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { isAuthenticated } from '$lib/utils/auth';
	import PDPAModal from '$lib/components/modals/PDPAModal.svelte';

	let showPDPA = $state(false);

	onMount(() => {
		if (isAuthenticated()) {
			goto('/app/jobs');
		}
	});
</script>

<div class="flex min-h-screen flex-col lg:flex-row">
	<!-- Left Side - Sign up section -->
	<div class="flex w-full items-center justify-center p-8 lg:w-3/6">
		<div class="w-full max-w-lg">
			<div class="mb-8 text-center">
				<h2 class="mb-2 text-2xl font-semibold text-gray-800">Sign up</h2>
				<div class="mb-3 flex items-baseline justify-center gap-1">
					<h1 class="text-2xl font-semibold text-gray-900">Job Applier</h1>
					<h1 class="text-2xl font-semibold text-green-700">3000</h1>
				</div>
				<p class="text-sm text-gray-500">Are you looking for a job or hiring talented graduates?</p>
			</div>

			<div class="mb-8 space-y-3">
				<!-- Option 1 - Student/Alumni -->
				<a
					href="/signup/student"
					class="block cursor-pointer rounded-lg border border-gray-200 p-4 transition-all duration-200 hover:border-green-600"
				>
					<div class="flex items-start space-x-3">
						<div class="flex-shrink-0">
							<div class="flex -space-x-2">
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-purple-100"
								></div>
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-blue-100"
								></div>
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-green-100"
								></div>
							</div>
						</div>
						<div class="flex-1">
							<h3 class="text-sm font-medium text-gray-900">I'm a student/alumni</h3>
							<p class="mt-1 text-xs text-gray-500">I'm a student/alumni looking for jobs.</p>
						</div>
					</div>
				</a>

				<!-- Option 2 - Recruiter -->
				<a
					href="/signup/company"
					class="block cursor-pointer rounded-lg border border-gray-200 p-4 transition-all duration-200 hover:border-green-600"
				>
					<div class="flex items-start space-x-3">
						<div class="flex-shrink-0">
							<div class="flex -space-x-2">
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-red-100"
								></div>
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-yellow-100"
								></div>
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full border-2 border-white bg-indigo-100"
								></div>
							</div>
						</div>
						<div class="flex-1">
							<h3 class="text-sm font-medium text-gray-900">I want to hire talents</h3>
							<p class="mt-1 text-xs text-gray-500">I'm a recruiter hiring talented graduates.</p>
						</div>
					</div>
				</a>
			</div>

			<div class="text-center">
				<p class="mt-8 text-center text-sm text-gray-500">
					By using our service, you consent to the processing of your Personal Data as described in
					our

					<button
						type="button"
						on:click={() => (showPDPA = true)}
						class="font-medium text-blue-600 hover:text-blue-700"
					>
						Privacy Notice
					</button>
				</p>

				<p class="text-sm text-gray-500">
					Already have an account?
					<a href="/login" class="font-medium text-green-600 hover:text-green-700">Log in</a>
				</p>
			</div>
		</div>
	</div>

	<!-- Right Side -->
	<div
		class="relative w-full overflow-hidden border-white bg-gradient-to-br from-green-50 to-blue-50 lg:w-3/5"
	>
		<div class="flex h-full flex-col justify-between p-12">
			<div class="max-w-2xl px-20 pt-20">
				<p class="mb-6 text-2xl leading-relaxed text-gray-700">
					Job Applier 3000 connects CPE students, alumni, and companies in a platform where students
					find opportunities and recruiters reach verified talent
				</p>
				<p class="text-2xl font-semibold text-gray-900">Join now to get started!</p>
			</div>

			<div class="flex justify-end">
				<div class="-mr-18 -mb-18 w-full max-w-4xl">
					<div class="border-5 overflow-hidden rounded-3xl border-black bg-white">
						<!-- Image placeholder -->
						<div class="h-160 flex w-full items-center justify-center bg-gray-200">
							<div class="text-center">
								<p class="text-lg text-gray-500">Dashboard Image</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
<PDPAModal bind:isVisible={showPDPA} />

<style>
	@keyframes fade-in {
		from {
			opacity: 0;
		}
	}

	@keyframes fade-out {
		to {
			opacity: 0;
		}
	}

	@keyframes slide-from-right {
		from {
			transform: translateX(30px);
		}
	}

	@keyframes slide-to-left {
		to {
			transform: translateX(-30px);
		}
	}

	:root::view-transition-old(root) {
		animation:
			90ms cubic-bezier(0.4, 0, 1, 1) both fade-out,
			300ms cubic-bezier(0.4, 0, 0.2, 1) both slide-to-left;
	}

	:root::view-transition-new(root) {
		animation:
			210ms cubic-bezier(0, 0, 0.2, 1) 90ms both fade-in,
			300ms cubic-bezier(0.4, 0, 0.2, 1) both slide-from-right;
	}
</style>
