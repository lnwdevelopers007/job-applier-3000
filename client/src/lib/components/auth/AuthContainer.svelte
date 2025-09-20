<script>
  export let maxWidth = 'max-w-lg';
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-b from-gray-50 to-white py-8">
  <div class="w-full {maxWidth} p-8 bg-white shadow-sm border border-gray-200 rounded-xl">
    <slot />
  </div>
</div>