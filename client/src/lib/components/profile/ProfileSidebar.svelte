<script>
  let { tabs = [], activeTab = '', onTabClick } = $props();
</script>

<div class="sticky top-8 self-start overflow-y-auto max-h-[calc(100vh-8rem)]">
    <nav class="space-y-2">
      {#each tabs as tab (tab.id)}
        <button
          onclick={() => onTabClick?.(tab.id)}
          class="w-full flex items-center gap-3 px-4 py-2.5 text-left text-sm rounded-lg transition-all duration-150 hover:cursor-pointer {activeTab === tab.id ? 'bg-green-50 text-green-700 font-medium' : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50 font-medium'}"
        >
          <span>{tab.label}</span>
        </button>
      {/each}
    </nav>
</div>