<script lang="ts">
	import { marked } from 'marked'; // 1. Import the marked library
	import SafeHTML from '$lib/utils/SafeHTML.svelte';

	import pdpaMarkdown from '$lib/assets/pdpa.md?raw'; // 2. Import the .md file as raw text
	export let isVisible: boolean = false;

	// 3. Convert the Markdown string into an HTML string
	const pdpaHtml: string = marked.parse(pdpaMarkdown, { async: false }) as string;
</script>

{#if isVisible}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30 p-4">
		<div class="w-full max-w-3xl rounded-lg bg-white p-6 shadow-lg">
			<div class="prose prose-lg max-h-[70vh] max-w-none overflow-y-auto pr-4">
				<SafeHTML html={pdpaHtml} />
			</div>

			<div class="mt-4 flex justify-end space-x-3 border-t pt-4">
				<button
					class="rounded bg-gray-100 px-4 py-2 text-gray-700 hover:bg-gray-200"
					on:click={() => (isVisible = false)}
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}
