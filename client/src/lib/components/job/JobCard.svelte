<script lang="ts">
  import { MapPin, Bookmark } from 'lucide-svelte';
  import Badge from './Badge.svelte';
  import JobCardSkeleton from './JobCardSkeleton.svelte';
  
  type Job = {
    id: string;
    title: string;
    company: string;
    location: string;
    logo: string;
    workType: string;
    workArrangement: string;
    posted: string;
    salary?: string;
  };

  interface Props {
    job?: Job;
    onclick?: () => void;
    onBookmark?: () => void;
    loading?: boolean;
    showSalary?: boolean;
  }

  let {
    job,
    onclick,
    onBookmark,
    loading = false,
    showSalary = false
  }: Props = $props();

  let isBookmarked = $state(false);

  function handleBookmark(event: Event) {
    event.stopPropagation();
    isBookmarked = !isBookmarked;
    onBookmark?.();
  }
  
</script>

{#if loading || !job}
  <JobCardSkeleton />
{:else}
<div
  role="button"
  tabindex="0"
  onclick={onclick}
  onkeydown={(e) => e.key === 'Enter' && onclick?.()}
  class="block rounded-lg border border-gray-200 p-4 hover:bg-gray-50 transition-all text-left w-full cursor-pointer"
>
  <div class="flex items-start gap-3">
    <img src={job.logo} alt={job.company} class="w-12 h-12 rounded flex-shrink-0">
    <div class="flex-1 min-w-0">
      <div class="flex items-start justify-between">
        <h3 class="text-md font-medium text-gray-900">{job.title}</h3>
        <button
          onclick={handleBookmark}
          onkeydown={(e) => e.key === 'Enter' && handleBookmark(e)}
          class="hover:bg-gray-100 rounded transition-colors"
          aria-label="Bookmark job"
        >
          <Bookmark class="w-5 h-5 {isBookmarked ? 'text-green-600 fill-current' : 'text-gray-400 hover:text-gray-600'} hover:cursor-pointer" />
        </button>
      </div>
      <div class="space-y-1 mb-3">
        <p class="text-sm text-gray-600">{job.company}</p>
        <div class="flex items-center gap-1 text-xs text-gray-600">
          <MapPin class="w-3 h-3" />
          <span>{job.location}</span>
        </div>
      </div>
    </div>
  </div>
  <div class="flex flex-wrap gap-1.5">
    <Badge variant="purple" text={job.workArrangement} />
    <Badge variant="info" text={job.workType} />
  </div>
</div>
{/if}