<script>
  let { 
    id = '', 
    label = '', 
    value = $bindable(''), 
    required = false, 
    disabled = false, 
    placeholder = 'Select...' 
  } = $props();
</script>

<div>
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 mb-2">{label}</label>
  {/if}
  <select
    {id}
    bind:value
    {required}
    {disabled}
    class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors"
  >
    <option value="" disabled selected>{placeholder}</option>
    <slot />
  </select>
</div>