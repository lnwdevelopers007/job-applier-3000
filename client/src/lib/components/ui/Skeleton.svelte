<script lang="ts">
  let {
    class: className = '',
    width = 'w-full',
    height = 'h-4'
  }: {
    class?: string;
    width?: string;
    height?: string;
  } = $props();
</script>

<div class="animate-pulse bg-gray-200 rounded {width} {height} {className}"></div>