<script lang="ts">
  import { MapPin, Bookmark } from 'lucide-svelte';
  import Badge from './Badge.svelte';
  import JobCardSkeleton from './JobCardSkeleton.svelte';
	import { formatRelativeTime } from '$lib/utils/datetime';
  
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
    tags?: string[];
  };

  interface Props {
    job?: Job;
    onclick?: () => void;
    onBookmark?: () => void;
    loading?: boolean;
  }

  let {
    job,
    onclick,
    onBookmark,
    loading = false
  }: Props = $props();

  let isBookmarked = $state(false);

  // Show max 3 skill tags, then "+n" for remaining
  const maxSkillTags = 3;
  const visibleSkills = $derived(job?.tags?.slice(0, maxSkillTags) || []);
  const remainingSkillsCount = $derived((job?.tags?.length || 0) - maxSkillTags);

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
  class="block rounded-lg bg-white border border-gray-200 p-3 hover:bg-gray-50 transition-all text-left w-full cursor-pointer"
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
  <div class="flex flex-wrap gap-1.5 mb-2">
    <Badge variant="purple" text={job.workArrangement} />
    <Badge variant="info" text={job.workType} />
    
    <!-- Required Skills -->
    {#if job.tags && job.tags.length > 0}
      {#each visibleSkills as skill}
        <Badge variant="secondary" text={skill} />
      {/each}
      
      {#if remainingSkillsCount > 0}
        <Badge variant="secondary" text="{remainingSkillsCount}+" />
      {/if}
    {/if}
  </div>
  <span class="text-xs text-gray-600">Posted {formatRelativeTime(job.posted)}</span>
</div>
{/if}