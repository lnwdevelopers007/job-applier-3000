<script>
  let { 
    id = '', 
    type = 'text', 
    label = '', 
    placeholder = '', 
    value = $bindable(''), 
    required = false, 
    autocomplete = undefined, 
    pattern = undefined, 
    disabled = false 
  } = $props();
</script>

<div>
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 mb-2">{label}</label>
  {/if}
  <input
    {id}
    {type}
    bind:value
    {placeholder}
    {required}
    {autocomplete}
    {pattern}
    {disabled}
    class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors"
  />
</div>