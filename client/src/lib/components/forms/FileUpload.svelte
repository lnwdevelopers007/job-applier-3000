<script>
	let {
		label = '',
		currentImage = null,
		accept = 'image/*',
		maxSize = '2MB',
		helpText = 'PNG, JPG or SVG',
		onFileSelect = () => {}
	} = $props();
	

	/**
	 * @type {HTMLInputElement}
	 */
	let fileInput;
	let previewImage = $state(null);
	
	// Display priority: previewImage (newly selected) > currentImage (existing)
	let displayImage = $derived(previewImage || currentImage);

	// @ts-expect-error: event target may be undefined, safe to ignore
	function handleFileSelect(event) {
		const file = event.target.files[0];
		if (file) {
			// Create preview URL for the selected file
			const reader = new FileReader();
			reader.onload = (e) => {
				previewImage = e.target.result;
			};
			reader.readAsDataURL(file);
			
			onFileSelect(file);
		}
	}
</script>

<div>
	{#if label}
		<!-- svelte-ignore a11y_label_has_associated_control -->
		<label class="text-sm font-medium text-gray-700 block mb-2">{label}</label>
	{/if}
	
	<div class="flex items-center gap-4">
		{#if displayImage}
			<img 
				src={displayImage} 
				alt="Current" 
				class="w-16 h-16 rounded-full object-cover"
			/>
		{:else}
			<div class="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center border-2 border-gray-200">
				<svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
				</svg>
			</div>
		{/if}
		
		<div>
			<button 
				onclick={() => fileInput.click()}
				type="button"
				class="text-sm font-medium text-gray-700 hover:text-gray-900 hover:cursor-pointer"
			>
				Click to upload
			</button>
			<span class="text-sm text-gray-500"> or drag and drop</span>
			<p class="text-xs text-gray-500 mt-1">{helpText} (max. {maxSize})</p>
			<input 
				bind:this={fileInput}
				type="file" 
				class="hidden" 
				{accept}
				onchange={handleFileSelect}
			/>
		</div>
	</div>
</div>