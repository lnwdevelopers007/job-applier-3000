<script lang="ts">
  import LayoutHeader from '$lib/components/LayoutHeader.svelte';
  
  export let data: { health: { ok: boolean } };
</script>

<div class="min-h-screen bg-gray-50">
  <LayoutHeader />
  
  <main class="container mx-auto px-6 py-12">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-gray-900 mb-2">Welcome</h1>
      
      <div class="mb-8">
        {#if data.health.ok}
          <p class="text-green-600">✅ System is operational</p>
        {:else}
          <p class="text-red-600">❌ Backend is not responding</p>
        {/if}
      </div>
    </div>
  </main>
</div>