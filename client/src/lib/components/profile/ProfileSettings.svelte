<script>
	import TextInput from '$lib/components/forms/TextInput.svelte';
	import Select from '$lib/components/forms/Select.svelte';
	import Textarea from '$lib/components/forms/Textarea.svelte';
	import FileUpload from '$lib/components/forms/FileUpload.svelte';
	
	let {
		tabs = [],
		activeTab = $bindable(''),
		userData = $bindable({}),
		onSave = () => {},
		title = 'Settings',
		userType = 'seeker'
	} = $props();
	
	let hasChanges = $state(false);
	let originalData = $state(null);
	let documentsChanged = $state(false);
	let originalDocuments = $state(null);
	
	// Initialize original data only once when component mounts
	$effect(() => {
		if (originalData === null) {
			originalData = JSON.parse(JSON.stringify(userData));
		}
		if (originalDocuments === null) {
			originalDocuments = JSON.parse(JSON.stringify(userData.documents || []));
		}
	});
	
	function handleTabChange(tabId) {
		activeTab = tabId;
	}
	
	// Use reactive computed values for change detection
	$effect(() => {
		if (originalData !== null) {
			hasChanges = JSON.stringify(userData) !== JSON.stringify(originalData);
		}
	});
	
	$effect(() => {
		if (originalDocuments !== null) {
			const currentDocs = userData.documents || [];
			documentsChanged = JSON.stringify(currentDocs) !== JSON.stringify(originalDocuments);
		}
	});
	
	function handleDocumentUpload(event) {
		const files = event.target.files;
		if (files && files.length > 0) {
			if (!userData.documents) {
				userData.documents = [];
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
			
			userData.documents = [...userData.documents, ...newDocs];
			event.target.value = '';
		}
	}
	
	function handleSave() {
		onSave(userData);
		originalData = JSON.parse(JSON.stringify(userData));
		hasChanges = false;
	}
	
	function handleSaveDocuments() {
		onSave({ documents: userData.documents });
		originalDocuments = JSON.parse(JSON.stringify(userData.documents || []));
		documentsChanged = false;
	}
	
	const activeTabData = $derived(tabs.find(tab => tab.id === activeTab));
</script>

<div class="bg-white rounded-lg border border-gray-200 p-6 min-h-[800px]">
	<h1 class="text-2xl font-semibold text-gray-900 p-3">{title}</h1>
	
	<div class="border-b border-gray-200">
		<nav class="flex -mb-px">
			{#each tabs as tab}
				<button
					onclick={() => handleTabChange(tab.id)}
					class="px-4 py-3 text-sm font-medium border-b-2 transition-all cursor-pointer {activeTab === tab.id ? 'border-green-600 text-green-700' : 'border-transparent text-gray-600 hover:text-gray-900 hover:border-gray-300'}"
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>
	
	<div class="flex justify-between items-start px-8 py-6 border-b border-gray-200">
		<div>
			<h2 class="text-lg font-medium text-gray-900">{activeTabData?.title || ''}</h2>
			<p class="text-sm text-gray-500 mt-1">{activeTabData?.description || ''}</p>
		</div>
		
		<div class="flex items-center gap-3">
			{#if activeTab === 'documents'}
				<button 
					onclick={() => document.getElementById('file-upload').click()}
					class="px-4 py-2 border border-gray-300 text-gray-700 font-medium text-sm rounded-lg hover:bg-gray-50 transition-colors flex items-center gap-2 cursor-pointer"
				>
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
					</svg>
					Add document
				</button>
			{/if}
			
			{#if activeTab === 'documents' && documentsChanged}
				<button
					onclick={handleSaveDocuments}
					class="px-4 py-2 bg-green-600 text-white font-medium text-sm rounded-lg hover:bg-green-700 transition-colors cursor-pointer"
				>
					Save changes
				</button>
			{:else if hasChanges && !activeTabData?.hideActions && activeTab !== 'documents'}
				<button
					onclick={handleSave}
					class="px-4 py-2 bg-green-600 text-white font-medium text-sm rounded-lg hover:bg-green-700 transition-colors cursor-pointer"
				>
					Save changes
				</button>
			{/if}
		</div>
	</div>
	
	<div class="tab-content">
		{#if activeTab === 'user'}
			{@render userTabContent()}
		{:else if activeTab === 'personal' && userType === 'seeker'}
			{@render personalTabContent()}
		{:else if activeTab === 'company' && userType === 'company'}
			{@render companyTabContent()}
		{:else if activeTab === 'documents'}
			{@render documentsTabContent()}
		{:else if activeTab === 'benefits' && userType === 'company'}
			{@render benefitsTabContent()}
		{/if}
	</div>
</div>

{#snippet userTabContent()}
	<div class="divide-y divide-gray-200">
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Full name</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.fullName}
						placeholder="Full name"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Email address</label>
				<div class="col-span-1">
					<TextInput 
						value={userData.email}
						type="email"
						readonly={true}
						icon='<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>'
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Profile Picture</label>
					<p class="text-xs text-gray-500 mt-1">This will be displayed on your profile.</p>
				</div>
				<div class="col-span-1">
					<FileUpload 
						currentImage={userData.avatar}
						onFileSelect={(file) => {
							userData.avatarFile = file;
						}}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Connected accounts</label>
					<p class="text-xs text-gray-500 mt-1">Manage your connected social accounts</p>
				</div>
				<div class="col-span-1">
					<div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 bg-white border border-gray-200 rounded-full flex items-center justify-center">
								<svg class="w-5 h-5" viewBox="0 0 24 24">
									<path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
									<path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
									<path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
									<path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
								</svg>
							</div>
							<div>
								<p class="text-sm font-medium text-gray-900">Google</p>
								<p class="text-xs text-gray-500">{userData.email}</p>
							</div>
						</div>
						<button 
							onclick={() => {
								userData.googleConnected = !userData.googleConnected;
							}}
							class="px-3 py-1.5 text-xs font-medium border rounded-md transition-colors cursor-pointer {userData.googleConnected ? 'text-red-600 border-red-200 hover:bg-red-50' : 'text-green-600 border-green-200 hover:bg-green-50'}"
						>
							{userData.googleConnected ? 'Disconnect' : 'Connect'}
						</button>
					</div>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Password</label>
					<p class="text-xs text-gray-500 mt-1">Update your login password</p>
				</div>
				<div class="col-span-1 space-y-4">
					<TextInput 
						label="Current password"
						type="password"
						bind:value={userData.currentPassword}
						placeholder="Enter current password"
					/>
					
					<TextInput 
						label="New password"
						type="password"
						bind:value={userData.newPassword}
						placeholder="Enter new password"
						helpText="Must be at least 8 characters."
					/>
					
					<TextInput 
						label="Confirm password"
						type="password"
						bind:value={userData.confirmPassword}
						placeholder="Confirm new password"
					/>
				</div>
			</div>
		</div>
	</div>
{/snippet}

{#snippet personalTabContent()}
	<div class="divide-y divide-gray-200">
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Location</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.location}
						placeholder="Melbourne, Australia"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Desired role</label>
					<p class="text-xs text-gray-500 mt-1">The type of role you're looking for</p>
				</div>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.desiredRole}
						placeholder="e.g., Product Designer, UX/UI Designer"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">About me</label>
					<p class="text-xs text-gray-500 mt-1">Write a short introduction</p>
				</div>
				<div class="col-span-1">
					<Textarea 
						bind:value={userData.aboutMe}
						placeholder="Tell us about yourself..."
						maxLength={500}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Phone number</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.phone}
						type="tel"
						placeholder="+1 (555) 000-0000"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Date of birth</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.dateOfBirth}
						type="date"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Gender</label>
				<div class="col-span-1">
					<Select 
						bind:value={userData.gender}
						placeholder="Select gender..."
						options={['Male', 'Female', 'Non-binary', 'Prefer not to say']}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">LinkedIn profile</label>
					<p class="text-xs text-gray-500 mt-1">Your professional profile URL</p>
				</div>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.linkedin}
						type="url"
						placeholder="linkedin.com/in/your-profile"
						icon='<svg class="h-5 w-5" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>'
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Portfolio website</label>
					<p class="text-xs text-gray-500 mt-1">Showcase your work</p>
				</div>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.portfolio}
						type="url"
						placeholder="yourwebsite.com"
						icon='<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" /></svg>'
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">GitHub profile</label>
					<p class="text-xs text-gray-500 mt-1">Your code repositories</p>
				</div>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.github}
						type="url"
						placeholder="github.com/username"
						icon='<svg class="h-5 w-5" fill="#181717" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>'
					/>
				</div>
			</div>
		</div>
	</div>
{/snippet}

{#snippet companyTabContent()}
	<div class="divide-y divide-gray-200">
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Company name</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.companyName}
						placeholder="Company name"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Company logo</label>
					<p class="text-xs text-gray-500 mt-1">This will be displayed on job postings and your profile</p>
				</div>
				<div class="col-span-1">
					<FileUpload 
						currentImage={userData.companyLogo}
						maxSize="2MB"
						onFileSelect={(file) => {
							userData.companyLogoFile = file;
						}}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Industry</label>
				<div class="col-span-1">
					<Select 
						bind:value={userData.industry}
						placeholder="Select industry..."
						options={[
							'Technology',
							'Healthcare',
							'Finance',
							'Education',
							'Manufacturing',
							'Retail',
							'Consulting',
							'Other'
						]}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Company size</label>
				<div class="col-span-1">
					<Select 
						bind:value={userData.companySize}
						placeholder="Select size..."
						options={[
							'1-10 employees',
							'11-50 employees',
							'51-200 employees',
							'201-500 employees',
							'501-1000 employees',
							'1000+ employees'
						]}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Founded</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.foundedYear}
						type="number"
						placeholder="Year founded"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Headquarters</label>
					<p class="text-xs text-gray-500 mt-1">Primary company location</p>
				</div>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.headquarters}
						placeholder="City, State/Country"
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">About us</label>
					<p class="text-xs text-gray-500 mt-1">Brief overview of your company and what you do</p>
				</div>
				<div class="col-span-1">
					<Textarea 
						bind:value={userData.aboutCompany}
						placeholder="Tell candidates about your company..."
						maxLength={1000}
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">Website</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.companyWebsite}
						type="url"
						placeholder="https://yourcompany.com"
						icon='<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" /></svg>'
					/>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-center">
				<label class="text-sm font-medium text-gray-700">LinkedIn</label>
				<div class="col-span-1">
					<TextInput 
						bind:value={userData.companyLinkedin}
						type="url"
						placeholder="linkedin.com/company/yourcompany"
						icon='<svg class="h-5 w-5" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>'
					/>
				</div>
			</div>
		</div>
	</div>
{/snippet}

{#snippet documentsTabContent()}
	<div class="divide-y divide-gray-200">
		<div class="px-8 py-5">
			<div class="grid grid-cols-4 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Documents</label>
					<p class="text-xs text-gray-500 mt-1">Upload documents for recruiters such as resume, cover letter, portfolio, transcripts, or certificates</p>
				</div>
				<div class="col-span-2">
					<div class="space-y-4">
						{#if userData.documents && userData.documents.length > 0}
							{#each userData.documents as doc, index}
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
												if (!userData.documents) userData.documents = [];
												const newDocs = userData.documents.filter((_, i) => i !== index);
												userData.documents = [...newDocs];
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
							<div class="text-center py-12">
								<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
								</svg>
								<h3 class="mt-4 text-sm font-medium text-gray-900">No documents</h3>
								<p class="mt-2 text-sm text-gray-500">Get started by uploading your first document.</p>
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
{/snippet}

{#snippet benefitsTabContent()}
	<div class="divide-y divide-gray-200">
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Health & wellness</label>
					<p class="text-xs text-gray-500 mt-1">Medical, dental, and wellness benefits</p>
				</div>
				<div class="col-span-1">
					<div class="space-y-3">
						{#each ['Health insurance', 'Dental insurance', 'Vision insurance', 'Gym membership'] as benefit}
							<label class="flex items-center">
								<input 
									type="checkbox" 
									checked={userData.benefits?.includes(benefit)}
									onchange={(e) => {
										if (!userData.benefits) userData.benefits = [];
										if (e.target.checked) {
											userData.benefits = [...userData.benefits, benefit];
										} else {
											userData.benefits = userData.benefits.filter(b => b !== benefit);
										}
									}}
									class="mr-2 text-green-600 rounded"
								/>
								<span class="text-sm text-gray-700">{benefit}</span>
							</label>
						{/each}
					</div>
				</div>
			</div>
		</div>
		
		<div class="px-8 py-5">
			<div class="grid grid-cols-3 gap-8 items-start">
				<div>
					<label class="text-sm font-medium text-gray-700">Work-life balance</label>
					<p class="text-xs text-gray-500 mt-1">Flexibility and time-off benefits</p>
				</div>
				<div class="col-span-1">
					<div class="space-y-3">
						{#each ['Remote work', 'Flexible hours', 'Unlimited PTO', 'Parental leave'] as benefit}
							<label class="flex items-center">
								<input 
									type="checkbox" 
									checked={userData.benefits?.includes(benefit)}
									onchange={(e) => {
										if (!userData.benefits) userData.benefits = [];
										if (e.target.checked) {
											userData.benefits = [...userData.benefits, benefit];
										} else {
											userData.benefits = userData.benefits.filter(b => b !== benefit);
										}
									}}
									class="mr-2 text-green-600 rounded"
								/>
								<span class="text-sm text-gray-700">{benefit}</span>
							</label>
						{/each}
					</div>
				</div>
			</div>
		</div>
	</div>
{/snippet}