<script lang="ts">
	// client/src/lib/components/files/FileUploadModal.svelte
	import Modal from '../ui/Modal.svelte';
	import Select from '../forms/Select.svelte';
	import { FileText, Upload, Loader } from 'lucide-svelte';
	import { fileService } from '$lib/services/fileService';

	let {
		isOpen = $bindable(false),
		userRole = 'jobSeeker',
		onUploadSuccess = (file: any) => {}
	} = $props();

	let selectedFile = $state<File | null>(null);
	let selectedCategory = $state('');
	let isUploading = $state(false);
	let error = $state<string | null>(null);
	let dragOver = $state(false);

	// Get valid categories based on user role
	const categories = $derived(fileService.getValidCategories(userRole));

	// Set default category when role changes
	$effect(() => {
		if (categories.length > 0 && !selectedCategory) {
			selectedCategory = categories[0].value;
		}
	});

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			validateAndSetFile(target.files[0]);
		}
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		dragOver = false;

		if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
			validateAndSetFile(event.dataTransfer.files[0]);
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		dragOver = true;
	}

	function handleDragLeave() {
		dragOver = false;
	}

	function validateAndSetFile(file: File) {
		error = null;

		// Validate file type
		if (file.type !== 'application/pdf') {
			error = 'Only PDF files are allowed';
			return;
		}

		// Validate file size (10MB)
		if (file.size > 10 * 1024 * 1024) {
			error = 'File size must be less than 10MB';
			return;
		}

		selectedFile = file;
	}

	function removeFile() {
		selectedFile = null;
		error = null;
	}

	async function handleUpload() {
		if (!selectedFile || !selectedCategory) return;

		isUploading = true;
		error = null;

		try {
			const result = await fileService.uploadFile(selectedFile, selectedCategory);
			onUploadSuccess(result);
			selectedFile = null;
			error = null;
			isOpen = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to upload file';
			console.error('Upload error:', err);
		} finally {
			isUploading = false;
		}
	}

	function handleClose() {
		if (isUploading) return;

		isOpen = false;
		selectedFile = null;
		error = null;

		// Reset category to default
		if (categories.length > 0) {
			selectedCategory = categories[0].value;
		}
	}
</script>

<Modal bind:isOpen size="md" onClose={handleClose} closeOnBackdrop={!isUploading}>
	<div class="p-6">
		<h2 class="mb-4 text-xl font-medium text-gray-900">Upload Document</h2>

		<!-- Category Selection -->
		<div class="mb-4">
			<!-- svelte-ignore a11y_label_has_associated_control -->
			<label class="mb-2 block text-sm font-medium text-gray-700"> Document Category </label>
			<Select
				bind:value={selectedCategory}
				options={categories}
				placeholder="Select category..."
				disabled={isUploading}
			/>
		</div>

		<!-- File Upload Area -->
		<div class="mb-4">
			<!-- svelte-ignore a11y_label_has_associated_control -->
			<label class="mb-2 block text-sm font-medium text-gray-700"> Select File </label>

			{#if !selectedFile}
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					class="rounded-lg border-2 border-dashed p-8 text-center transition-colors {dragOver
						? 'border-green-500 bg-green-50'
						: 'border-gray-300 hover:border-gray-400'}"
					ondrop={handleDrop}
					ondragover={handleDragOver}
					ondragleave={handleDragLeave}
				>
					<Upload class="mx-auto mb-3 h-12 w-12 text-gray-400" />
					<p class="mb-2 text-sm text-gray-600">Drag and drop your PDF here, or</p>
					<label
						class="inline-block cursor-pointer rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50"
					>
						Browse Files
						<input
							type="file"
							accept="application/pdf"
							class="hidden"
							onchange={handleFileSelect}
							disabled={isUploading}
						/>
					</label>
					<p class="mt-3 text-xs text-gray-500">PDF only, max 10MB</p>
				</div>
			{:else}
				<div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
					<div class="flex items-center gap-3">
						<div
							class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded bg-red-100"
						>
							<FileText class="h-5 w-5 text-red-600" />
						</div>
						<div class="min-w-0 flex-1">
							<p class="truncate text-sm font-medium text-gray-900">{selectedFile.name}</p>
							<p class="text-xs text-gray-500">{fileService.formatFileSize(selectedFile.size)}</p>
						</div>
						{#if !isUploading}
							<button
								onclick={removeFile}
								class="text-gray-400 transition-colors hover:text-red-600"
								aria-label="Remove file"
							>
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M6 18L18 6M6 6l12 12"
									/>
								</svg>
							</button>
						{/if}
					</div>
				</div>
			{/if}
		</div>

		<!-- Error Message -->
		{#if error}
			<div class="mb-4 rounded-md border border-red-200 bg-red-50 p-3">
				<p class="text-sm text-red-700">{error}</p>
			</div>
		{/if}

		<!-- Actions -->
		<div class="flex items-center justify-end gap-3">
			<button
				onclick={handleClose}
				disabled={isUploading}
				class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:opacity-50"
			>
				Cancel
			</button>
			<button
				onclick={handleUpload}
				disabled={!selectedFile || !selectedCategory || isUploading}
				class="flex items-center gap-2 rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-green-700 disabled:opacity-50"
			>
				{#if isUploading}
					<Loader class="h-4 w-4 animate-spin" />
					Uploading...
				{:else}
					Upload Document
				{/if}
			</button>
		</div>
	</div>
</Modal>
