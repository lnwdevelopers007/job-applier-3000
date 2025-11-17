<script lang="ts">
	import { FileText, Download, Eye } from 'lucide-svelte';
	import FilePreviewModal from '$lib/components/files/FilePreviewModal.svelte';

	interface FileInfo {
		id: string;
		filename: string;
	}

	interface Props {
		files: FileInfo[];
		userName: string;
	}

	let { files, userName }: Props = $props();

	let showPreviewModal = $state(false);
	let selectedFileId = $state('');
	let selectedFilename = $state('');

	function handleViewDocument(fileId: string, filename: string) {
		selectedFileId = fileId;
		selectedFilename = filename;
		showPreviewModal = true;
	}

	function handleDownloadDocument(fileId: string, filename: string) {
		// Download document using the file ID
		const downloadUrl = `/files/download/${fileId}`;
		const link = document.createElement('a');
		link.href = downloadUrl;
		link.download = filename;
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}

	const hasDocuments = $derived(files && files.length > 0);
</script>

{#if hasDocuments}
	<!-- Dropdown menu for documents -->
	<div class="relative group">
		<div class="flex items-center gap-1 cursor-pointer border border-gray-300 px-3 py-1.5 rounded-xl hover:bg-gray-50 transition-colors">
			<FileText class="w-4 h-4 text-gray-500" />
			<span class="text-sm text-gray-700">{files.length} file{files.length === 1 ? '' : 's'}</span>
		</div>
		
		<!-- Dropdown content -->
		<div class="absolute top-full right-0 mt-1 w-64 bg-white border border-gray-200 rounded-lg shadow-lg z-50 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200">
			<div class="p-3">
				<div class="text-xs font-medium text-gray-500 mb-2">{userName}'s Documents</div>
				<div class="space-y-1 max-h-48 overflow-y-auto">
					{#each files as file (file.id)}
						<div class="flex items-center justify-between p-2 rounded">
							<div class="flex items-center gap-2 flex-1 min-w-0">
								<FileText class="w-3 h-3 text-gray-400 flex-shrink-0" />
								<span class="text-xs text-gray-700 truncate" title={file.filename}>
									{file.filename}
								</span>
							</div>
							<div class="flex items-center gap-1 ml-2">
								<button
									onclick={() => handleViewDocument(file.id, file.filename)}
									class="p-1 rounded hover:bg-gray-100"
									title="Preview document"
								>
									<Eye class="w-3 h-3 text-gray-500" />
								</button>
								<button
									onclick={() => handleDownloadDocument(file.id, file.filename)}
									class="p-1 rounded hover:bg-gray-100"
									title="Download document"
								>
									<Download class="w-3 h-3 text-gray-500" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
{:else}
	<span class="text-sm text-gray-400">No files</span>
{/if}

<FilePreviewModal
	bind:isOpen={showPreviewModal}
	fileId={selectedFileId}
	filename={selectedFilename}
/>