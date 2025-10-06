<script lang="ts">
	import { onMount } from 'svelte';

	let selectedFile: File | null = null;
	let category = 'resume';
	let userID = ''; // Test user ObjectID
	let userRole = 'jobSeeker'; // or 'company'
	let uploadedFiles: any[] = [];
	let uploading = false;
	let error = '';
	let success = '';

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			selectedFile = target.files[0];
			error = '';
			success = '';

			// Validate file
			if (selectedFile.type !== 'application/pdf') {
				error = 'Only PDF files are allowed';
				selectedFile = null;
				return;
			}

			if (selectedFile.size > 10 * 1024 * 1024) {
				error = 'File size must be less than 10MB';
				selectedFile = null;
				return;
			}
		}
	}

	async function uploadFile() {
		if (!selectedFile || !userID) {
			error = 'Please select a file and enter a user ID';
			return;
		}

		uploading = true;
		error = '';
		success = '';

		const formData = new FormData();
		formData.append('file', selectedFile);
		formData.append('category', category);
		formData.append('userID', userID);
		formData.append('userRole', userRole);

		try {
			const response = await fetch('http://localhost:8080/files/upload', {
				method: 'POST',
				body: formData
			});

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Upload failed');
			}

			const result = await response.json();
			uploadedFiles = [...uploadedFiles, result];
			selectedFile = null;
			success = `File "${result.filename}" uploaded successfully!`;

			// Reset file input
			const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement;
			if (fileInput) fileInput.value = '';
		} catch (err: any) {
			error = err.message;
		} finally {
			uploading = false;
		}
	}

	async function downloadFile(fileId: string, filename: string) {
		try {
			const response = await fetch(
				`http://localhost:8080/files/download/${fileId}?requestingUserID=${userID}`
			);

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Download failed');
			}

			const blob = await response.blob();
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = filename;
			a.click();
			window.URL.revokeObjectURL(url);

			success = `File "${filename}" downloaded successfully!`;
		} catch (err: any) {
			error = err.message;
		}
	}

	async function loadUserFiles() {
		if (!userID) return;

		try {
			const response = await fetch(
				`http://localhost:8080/files/user/${userID}?requestingUserID=${userID}`
			);
			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Failed to load files');
			}

			const data = await response.json();
			uploadedFiles = data.files || [];
		} catch (err: any) {
			error = err.message;
		}
	}

	async function deleteFile(fileId: string, filename: string) {
		if (!confirm(`Are you sure you want to delete "${filename}"?`)) return;

		try {
			const response = await fetch(
				`http://localhost:8080/files/${fileId}?requestingUserID=${userID}`,
				{
					method: 'DELETE'
				}
			);

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Delete failed');
			}

			uploadedFiles = uploadedFiles.filter((f) => f.id !== fileId);
			success = `File "${filename}" deleted successfully!`;
		} catch (err: any) {
			error = err.message;
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
	}

	function getCategoryLabel(category: string): string {
		const labels: Record<string, string> = {
			resume: 'Resume',
			cover_letter: 'Cover Letter',
			certification: 'Certification',
			verification: 'Verification'
		};
		return labels[category] || category;
	}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="container mx-auto px-4 max-w-6xl">
		<h1 class="text-4xl font-bold mb-2 text-gray-800">File Upload Test Page</h1>
		<p class="text-gray-600 mb-8">Test MongoDB file storage with role-based validation</p>

		<!-- Alerts -->
		{#if error}
			<div class="bg-red-50 border-l-4 border-red-500 text-red-700 px-4 py-3 rounded mb-6">
				<p class="font-medium">Error</p>
				<p class="text-sm">{error}</p>
			</div>
		{/if}

		{#if success}
			<div class="bg-green-50 border-l-4 border-green-500 text-green-700 px-4 py-3 rounded mb-6">
				<p class="font-medium">Success</p>
				<p class="text-sm">{success}</p>
			</div>
		{/if}

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Upload Form -->
			<div class="bg-white shadow-lg rounded-lg overflow-hidden">
				<div class="bg-blue-600 text-white px-6 py-4">
					<h2 class="text-xl font-semibold">Upload File</h2>
				</div>
				<div class="p-6 space-y-4">
					<div>
						<label for="user-id-input" class="block text-gray-700 text-sm font-bold mb-2">
							User ID (ObjectID) *
						</label>
						<input
							id="user-id-input"
							type="text"
							bind:value={userID}
							placeholder="e.g., 507f1f77bcf86cd799439011"
							class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
						<p class="text-xs text-gray-500 mt-1">
							Enter a valid MongoDB ObjectID for testing
						</p>
					</div>

					<div>
						<label for="user-role-select" class="block text-gray-700 text-sm font-bold mb-2">
							User Role
						</label>
						<select
							id="user-role-select"
							bind:value={userRole}
							on:change={() => {
								// Reset category when role changes
								if (userRole === 'jobSeeker') {
									category = 'resume';
								} else {
									category = 'verification';
								}
							}}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							<option value="jobSeeker">Job Seeker</option>
							<option value="company">Company</option>
						</select>
					</div>

					<div>
						<label for="category-select" class="block text-gray-700 text-sm font-bold mb-2">Category</label>
						<select
							id="category-select"
							bind:value={category}
							class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							{#if userRole === 'jobSeeker'}
								<option value="resume">Resume</option>
								<option value="cover_letter">Cover Letter</option>
								<option value="certification">Certification</option>
							{:else}
								<option value="verification">Verification</option>
								<option value="certification">Certification</option>
							{/if}
						</select>
						<p class="text-xs text-gray-500 mt-1">
							{#if userRole === 'jobSeeker'}
								Job seekers can upload: Resume, Cover Letter, Certification
							{:else}
								Companies can upload: Verification, Certification
							{/if}
						</p>
					</div>

					<div>
						<label for="file-input" class="block text-gray-700 text-sm font-bold mb-2">
							Select PDF File (Max 10MB)
						</label>
						<input
							id="file-input"
							type="file"
							accept=".pdf,application/pdf"
							on:change={handleFileSelect}
							class="w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
						/>
					</div>

					{#if selectedFile}
						<div class="bg-gray-50 p-3 rounded text-sm text-gray-700">
							<p class="font-medium">Selected File:</p>
							<p>{selectedFile.name}</p>
							<p class="text-gray-500">{formatFileSize(selectedFile.size)}</p>
						</div>
					{/if}

					<button
						on:click={uploadFile}
						disabled={!selectedFile || !userID || uploading}
						class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-4 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{uploading ? 'Uploading...' : 'Upload File'}
					</button>

					<button
						on:click={loadUserFiles}
						disabled={!userID}
						class="w-full bg-gray-600 hover:bg-gray-700 text-white font-bold py-3 px-4 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Load My Files
					</button>
				</div>
			</div>

			<!-- Files List -->
			<div class="bg-white shadow-lg rounded-lg overflow-hidden">
				<div class="bg-green-600 text-white px-6 py-4">
					<h2 class="text-xl font-semibold">My Files ({uploadedFiles.length})</h2>
				</div>
				<div class="p-6">
					{#if uploadedFiles.length === 0}
						<div class="text-center py-8 text-gray-500">
							<svg
								class="mx-auto h-12 w-12 text-gray-400 mb-4"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
								/>
							</svg>
							<p class="text-lg mb-2">No files uploaded yet</p>
							<p class="text-sm">Upload a file or load your files</p>
						</div>
					{:else}
						<div class="space-y-3 max-h-[600px] overflow-y-auto">
							{#each uploadedFiles as file}
								<div
									class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
								>
									<div class="flex items-start justify-between mb-3">
										<div class="flex-1">
											<div class="flex items-center gap-2 mb-2">
												<svg
													class="h-5 w-5 text-red-500"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
													/>
												</svg>
												<p class="font-semibold text-gray-800">{file.filename}</p>
											</div>
											<div class="text-xs text-gray-600 space-y-1 ml-7">
												<div class="flex items-center gap-4">
													<span
														class="inline-block px-2 py-1 rounded text-xs font-medium bg-blue-100 text-blue-800"
													>
														{getCategoryLabel(file.category)}
													</span>
													<span>{formatFileSize(file.size)}</span>
												</div>
												<p>
													<span class="font-medium">Uploaded:</span>
													{new Date(file.uploadDate).toLocaleString()}
												</p>
												<p class="text-gray-400 truncate text-[10px]">
													<span class="font-medium">ID:</span>
													{file.id}
												</p>
											</div>
										</div>
									</div>
									<div class="flex gap-2">
										<button
											on:click={() => downloadFile(file.id, file.filename)}
											class="flex-1 bg-green-500 hover:bg-green-600 text-white text-sm font-bold py-2 px-3 rounded transition-colors flex items-center justify-center gap-2"
										>
											<svg
												class="h-4 w-4"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
												/>
											</svg>
											Download
										</button>
										<button
											on:click={() => deleteFile(file.id, file.filename)}
											class="flex-1 bg-red-500 hover:bg-red-600 text-white text-sm font-bold py-2 px-3 rounded transition-colors flex items-center justify-center gap-2"
										>
											<svg
												class="h-4 w-4"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												/>
											</svg>
											Delete
										</button>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Info Boxes -->
		<div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
			<!-- Testing Instructions -->
			<div class="bg-blue-50 border border-blue-200 rounded-lg p-5">
				<h3 class="font-semibold text-blue-900 mb-3 flex items-center gap-2">
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
				<ol class="text-sm text-blue-800 space-y-2 list-decimal list-inside">
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
			<div class="bg-purple-50 border border-purple-200 rounded-lg p-5">
				<h3 class="font-semibold text-purple-900 mb-3 flex items-center gap-2">
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
				<div class="text-sm text-purple-800 space-y-3">
					<div>
						<p class="font-semibold mb-1">👤 Job Seeker can upload:</p>
						<ul class="list-disc list-inside ml-2 space-y-1">
							<li>Resume</li>
							<li>Cover Letter</li>
							<li>Certification</li>
						</ul>
					</div>
					<div>
						<p class="font-semibold mb-1">🏢 Company can upload:</p>
						<ul class="list-disc list-inside ml-2 space-y-1">
							<li>Verification</li>
							<li>Certification</li>
						</ul>
					</div>
					<div class="bg-purple-100 p-2 rounded mt-2">
						<p class="font-medium">🔒 Authorization:</p>
						<p class="text-xs mt-1">Users can only access their own files</p>
					</div>
				</div>
			</div>
		</div>

		<!-- Technical Details -->
		<div class="mt-6 bg-gray-100 border border-gray-300 rounded-lg p-5">
			<h3 class="font-semibold text-gray-900 mb-3">Technical Implementation Details</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
				<div>
					<p class="font-medium text-gray-700 mb-1">Storage:</p>
					<ul class="text-gray-600 space-y-1">
						<li>• Binary storage ([]byte)</li>
						<li>• Max 10MB per file</li>
						<li>• PDF only</li>
					</ul>
				</div>
				<div>
					<p class="font-medium text-gray-700 mb-1">Schema:</p>
					<ul class="text-gray-600 space-y-1">
						<li>• Linked to users collection</li>
						<li>• Role-based validation</li>
						<li>• Metadata included</li>
					</ul>
				</div>
				<div>
					<p class="font-medium text-gray-700 mb-1">Security:</p>
					<ul class="text-gray-600 space-y-1">
						<li>• Owner-only access</li>
						<li>• Auth middleware ready</li>
						<li>• File type validation</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</div>