<script lang="ts">
	export let things: any[] = [];
	export let table_header: string[] = [];
	export let row_attributes: string[] = [];
	export let selectedThing: number | null = null;
	export let selectRow: (index: number) => void;
	export let handleAction: (action: any, thing: any) => void;
</script>

<!-- Jobs Table -->
<div class="mt-6 overflow-x-auto rounded-lg bg-white shadow">
	<table class="min-w-full text-left text-sm text-gray-600">
		<thead class="bg-gray-100 text-sm font-semibold text-gray-700">
			<tr>
				{#each table_header as header}
					<th class="px-4 py-3">{header}</th>
				{/each}
				<th class="px-4 py-3">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each things as thing, i}
				<tr
					class="cursor-pointer border-t {selectedThing === i ? 'bg-green-200' : ''}"
					onclick={() => selectRow(i)}
				>
					{#each row_attributes as attr, i}
						{#if i == 0}
							<td class="px-4 py-3 font-medium text-gray-900">{thing[attr]}</td>
						{:else}
							<td class="px-4 py-3">{thing[attr]}</td>
						{/if}
					{/each}

					<td class="space-x-2 px-4 py-3">
						{#each thing.actions as action}
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
