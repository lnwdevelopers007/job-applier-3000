<script lang="ts">
	export let isVisible: boolean = false;
	export let thing: string = '';
	export let thingName: string = '';
	export let actionName: string = '';
	export let isActionInProgress: boolean = false;
	export let reasonForAction: string = '';
	export let action: () => void;
</script>

{#if isVisible}
	<div class="fixed inset-0 z-50 flex items-center justify-center">
		<div class="w-96 rounded-lg bg-white p-6 shadow-lg">
			<h2 class="mb-4 text-lg font-semibold text-gray-800">
				{actionName}
				{thing}: {thingName}
			</h2>
			<p class="mb-2 block text-sm text-gray-600">Reason:</p>

			<textarea
				bind:value={reasonForAction}
				input={(e: any) => {
					if (e.target.value.length > 500) {
						reasonForAction = e.target.value.slice(0, 500);
					} else {
						reasonForAction = e.target.value;
					}
				}}
				rows="4"
				maxlength="500"
				placeholder="Enter the reason for this {actionName.toLowerCase()} action..."
				class="w-full rounded-md border border-gray-300 p-2 text-gray-800 focus:outline-none focus:ring-2 focus:ring-red-400"
			></textarea>

			<div class="mt-1 flex items-center justify-between text-sm">
				<span class={reasonForAction.length >= 480 ? 'text-red-500' : 'text-gray-500'}>
					{reasonForAction.length}/500 characters
				</span>
				{#if reasonForAction.length >= 500}
					<span class="font-medium text-red-600">Limit reached</span>
				{/if}
			</div>

			<div class="mt-4 flex justify-end space-x-3">
				<button
					class="rounded bg-gray-100 px-4 py-2 text-gray-700 hover:bg-gray-200"
					onclick={() => (isVisible = false)}
				>
					Cancel
				</button>
				<button
					class="rounded bg-red-600 px-4 py-2 text-white hover:bg-red-700 disabled:opacity-50"
					disabled={!reasonForAction.trim() || isActionInProgress}
					onclick={action}
				>
					Confirm {actionName}
				</button>
			</div>
		</div>
	</div>
{/if}
