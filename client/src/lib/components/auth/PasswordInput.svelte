<script>
  import { Eye, EyeOff, Lock } from 'lucide-svelte';
  
  let { 
    id = 'password', 
    label = 'Password', 
    placeholder = 'Enter password...', 
    value = $bindable(''), 
    showPassword = $bindable(false), 
    required = true, 
    showIcon = false 
  } = $props();
</script>

<div>
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-700 mb-2">{label}</label>
  {/if}
  <div class="relative">
    <input
      {id}
      type={showPassword ? 'text' : 'password'}
      bind:value
      {placeholder}
      {required}
      class="w-full {showIcon ? 'pl-10' : 'px-4'} py-3 pr-12 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors"
    />
    {#if showIcon}
      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 w-5 h-5" />
    {/if}
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