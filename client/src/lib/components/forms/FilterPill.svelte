<script lang="ts">
  import { ChevronDown, ArrowUpDown } from 'lucide-svelte';
  import { clickOutside } from '$lib/utils/clickOutside';

  interface Option {
    value: string;
    label: string;
  }

  interface Props {
    label: string;
    options: Option[];
    selectedValue: string;
    onSelectionChange?: (value: string) => void;
    type?: 'filter' | 'sort';
  }

  let {
    label,
    options,
    selectedValue = $bindable(''),
    onSelectionChange,
    type = 'filter'
  }: Props = $props();

  let isOpen = $state(false);
  let dropdownRef: HTMLDivElement;

  const isActive = $derived(selectedValue !== '');
  const displayLabel = $derived(() => {
    if (!selectedValue) {
      return type === 'sort' ? 'Sort by' : label;
    }
    const option = options.find(opt => opt.value === selectedValue);
    return option?.label || label;
  });

  const colorClasses = $derived(() => {
    if (type === 'sort') {
      return isActive 
        ? 'bg-blue-600 text-white hover:bg-blue-700' 
        : 'bg-white text-gray-700 border border-gray-200 hover:bg-gray-50 hover:border-gray-300';
    } else {
      return isActive 
        ? 'bg-green-600 text-white hover:bg-green-700' 
        : 'bg-white text-gray-700 border border-gray-200 hover:bg-gray-50 hover:border-gray-300';
    }
  });

  const highlightClasses = $derived(() => {
    return type === 'sort' ? 'bg-blue-50 text-blue-700' : 'bg-green-50 text-green-700';
  });

  function toggleDropdown() {
    isOpen = !isOpen;
  }

  function handleOptionSelect(value: string) {
    selectedValue = selectedValue === value ? '' : value;
    onSelectionChange?.(selectedValue);
    isOpen = false;
  }

  function handleClickOutside() {
    isOpen = false;
  }
</script>

<div class="relative" bind:this={dropdownRef} use:clickOutside={handleClickOutside}>
  <button
    type="button"
    onclick={toggleDropdown}
    class={`inline-flex items-center px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${colorClasses()}`}
  >
    {#if type === 'sort'}
      <ArrowUpDown class="mr-2 h-4 w-4" />
    {/if}
    {displayLabel()}
    <ChevronDown class={`ml-2 h-4 w-4 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
  </button>

  {#if isOpen}
    <div class="absolute top-full left-0 mt-1 w-64 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
      <div class="p-2">
        <div class="space-y-1 max-h-60 overflow-y-auto">
          <button
            type="button"
            onclick={() => handleOptionSelect('')}
            class={`w-full text-left px-3 py-2 text-sm rounded hover:bg-gray-50 ${
              selectedValue === '' ? highlightClasses() : 'text-gray-700'
            }`}
          >
            {type === 'sort' ? 'Default' : `All ${label}`}
          </button>
          {#each options as option}
            <button
              type="button"
              onclick={() => handleOptionSelect(option.value)}
              class={`w-full text-left px-3 py-2 text-sm rounded hover:bg-gray-50 ${
                selectedValue === option.value ? `${highlightClasses()} font-medium` : 'text-gray-700'
              }`}
            >
              {option.label}
            </button>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>