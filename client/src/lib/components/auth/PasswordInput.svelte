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
    <label for={id} class="block text-sm font-medium text-gray-700 mb-1.5">{label}</label>
  {/if}
  <div class="relative">
    <input
      {id}
      type={showPassword ? 'text' : 'password'}
      bind:value
      {placeholder}
      {required}
      class="w-full pl-10 pr-10 py-2.5 text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent outline-none transition-all"
    />
    <div class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
      </svg>
    </div>
    <button
      type="button"
      onclick={() => showPassword = !showPassword}
      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
      aria-label={showPassword ? 'Hide password' : 'Show password'}
    >
      {#if showPassword}
        <EyeOff class="w-4 h-4" />
      {:else}
        <Eye class="w-4 h-4" />
      {/if}
    </button>
  </div>
</div>