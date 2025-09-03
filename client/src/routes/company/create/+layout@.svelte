<script>
  import '../../../app.css';
  import LayoutHeader from '$lib/components/LayoutHeader.svelte';
  
  let { children } = $props();
</script>

<div class="bg-white">
  <LayoutHeader userName="Company" userRole="Recruiter" />

  <main class="flex-1">
    {@render children?.()}
  </main>
</div>