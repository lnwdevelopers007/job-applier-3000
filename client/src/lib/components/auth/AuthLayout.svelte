<script>
  import { ArrowLeft } from 'lucide-svelte';
  
  let { backHref = '', backText = 'Back' } = $props();
</script>

<div class="min-h-screen flex">
  <div class="hidden lg:flex lg:w-2/3 bg-gradient-to-br from-green-50 to-blue-50">
    <slot name="left-side" />
  </div>

  <div class="w-full lg:w-1/2 flex items-center justify-center p-8 bg-white relative">
    {#if backHref}
      <a href={backHref} class="absolute top-8 left-8 flex items-center gap-2 text-gray-500 hover:text-gray-700 transition-colors hover:cursor-pointer">
        <ArrowLeft class="w-4 h-4" />
        <span class="text-sm">{backText}</span>
      </a>
    {/if}
    
    <div class="w-full max-w-lg">
      <slot />
    </div>
  </div>
</div>