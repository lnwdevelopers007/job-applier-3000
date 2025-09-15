<script>
  export let maxWidth = 'max-w-lg';
</script>

<div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-blue-50 py-12">
  <div class="w-full {maxWidth} p-10 bg-white  border border-gray-200 shadow-sm">
    <slot />
  </div>
</div>