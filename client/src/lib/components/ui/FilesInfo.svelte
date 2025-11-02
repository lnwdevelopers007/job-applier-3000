<!-- Info Boxes -->
<div class="mt-8 grid grid-cols-1 gap-6 md:grid-cols-2">
	<!-- Testing Instructions -->
	<div class="rounded-lg border border-blue-200 bg-blue-50 p-5">
		<h3 class="mb-3 flex items-center gap-2 font-semibold text-blue-900">
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			Testing Instructions
		</h3>
		<ol class="list-inside list-decimal space-y-2 text-sm text-blue-800">
			<li>Enter a valid MongoDB ObjectID for the User ID field</li>
			<li>Select user role (Job Seeker or Company)</li>
			<li>Available categories will update based on role</li>
			<li>Select a PDF file (max 10MB)</li>
			<li>Click "Upload File"</li>
			<li>Test downloading to verify file integrity</li>
			<li>Try uploading wrong category for role (should fail)</li>
			<li>Use "Load My Files" to query files by user</li>
		</ol>
	</div>

	<!-- Role-Based Rules -->
	<div class="rounded-lg border border-purple-200 bg-purple-50 p-5">
		<h3 class="mb-3 flex items-center gap-2 font-semibold text-purple-900">
			<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
				/>
			</svg>
			Role-Based Validation Rules
		</h3>
		<div class="space-y-3 text-sm text-purple-800">
			<div>
				<p class="mb-1 font-semibold">👤 Job Seeker can upload:</p>
				<ul class="ml-2 list-inside list-disc space-y-1">
					<li>Resume</li>
					<li>Cover Letter</li>
					<li>Certification</li>
				</ul>
			</div>
			<div>
				<p class="mb-1 font-semibold">🏢 Company can upload:</p>
				<ul class="ml-2 list-inside list-disc space-y-1">
					<li>Verification</li>
					<li>Certification</li>
				</ul>
			</div>
			<div class="mt-2 rounded bg-purple-100 p-2">
				<p class="font-medium">🔒 Authorization:</p>
				<p class="mt-1 text-xs">Users can only access their own files</p>
			</div>
		</div>
	</div>
</div>

<!-- Technical Details -->
<div class="mt-6 rounded-lg border border-gray-300 bg-gray-100 p-5">
	<h3 class="mb-3 font-semibold text-gray-900">Technical Implementation Details</h3>
	<div class="grid grid-cols-1 gap-4 text-sm md:grid-cols-3">
		<div>
			<p class="mb-1 font-medium text-gray-700">Storage:</p>
			<ul class="space-y-1 text-gray-600">
				<li>• Binary storage ([]byte)</li>
				<li>• Max 10MB per file</li>
				<li>• PDF only</li>
			</ul>
		</div>
		<div>
			<p class="mb-1 font-medium text-gray-700">Schema:</p>
			<ul class="space-y-1 text-gray-600">
				<li>• Linked to users collection</li>
				<li>• Role-based validation</li>
				<li>• Metadata included</li>
			</ul>
		</div>
		<div>
			<p class="mb-1 font-medium text-gray-700">Security:</p>
			<ul class="space-y-1 text-gray-600">
				<li>• Owner-only access</li>
				<li>• Auth middleware ready</li>
				<li>• File type validation</li>
			</ul>
		</div>
	</div>
</div>
