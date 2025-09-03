<script>
  import '../../../app.css';
  import LayoutHeader from '$lib/components/LayoutHeader.svelte';
  
  let { children } = $props();
</script>

<div class="bg-white min-h-screen">
  <LayoutHeader />

  <div class="flex">
    <main class="flex-1">
      {@render children?.()}
    </main>
  </div>
</div>