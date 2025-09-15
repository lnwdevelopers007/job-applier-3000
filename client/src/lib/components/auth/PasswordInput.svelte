<script>
  import { Eye, EyeOff } from 'lucide-svelte';
  
  export let id = 'password';
  export let label = 'Password';
  export let placeholder = 'Enter password';
  export let value = '';
  export let showPassword = false;
  export let required = true;
</script>

<div class="relative">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 mb-2">{label}</label>
  {/if}
  <input
    {id}
    type={showPassword ? 'text' : 'password'}
    bind:value
    {placeholder}
    {required}
    class="w-full px-4 py-3 text-sm pr-12 border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
  />
  <button
    type="button"
    onclick={() => showPassword = !showPassword}
    class="absolute right-3 {label ? 'top-[38px]' : 'top-1/2 -translate-y-1/2'} text-gray-500 hover:text-gray-700"
    aria-label={showPassword ? 'Hide password' : 'Show password'}
  >
    {#if showPassword}
      <EyeOff class="w-4 h-4" />
    {:else}
      <Eye class="w-4 h-4" />
    {/if}
  </button>
</div>