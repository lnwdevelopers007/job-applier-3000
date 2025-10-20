<script lang="ts">
  import Skeleton from '../ui/Skeleton.svelte';
</script>

<div class="block rounded-lg border border-gray-200 p-4 w-full">
  <div class="flex items-start gap-3 mb-2">
    <Skeleton width="w-12" height="h-12" class="rounded flex-shrink-0" />
    <div class="flex-1 min-w-0">
      <div class="flex items-start justify-between mb-1">
        <Skeleton width="w-3/4" height="h-5" />
      </div>
      <div class="flex items-center gap-1 mb-3">
        <Skeleton width="w-1/2" height="h-4" />
        <div class="flex items-center gap-1 text-xs text-gray-600">
          <Skeleton width="w-3" height="h-3" />
          <Skeleton width="w-16" height="h-3" />
        </div>
      </div>
    </div>
  </div>
  <div class="flex flex-wrap gap-1.5">
    <Skeleton width="w-16" height="h-6" class="rounded-md" />
    <Skeleton width="w-20" height="h-6" class="rounded-md" />
  </div>
</div>