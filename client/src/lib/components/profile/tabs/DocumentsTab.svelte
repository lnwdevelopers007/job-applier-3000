<script>
	let {
		documents = $bindable([])
		// onSave = () => {} // Unused - saving handled by parent
	} = $props();
	
	function handleDocumentUpload(event) {
		const files = event.target.files;
		if (files && files.length > 0) {
			if (!documents) {
				documents = [];
			}
			
			const newDocs = [];
			for (let file of files) {
				const colors = ['red', 'blue', 'purple', 'green', 'yellow', 'indigo'];
				const randomColor = colors[Math.floor(Math.random() * colors.length)];
				
				newDocs.push({
					name: file.name,
					size: `${(file.size / (1024 * 1024)).toFixed(1)} MB`,
					uploadDate: new Date().toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }),
					color: randomColor,
					file: file
				});
			}
			
			documents = [...documents, ...newDocs];
			event.target.value = '';
		}
	}
</script>

<div class="divide-y divide-gray-200">
	<div class="px-8 py-5">
		<div class="grid grid-cols-4 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Documents</label>
				<p class="text-xs text-gray-500 mt-1">Upload documents for recruiters such as resume, cover letter, portfolio, transcripts, or certificates</p>
			</div>
			<div class="col-span-2">
				<div class="space-y-4">
					{#if documents && documents.length > 0}
						{#each documents as doc, index (index)}
							<div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
								<div class="flex items-center gap-4">
									<div class="w-12 h-12 bg-{doc.color || 'gray'}-100 rounded-lg flex items-center justify-center">
										<svg class="w-6 h-6 text-{doc.color || 'gray'}-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
										</svg>
									</div>
									<div>
										<p class="text-sm font-medium text-gray-900">{doc.name}</p>
										<p class="text-xs text-gray-500">{doc.size} • Uploaded on {doc.uploadDate}</p>
									</div>
								</div>
								<div class="flex items-center gap-2">
									<button 
										class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors cursor-pointer"
										title="Download"
										aria-label="Download document"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
										</svg>
									</button>
									<button 
										onclick={() => {
											if (!documents) documents = [];
											const newDocs = documents.filter((_, i) => i !== index);
											documents = [...newDocs];
										}}
										class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
										title="Delete"
										aria-label="Delete document"
									>
										<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
										</svg>
									</button>
								</div>
							</div>
						{/each}
					{:else}
						<div class="text-center">
							<h3 class="mt-4 text-sm font-medium text-gray-900">No documents</h3>
							<p class="mt-2 text-sm text-gray-500 mb-3">Get started by uploading your first document.</p>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
	
	<input 
		type="file" 
		id="file-upload" 
		class="hidden" 
		accept=".pdf,.doc,.docx" 
		multiple
		onchange={handleDocumentUpload}
	/>
</div>