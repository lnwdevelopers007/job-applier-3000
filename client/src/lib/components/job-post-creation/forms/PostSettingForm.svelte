<script>
	let { 
		formData = $bindable(),
		validationErrors = {},
		showValidationErrors = false
	} = $props();
	
	if (!formData.applicationRequirements) {
		formData.applicationRequirements = {
			resume: false,
			transcript: false,
			portfolio: false,
			linkedin: false
		};
	}
	
	if (!formData.notificationSettings) {
		formData.notificationSettings = {
			notifyOnApply: false
		};
	}
</script>

<div class="space-y-6">	
	<!-- Posting Timeline -->
	<div>
		<h3 class="mb-4 text-sm font-medium text-gray-900">Posting Timeline</h3>
		<div class="grid grid-cols-2 gap-4">
			<div>
				<label for="postingOpenDate" class="mb-2 block text-sm text-gray-700">
					Posting Open Date <span class="text-red-500">*</span>
				</label>
				<input
					id="postingOpenDate"
					type="date"
					bind:value={formData.postingOpenDate}
					placeholder="DD/MM/YYYY"
					class="w-full rounded-md border px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-500 {showValidationErrors && validationErrors.postingOpenDate ? 'border-red-500' : 'border-gray-300'}"
				/>
				{#if showValidationErrors && validationErrors.postingOpenDate}
					<p class="mt-1 text-xs text-red-600">{validationErrors.postingOpenDate}</p>
				{:else}
					<p class="mt-1 text-xs text-gray-500">When should this job posting go live?</p>
				{/if}
			</div>
			
			<div>
				<label for="postingCloseDate" class="mb-2 block text-sm text-gray-700">
					Posting Close Date
				</label>
				<input
					id="postingCloseDate"
					type="date"
					bind:value={formData.postingCloseDate}
					placeholder="DD/MM/YYYY"
					class="w-full rounded-md border px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-500 {showValidationErrors && validationErrors.postingCloseDate ? 'border-red-500' : 'border-gray-300'}"
				/>
				{#if showValidationErrors && validationErrors.postingCloseDate}
					<p class="mt-1 text-xs text-red-600">{validationErrors.postingCloseDate}</p>
				{:else}
					<p class="mt-1 text-xs text-gray-500">When should applications close? (optional)</p>
				{/if}
			</div>
		</div>
	</div>

	<!-- Application Requirements -->
	<!-- <div>
		<h3 class="mb-4 text-sm font-medium text-gray-900">Application Requirements</h3>
		<div class="space-y-3">
			<label class="flex items-center">
				<input
					type="checkbox"
					bind:checked={formData.applicationRequirements.resume}
					class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
				/>
				<span class="ml-2 text-sm text-gray-700">Resume/CV required</span>
			</label>

			<label class="flex items-center">
				<input
					type="checkbox"
					bind:checked={formData.applicationRequirements.transcript}
					class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
				/>
				<span class="ml-2 text-sm text-gray-700">Transcript required</span>
			</label>

			<label class="flex items-center">
				<input
					type="checkbox"
					bind:checked={formData.applicationRequirements.portfolio}
					class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
				/>
				<span class="ml-2 text-sm text-gray-700">Portfolio/work samples required</span>
			</label>

			<label class="flex items-center">
				<input
					type="checkbox"
					bind:checked={formData.applicationRequirements.linkedin}
					class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
				/>
				<span class="ml-2 text-sm text-gray-700">LinkedIn profile required</span>
			</label>
		</div>
	</div> -->

	<!-- Screening Questions -->
	<!-- <div>
		<h3 class="mb-2 text-sm font-medium text-gray-900">Screening Questions</h3>
		<textarea
			bind:value={formData.screeningQuestions}
			placeholder="Add questions to help filter candidates."
			rows="3"
			class="w-full resize-none rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-500"
		></textarea>
		<p class="mt-1 text-xs text-gray-500">
			Optional questions to help screen applicants before interviews
		</p>
	</div> -->

	<!-- Notification Settings -->
	<div>
		<h3 class="mb-4 text-sm font-medium text-gray-900">Notification Settings</h3>
		<div class="space-y-3">
			<label class="flex items-center">
				<input
					type="checkbox"
					bind:checked={formData.emailNotifications}
					class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-green-500"
				/>
				<span class="ml-2 text-sm text-gray-700">Notify me when candidates apply</span>
			</label>
		</div>
	</div>
</div>