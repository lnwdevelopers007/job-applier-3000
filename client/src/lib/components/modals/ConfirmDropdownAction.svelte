<script lang="ts">
	interface DropdownData {
		name: string;
		values: any[];
		defaultVal: any;
	}

	interface LinkData {
		title: string;
		link: string;
	}

	let {
		isVisible = $bindable(),
		actionName,
		action,
		showConfirmButton = $bindable(),
		dropdowns = $bindable(),
		linkData
	} = $props<{
		isVisible?: boolean;
		actionName: string;
		action: (selectedVals: Record<string, any>) => void;
		showConfirmButton?: boolean;
		dropdowns?: DropdownData[];
		linkData?: LinkData | null;
	}>();

	let selectedVals = $state(
		Object.fromEntries(dropdowns.map((d: DropdownData) => [d.name, d.defaultVal]))
	);

	$effect(() => {
		if (isVisible) {
			selectedVals = Object.fromEntries(dropdowns.map((d: DropdownData) => [d.name, d.defaultVal]));
		}
	});

	function onConfirm() {
		action(selectedVals);
	}
</script>

{#if isVisible}
	<div class="fixed inset-0 z-50 flex items-center justify-center">
		<div class="w-96 rounded-lg bg-white p-6 shadow-lg">
			<h2 class="mb-4 text-lg font-semibold text-gray-800">{actionName}</h2>
			{#each dropdowns as dropdown}
				<div>
					<label class="mb-2 block text-sm text-gray-600"
						>{dropdown.name}

						<select bind:value={selectedVals[dropdown.name]}>
							{#each dropdown.values as val}
								<option value={val}>{val}</option>
							{/each}
						</select>
					</label>
				</div>
			{/each}

			{#if linkData}
				<a href={linkData.link} class="text-blue-600 underline">
					{linkData.title}
				</a>
			{/if}
			<div class="mt-4 flex justify-end space-x-3">
				<button
					class="rounded bg-gray-100 px-4 py-2 text-gray-700 hover:bg-gray-200"
					onclick={() => (isVisible = false)}
				>
					Cancel
				</button>
				<button
					class="rounded bg-yellow-600 px-4 py-2 text-white hover:bg-yellow-700 disabled:opacity-50"
					disabled={showConfirmButton}
					onclick={onConfirm}
				>
					Confirm
				</button>
			</div>
		</div>
	</div>
{/if}
