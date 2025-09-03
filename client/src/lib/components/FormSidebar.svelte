<script>
  let { sections = [], activeSection = '', onSectionClick, onPreview, onPublish, actionLabel = 'Update Job Post' } = $props();
</script>

<div class="w-72 bg-white border border-gray-200 rounded-lg sticky top-24 self-start overflow-y-auto max-h-[calc(100vh-8rem)]">
  <div class="p-5">
    <div class="mb-8">
      <h3 class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-5">Form Sections</h3>
      <nav class="space-y-2">
        {#each sections as section, index}
          <button
            onclick={() => onSectionClick?.(section.id)}
            class="w-full flex items-center gap-3 px-3 py-2.5 text-left text-sm rounded-md transition-all duration-150 {activeSection === section.id ? 'bg-green-50 text-green-700 font-medium' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'}"
          >
            <div class="flex-shrink-0 w-2 h-2 rounded-full {activeSection === section.id ? 'bg-green-600' : 'bg-gray-300'}"></div>
            <span>{section.name}</span>
          </button>
        {/each}
      </nav>
    </div>

    <div class="border-t border-gray-200 pt-5 space-y-3">
      <button
        onclick={onPreview}
        class="w-full px-4 py-2.5 text-sm text-gray-600 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors font-medium"
      >
        Preview
      </button>
      
      <button
        onclick={onPublish}
        class="w-full px-4 py-2.5 text-sm bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors font-medium"
      >
        {actionLabel}
      </button>
    </div>
  </div>
</div>