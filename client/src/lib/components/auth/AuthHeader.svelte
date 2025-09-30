<script>
  import { goto } from '$app/navigation';
  let { showLogo = true, title = '', subtitle = '', centered = true } = $props();
</script>

<div class="{centered ? 'text-center' : ''} {showLogo ? 'mb-12' : 'mb-8'}">
  {#if showLogo}
    <div class="flex items-baseline {centered ? 'justify-center' : ''} gap-1 {title ? 'mb-8' : ''}">
      <button 
        class="flex items-baseline gap-1 group cursor-pointer focus:outline-none focus:opacity-80"
        onclick={() => goto('/')}
        onkeydown={(e) => e.key === 'Enter' && goto('/')}
        aria-label="Go to home page"
      >
        <span class="text-2xl font-semibold text-gray-900">Job Applier</span>
        <span class="text-2xl font-semibold text-green-700">3000</span>
      </button>
    </div>
  {/if}
  
  {#if title}
    <h1 class="text-3xl font-semibold text-gray-900 mb-2">{title}</h1>
  {/if}
  
  {#if subtitle}
    <p class="text-sm text-gray-500">{subtitle}</p>
  {/if}
</div>