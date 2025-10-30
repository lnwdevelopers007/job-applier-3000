<script lang="ts">
  import { AlertCircle, ChevronDown } from 'lucide-svelte';
  
  interface Option {
    value: string;
    label: string;
  }
  
  interface Props {
    label?: string;
    value?: string;
    options: Option[];
    placeholder?: string;
    required?: boolean;
    error?: string;
    showError?: boolean;
    class?: string;
    disabled?: boolean;
  }
  
  let { 
    label,
    value = $bindable(''),
    options,
    placeholder = 'Select an option',
    required = false,
    error = '',
    showError = false,
    class: className = '',
    disabled = false
  }: Props = $props();
  
  // Track if user has started typing to hide error
  let hideError = $state(false);
  // Track hover state for error icon
  let isErrorHovered = $state(false);
  
  // Hide error when user makes selection
  function handleChange() {
    hideError = true;
  }
  
  // Reset hideError when new validation occurs
  $effect(() => {
    if (showError && error) {
      hideError = false;
    }
  });
</script>

<div class={className}>
  {#if label}
    <!-- svelte-ignore a11y_label_has_associated_control -->
    <label class="block text-sm font-medium text-gray-700 mb-1">
      {label}
      {#if required}
        <span class="text-red-500">*</span>
      {/if}
    </label>
  {/if}
  
  <div class="relative">
    <select
      bind:value
      {disabled}
      onchange={handleChange}
      class="w-full px-4 {showError && error && !hideError ? 'pr-10' : 'pr-10'} py-2 text-sm border rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all appearance-none bg-white {disabled ? 'bg-gray-50 text-gray-500 cursor-not-allowed' : ''} {showError && error && !hideError ? 'border-red-500' : 'border-gray-300'}"
    >
      {#if placeholder}
        <option value="" disabled>{placeholder}</option>
      {/if}
      {#each options as option (option.value)}
        <option value={option.value}>{option.label}</option>
      {/each}
    </select>
    
    <!-- Dropdown arrow -->
    <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
      {#if showError && error && !hideError}
        <div
          class="pointer-events-auto"
          onmouseenter={() => isErrorHovered = true}
          onmouseleave={() => isErrorHovered = false}
          role="img"
          aria-label="Error"
        >
          <AlertCircle class="w-5 h-5 text-white mr-1 cursor-help" fill="red" />
        </div>
      {/if}
      <ChevronDown class="w-4 h-4 text-gray-400" />
    </div>
    
    {#if showError && error && !hideError && isErrorHovered}
      <!-- Floating error tooltip positioned below select - only show on hover -->
      <div class="absolute z-50 right-1 top-full mt-1 px-3 py-2 bg-red-500 text-white text-xs rounded-md shadow-sm whitespace-nowrap transition-opacity duration-200">
        <!-- Arrow pointing up to the alert icon -->
        <div class="absolute -top-1 right-4 w-2 h-2 bg-red-500 transform rotate-45"></div>
        {error}
      </div>
    {/if}
  </div>
</div>