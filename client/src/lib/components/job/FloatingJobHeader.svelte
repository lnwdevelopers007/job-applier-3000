<script lang="ts">
	import { Bookmark, Share2 } from 'lucide-svelte';
	
	interface Job {
		id: string;
		title: string;
		company: string;
		logo: string;
	}
	
	let {
		show = false,
		job,
		onApply,
		onBookmark,
		onShare,
		isBookmarked = false
	}: {
		show: boolean;
		job: Job;
		onApply: () => void;
		onBookmark?: () => void;
		onShare?: () => void;
		isBookmarked?: boolean;
	} = $props();
</script>

{#if job}
	<div class="fixed top-0 left-0 right-0 bg-white border-b border-gray-200 shadow-sm z-50 transition-all duration-300 ease-in-out {show ? 'translate-y-0 opacity-100' : '-translate-y-full opacity-0'}">
		<div class="max-w-7xl mx-auto px-4 py-3">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-4">
					<img src={job.logo} alt={job.company} class="w-10 h-10 rounded-lg hidden sm:block">
					<div>
						<h3 class="font-medium text-gray-900">{job.title}</h3>
						<p class="text-sm text-gray-600">{job.company}</p>
					</div>
				</div>
				<div class="flex items-center gap-2">
					<button 
						onclick={onApply}
						class="px-6 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 active:scale-[0.98] transition-all duration-150 hover:cursor-pointer"
					>
						Apply now
					</button>
					{#if onBookmark}
						<button 
							onclick={onBookmark}
							class="p-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center {isBookmarked ? 'text-green-700' : 'text-gray-700 bg-white'} hover:cursor-pointer"
						>
							<Bookmark class="w-4 h-4 {isBookmarked ? 'fill-current' : ''}" />
						</button>
					{/if}
					{#if onShare}
						<button 
							onclick={onShare}
							class="p-2 text-sm font-medium border border-gray-200 rounded-md hover:bg-gray-50 transition-colors flex items-center hover:cursor-pointer"
						>
							<Share2 class="w-4 h-4" />
						</button>
					{/if}
				</div>
			</div>
		</div>
	</div>
{/if}