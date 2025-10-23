<script lang="ts">
  let {
    skill,
    onclick
  }: {
    skill: string;
    onclick?: () => void;
  } = $props();
</script>

{#if onclick}
  <button
    {onclick}
    class="inline-flex items-center px-2.5 py-1 rounded-md text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 transition-colors cursor-pointer"
  >
    {skill}
  </button>
{:else}
  <span class="inline-flex items-center px-2.5 py-1 rounded-md text-sm font-medium bg-gray-100 text-gray-700">
    {skill}
  </span>
{/if}