<h1>Welcome to SvelteKit</h1>

{#if data.health.ok}
  <p>✅ Backend is healthy!</p>
{:else}
  <p>❌ Backend is not responding</p>
{/if}

<p>
  Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation
</p>

<script lang="ts">
  export let data: { health: { ok: boolean } };
</script>