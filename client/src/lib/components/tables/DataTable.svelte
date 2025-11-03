<script lang="ts" generics="TData">
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	type Column<T> = {
		id: string;
		accessorKey?: keyof T;
		header: string;
		cell?: (item: T) => any;
		width?: string;
		cellClass?: string;
	};

	type Props = {
		data: TData[];
		columns: Column<TData>[];
		showPagination?: boolean;
		pageSize?: number;
		loading?: boolean;
		emptyMessage?: string;
		emptyDescription?: string;
		rowClass?: string;
		defaultCellClass?: string;
		class?: string;
	};

	let {
		data = $bindable(),
		columns,
		showPagination = true,
		pageSize = 10,
		loading = false,
		emptyMessage = 'No data found',
		emptyDescription = 'No items available',
		rowClass = 'hover:bg-gray-50 transition-colors',
		defaultCellClass = 'px-6 py-4 whitespace-nowrap',
		class: className = ''
	}: Props = $props();

	let currentPage = $state(0);
	const totalPages = $derived(Math.ceil(data.length / pageSize));

	const paginatedData = $derived(() => {
		const start = currentPage * pageSize;
		const end = start + pageSize;
		return data.slice(start, end);
	});

	// Reset to first page when data changes
	$effect(() => {
		void data; // Explicitly mark as intentional side effect
		currentPage = 0;
	});

	function getCellValue(column: Column<TData>, item: TData) {
		if (column.cell && typeof column.cell === 'function') {
			return column.cell(item);
		}
		if (column.accessorKey) {
			return item[column.accessorKey];
		}
		return '';
	}

</script>

<div class={`bg-white ${className}`}>
	{#if !loading && data.length === 0}
		<div class="flex flex-col items-center justify-center py-16">
			<p class="text-gray-900 font-medium mb-1">{emptyMessage}</p>
			<p class="text-sm text-gray-500">{emptyDescription}</p>
		</div>
	{:else}
		<div class="overflow-x-auto rounded-tl-lg rounded-tr-lg border border-gray-200">
			<table class="w-full divide-y divide-gray-200" style="table-layout: fixed;">
				<thead class="bg-slate-50">
					<tr>
						{#each columns as column (column.id)}
							<th
								class="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider"
								style={column.width ? `width: ${column.width};` : ''}
							>
								{column.header}
							</th>
						{/each}
					</tr>
				</thead>
				<tbody class="bg-white divide-y divide-gray-200">
					{#if loading}
						{#each Array.from({ length: pageSize }, (_, index) => index) as index (index)}
							<tr>
								{#each columns as column (column.id)}
									<td class={column.cellClass || defaultCellClass}>
										<!-- Safe skeleton rendering without @html -->
										{#if column.id === 'job' || column.id === 'title'}
											<div class="flex items-center">
												<div class="flex-shrink-0 h-12 w-12">
													<div class="h-12 w-12 bg-gray-200 rounded-lg animate-pulse"></div>
												</div>
												<div class="ml-4 space-y-1">
													<div class="h-3 bg-gray-200 rounded-full w-24 animate-pulse"></div>
													<div class="h-3 bg-gray-200 rounded-full w-32 animate-pulse"></div>
												</div>
											</div>
										{:else}
											<div class="h-4 bg-gray-200 rounded-full w-full animate-pulse"></div>
										{/if}
									</td>
								{/each}
							</tr>
						{/each}
					{:else}
						{#each paginatedData() as item, itemIndex (item.id || itemIndex)}
							<tr class={rowClass}>
								{#each columns as column (column.id)}
									{@const value = getCellValue(column, item)}
									<td class={column.cellClass || defaultCellClass}>
										{#if value && typeof value === 'object' && 'component' in value}
											{@const Component = value.component}
											<Component {...value.props} />
										{:else}
											<div class="text-sm text-gray-900 truncate" title={String(value)}>
												{value ?? ''}
											</div>
										{/if}
									</td>
								{/each}
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>

		{#if !loading && showPagination && totalPages > 1}
			<div class="px-6 py-4 border border-t-0 border-gray-200 rounded-bl-lg rounded-br-lg bg-white">
				<div class="flex items-center justify-between">
					<div class="text-sm text-gray-600">
						Showing <span class="font-medium text-gray-900">{currentPage * pageSize + 1}</span> to
						<span class="font-medium text-gray-900"
							>{Math.min((currentPage + 1) * pageSize, data.length)}</span
						>
						of
						<span class="font-medium text-gray-900">{data.length}</span> results
					</div>

					<div class="flex items-center space-x-2">
						<button
							onclick={() => (currentPage = Math.max(0, currentPage - 1))}
							disabled={currentPage === 0}
							class="flex items-center px-2 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
						>
							<ChevronLeft class="w-4 h-4" />
						</button>

						<div class="flex items-center space-x-1">
							{#each Array.from({ length: totalPages }, (_, i) => i) as i (i)}
								{#if i < 2 || i >= totalPages - 2 || (i >= currentPage - 1 && i <= currentPage + 1)}
									<button
										onclick={() => (currentPage = i)}
										class="w-9 h-9 flex items-center justify-center text-sm font-medium rounded-lg transition-colors {currentPage ===
										i
											? 'bg-green-50 border border-green-500 text-gray-900'
											: 'text-gray-600 bg-white border border-gray-300 hover:bg-gray-50 hover:text-gray-700'}"
									>
										{i + 1}
									</button>
								{:else if (i === 2 && currentPage > 3) || (i === totalPages - 3 && currentPage < totalPages - 4)}
									<span class="px-2 py-2 text-sm text-gray-400">...</span>
								{/if}
							{/each}
						</div>

						<button
							onclick={() => (currentPage = Math.min(totalPages - 1, currentPage + 1))}
							disabled={currentPage === totalPages - 1}
							class="flex items-center px-2 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
						>
							<ChevronRight class="w-4 h-4" />
						</button>
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>