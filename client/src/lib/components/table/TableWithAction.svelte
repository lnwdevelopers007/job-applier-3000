<script lang="ts">
	export let things: any[] = [];
	export let tableHeader: string[] = [];
	export let rowAttributes: string[] = [];
	export let handleAction: (action: any, thing: any) => void;

	let selectedThingIndex: number | null = null;
	function selectRow(index: number) {
		selectedThingIndex = selectedThingIndex === index ? null : index;
	}
</script>

<!-- Jobs Table -->
<div class="mt-6 overflow-x-auto rounded-lg bg-white shadow">
	<table class="min-w-full text-left text-sm text-gray-600">
		<thead class="bg-gray-100 text-sm font-semibold text-gray-700">
			<tr>
				{#each tableHeader as header, i (i)}
					<th class="px-4 py-3">{header}</th>
				{/each}
				<th class="px-4 py-3">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each things as thing, i (i)}
				<tr
					class="cursor-pointer border-t {selectedThingIndex === i ? 'bg-green-200' : ''}"
					onclick={() => selectRow(i)}
				>
					{#each rowAttributes as attr, i (i)}
						{#if i == 0}
							<td class="px-4 py-3 font-medium text-gray-900"
								>{typeof thing[attr] === 'boolean' ? String(thing[attr]) : thing[attr]}</td
							>
						{:else}
							<td class="px-4 py-3"
								>{typeof thing[attr] === 'boolean' ? String(thing[attr]) : thing[attr]}</td
							>
						{/if}
					{/each}

					<td class="space-x-2 px-4 py-3">
						{#each thing.actions as action, i (i)}
							<button
								class="border-1 rounded border-gray-500 bg-gray-100 px-3 py-1 text-sm text-gray-900 enabled:hover:bg-gray-300 disabled:cursor-not-allowed disabled:opacity-50"
								disabled={action.disabled}
								onclick={(e) => {
									e.stopPropagation();
									handleAction(action, thing);
								}}
							>
								{action.label}
							</button>
						{/each}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
