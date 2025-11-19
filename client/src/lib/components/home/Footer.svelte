<script lang="ts">
	import { goto } from '$app/navigation';
	import PDPAModal from '../modals/PDPAModal.svelte';
	let showPDPA = $state(false);
	const currentYear = new Date().getFullYear();
</script>

<footer class="bg-white px-4 py-12 sm:px-6 lg:px-8">
	<div class="container mx-auto max-w-7xl">
		<div class="mb-8 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-4">
			<div>
				<button
					class="group mb-4 flex cursor-pointer items-start space-x-1 transition-opacity hover:opacity-80 focus:opacity-80 focus:outline-none"
					onclick={() => goto('/')}
					onkeydown={(e) => e.key === 'Enter' && goto('/')}
					aria-label="Go to home page"
				>
					<h2 class="text-lg font-semibold text-black">Job Applier</h2>
					<h2 class="text-lg font-semibold text-green-700">3000</h2>
				</button>
				<p class="text-sm text-gray-600">
					Your gateway to amazing tech opportunities at KU Computer Engineering
				</p>
			</div>
		</div>

		<div class="border-t border-gray-200 pt-8 text-center">
			<p class="text-sm text-gray-500">Copyright © {currentYear} lnwdevelopers007.</p>
			<p class="text-sm text-gray-500">
				This website is licensed under the <a
					href="https://www.gnu.org/licenses/agpl-3.0.txt"
					class="text-blue-500 hover:text-blue-700">GNU Affero General Public License Version 3.</a
				>
			</p>
			<p class="text-sm text-gray-500">
				The source code is available <a
					href="https://github.com/lnwdevelopers007/job-applier-3000/"
					class="text-blue-500 hover:text-blue-700">Here.</a
				>
			</p>
			<p class="text-sm text-gray-500">
				By using our service, you consent to the processing of your Personal Data as described in
				our

				<button
					type="button"
					onclick={() => (showPDPA = true)}
					class="font-medium text-blue-500 hover:text-blue-700"
				>
					Privacy Notice.
				</button>
			</p>
		</div>
	</div>
</footer>

<PDPAModal bind:isVisible={showPDPA} />
