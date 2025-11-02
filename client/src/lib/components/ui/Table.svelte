<script lang="ts">
  import { ChevronLeft, ChevronRight, EllipsisVertical } from 'lucide-svelte';
  import StatusBadge from './StatusBadge.svelte';
  import Dropdown from './Dropdown.svelte';

  interface Column {
    key: string;
    label: string;
    width: string;
    align?: 'left' | 'center' | 'right';
    cellClass?: string;
    // Cell content rendering
    type?: 'text' | 'badge' | 'job' | 'actions' | 'custom';
    // For badge type
    badgeType?: 'status';
    // For job type (company logo + title)
    jobFields?: {
      logo: string;
      company: string;
      title: string;
    };
    // For actions type
    actions?: {
      onAction?: (action: string, row: any) => void;
      items?: Array<{ label: string; variant?: string }>;
    };
  }

  interface Props {
    columns: Column[];
    data: any[];
    loading?: boolean;
    currentPage?: number;
    totalPages?: number;
    itemsPerPage?: number;
    onPageChange?: (page: number) => void;
    emptyMessage?: string;
    emptyDescription?: string;
    rowClass?: string;
    defaultCellClass?: string;
    children?: any; // For custom cell rendering
    onAction?: (action: string, row: any) => void;
  }

  let {
    columns,
    data,
    loading = false,
    currentPage = 1,
    totalPages = 1,
    itemsPerPage = 10,
    onPageChange,
    emptyMessage = "No data found",
    emptyDescription = "No items available",
    rowClass = "hover:bg-gray-50 transition-colors",
    defaultCellClass = "px-6 py-4 whitespace-nowrap",
    children,
    onAction
  }: Props = $props();

  // State for dropdowns
  let openDropdown = $state<string | null>(null);
  let dropdownTriggers = $state<Record<string, HTMLElement>>({});

  // Calculate paginated data
  const paginatedData = $derived(
    data.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
  );

  function handlePageChange(page: number) {
    if (onPageChange) {
      onPageChange(page);
    }
  }

  function renderSkeletonCell(column: Column) {
    if (column.type === 'job') {
      return `
        <div class="flex items-center">
          <div class="flex-shrink-0 h-12 w-12">
            <div class="h-12 w-12 bg-gray-200 rounded-lg animate-pulse"></div>
          </div>
          <div class="ml-4 space-y-1">
            <div class="h-3 bg-gray-200 rounded-full w-24 animate-pulse"></div>
            <div class="h-3 bg-gray-200 rounded-full w-32 animate-pulse"></div>
          </div>
        </div>
      `;
    } else if (column.type === 'badge') {
      return '<div class="inline-flex items-center px-2.5 py-1.5 rounded-full bg-gray-200 w-16 h-6 animate-pulse"></div>';
    } else if (column.key === 'actions') {
      return '<div class="flex justify-center"><div class="w-8 h-8 bg-gray-200 rounded animate-pulse"></div></div>';
    } else {
      return '<div class="h-4 bg-gray-200 rounded-full w-20 animate-pulse"></div>';
    }
  }
</script>

<div class="bg-white">
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
            {#each columns as column}
              <th 
                class="px-6 py-3 text-{column.align || 'left'} text-xs font-medium text-gray-600 uppercase tracking-wider" 
                style="width: {column.width};"
              >
                {column.label}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#if loading}
            {#each Array(itemsPerPage) as _}
              <tr>
                {#each columns as column}
                  <td class={column.cellClass || defaultCellClass}>
                    {@html renderSkeletonCell(column)}
                  </td>
                {/each}
              </tr>
            {/each}
          {:else if children}
            <!-- Custom rendering with children slot -->
            {@render children({data: paginatedData, columns})}
          {:else}
            <!-- Auto-rendering based on column types -->
            {#each paginatedData as row (row.id)}
              <tr class={rowClass}>
                {#each columns as column}
                  <td class={column.cellClass || defaultCellClass}>
                    {#if column.type === 'badge' && column.badgeType === 'status'}
                      <StatusBadge status={row[column.key]} size="sm" />
                    {:else if column.type === 'job' && column.jobFields}
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-12 w-12">
                          <img 
                            class="h-12 w-12 rounded-lg p-1 object-contain" 
                            src={row[column.jobFields.logo]} 
                            alt={row[column.jobFields.company]}
                          />
                        </div>
                        <div class="ml-4 min-w-0 flex-1">
                          <div class="text-xs text-gray-500 truncate" title={row[column.jobFields.company]}>{row[column.jobFields.company]}</div>
                          <div class="text-sm font-medium text-gray-900 truncate" title={row[column.jobFields.title]}>{row[column.jobFields.title]}</div>
                        </div>
                      </div>
                    {:else if column.type === 'actions' && column.actions}
                      <div class="relative">
                        <button
                          bind:this={dropdownTriggers[row.id]}
                          onclick={(e) => {
                            e.stopPropagation();
                            openDropdown = openDropdown === row.id ? null : row.id;
                          }}
                          class="p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
                        >
                          <EllipsisVertical class="w-4 h-4" />
                        </button>
                        <Dropdown
                          items={column.actions.items?.map(item => ({
                            label: item.label,
                            action: () => onAction?.(item.label, row),
                            variant: item.variant as 'default' | 'danger' | undefined
                          })) || []}
                          isOpen={openDropdown === row.id}
                          onClose={() => openDropdown = null}
                          triggerElement={dropdownTriggers[row.id]}
                        />
                      </div>
                    {:else}
                      <div class="text-sm text-gray-900 truncate" title={row[column.key]}>
                        {row[column.key] ?? ''}
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

    <!-- Pagination -->
    {#if !loading && totalPages > 1}
      <div class="px-6 py-4 border border-t-0 border-gray-200 rounded-bl-lg rounded-br-lg bg-white">
        <div class="flex items-center justify-between">
          <!-- Results info -->
          <div class="text-sm text-gray-600">
            Showing <span class="font-medium text-gray-900">{(currentPage - 1) * itemsPerPage + 1}</span> to 
            <span class="font-medium text-gray-900">{Math.min(currentPage * itemsPerPage, data.length)}</span> of 
            <span class="font-medium text-gray-900">{data.length}</span> results
          </div>
          
          <!-- Pagination controls -->
          <div class="flex items-center space-x-2">
            <button
              onclick={() => handlePageChange(Math.max(1, currentPage - 1))}
              disabled={currentPage === 1}
              class="flex items-center px-2 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 hover:text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronLeft class="w-4 h-4" />
            </button>
            
            <div class="flex items-center space-x-1">
              {#each Array(totalPages) as _, i}
                {#if i < 2 || i >= totalPages - 2 || (i >= currentPage - 2 && i <= currentPage)}
                  <button
                    onclick={() => handlePageChange(i + 1)}
                    class="w-9 h-9 flex items-center justify-center text-sm font-medium rounded-lg transition-colors {
                      currentPage === i + 1
                        ? 'bg-green-50 border border-green-500 text-gray-900'
                        : 'text-gray-600 bg-white border border-gray-300 hover:bg-gray-50 hover:text-gray-700'
                    }"
                  >
                    {i + 1}
                  </button>
                {:else if i === 2 || i === totalPages - 3}
                  <span class="px-2 py-2 text-sm text-gray-400">...</span>
                {/if}
              {/each}
            </div>
            
            <button
              onclick={() => handlePageChange(Math.min(totalPages, currentPage + 1))}
              disabled={currentPage === totalPages}
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