<script lang="ts">
  import Skeleton from '../ui/Skeleton.svelte';
  import JobCardSkeleton from './JobCardSkeleton.svelte';
</script>

<div class="mx-auto pb-10">
  <div class="flex gap-6">
    
    <!-- Main Content Container -->
    <div class="flex-1">
      
      <!-- Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <!-- Left Column - Combined Job Header and Details Skeleton -->
        <div class="lg:col-span-2">
          <div class="p-8">
            
            <!-- Job Header Section Skeleton -->
            <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-8 pb-8 border-b border-gray-100">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-3">
                  <Skeleton width="w-10" height="h-10" class="rounded-lg" />
                  <div>
                    <Skeleton width="w-32" height="h-5" />
                  </div>
                </div>
                <Skeleton width="w-2/4" height="h-8" class="mb-3" />

                <div class="flex flex-wrap items-center gap-2 mb-4">
                  <Skeleton width="w-24" height="h-4" />
                  <Skeleton width="w-2" height="h-4" />
                  <Skeleton width="w-20" height="h-4" />
                </div>

                <div class="flex flex-wrap gap-2">
                  <Skeleton width="w-20" height="h-6" class="rounded-md" />
                  <Skeleton width="w-16" height="h-6" class="rounded-md" />
                </div>
              </div>


            </div>

            <!-- Job Description Section Skeleton -->
            <Skeleton width="w-32" height="h-6" class="mb-4" />
            <div class="space-y-3">
              <Skeleton width="w-full" height="h-4" />
              <Skeleton width="w-full" height="h-4" />
              <Skeleton width="w-3/4" height="h-4" />
              <Skeleton width="w-full" height="h-4" />
              <Skeleton width="w-5/6" height="h-4" />
            </div>

            <!-- Skills Section Skeleton -->
            <div class="mt-6 pt-6">
              <Skeleton width="w-16" height="h-6" class="mb-4" />
              <div class="flex flex-wrap gap-2">
                <Skeleton width="w-20" height="h-8" class="rounded-md" />
                <Skeleton width="w-24" height="h-8" class="rounded-md" />
                <Skeleton width="w-16" height="h-8" class="rounded-md" />
                <Skeleton width="w-20" height="h-8" class="rounded-md" />
                <Skeleton width="w-18" height="h-8" class="rounded-md" />
              </div>
            </div>

            <!-- About the Company Section Skeleton -->
            <div class="mt-6 pt-6">
              <Skeleton width="w-40" height="h-6" class="mb-4" />
              <div class="flex items-start gap-3 bg-slate-50 rounded-lg border border-gray-200 p-6">
                <Skeleton width="w-10" height="h-10" class="rounded" />
                <div class="flex-1">
                  <Skeleton width="w-32" height="h-5" class="mb-2" />
                  <Skeleton width="w-40" height="h-4" class="mb-2" />
                  <div class="flex items-center gap-2 mb-3">
                    <Skeleton width="w-24" height="h-4" />
                    <Skeleton width="w-20" height="h-4" />
                  </div>
                  <Skeleton width="w-full" height="h-4" class="mb-1" />
                  <Skeleton width="w-3/4" height="h-4" class="mb-2" />
                  <Skeleton width="w-20" height="h-4" />
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column Skeleton -->
        <div class="lg:col-span-1">
          <!-- Similar Jobs Skeleton -->
          <div class="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden p-6 mb-6">
            <Skeleton width="w-32" height="h-6" class="mb-3" />
            <div class="space-y-3">
              <JobCardSkeleton />
              <JobCardSkeleton />
              <JobCardSkeleton />
            </div>
          </div>

          <!-- Other Jobs Skeleton -->
          <div class="bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden p-6">
            <Skeleton width="w-40" height="h-6" class="mb-3" />
            <div class="space-y-3">
              <JobCardSkeleton />
              <JobCardSkeleton />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>