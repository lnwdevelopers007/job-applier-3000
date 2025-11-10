<script lang="ts">
	import { FileText, Download, BriefcaseBusiness, Globe, MapPin, StickyNote, Check, X } from 'lucide-svelte';
	import SkillTag from '$lib/components/job/SkillTag.svelte';
	import NotesModal from '$lib/components/ui/NotesModal.svelte';
	
	let { 
		userData = {},
		showApplicationInfo = false,
		isPreviewMode = false,
		appliedJobTitle = '',
		onAccept = () => {},
		onReject = () => {},
		onNotes = () => {},
		isUpdatingStatus = false,
		candidateStatus = '',
		candidateId = ''
	} = $props();
	
	// Notes modal state
	let showNotesModal = $state(false);
	
	function handleNotesClick() {
		showNotesModal = true;
		// Still call the parent onNotes function if provided
		onNotes();
	}
	
	// Skills from userData
	const skills = $derived(userData.skills && Array.isArray(userData.skills) ? userData.skills : []);
	
	// Documents from userData 
	const documents = $derived(userData.documents && Array.isArray(userData.documents) ? userData.documents : []);

	// Helper function to display value or dash
	const displayValue = (value?: string) => value && value.trim() ? value : '-';
	
	// Calculate age from date of birth
	const calculateAge = (dateOfBirth?: string) => {
		if (!dateOfBirth) return null;
		const today = new Date();
		const birthDate = new Date(dateOfBirth);
		let age = today.getFullYear() - birthDate.getFullYear();
		const monthDiff = today.getMonth() - birthDate.getMonth();
		if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
			age--;
		}
		return age;
	};
	
	const age = $derived(calculateAge(userData.dateOfBirth));
	
	
	// Helper to check if status allows actions
	const canTakeAction = $derived(candidateStatus === 'Pending');
	
	// Helper function to ensure URL has proper protocol
	function ensureHttps(url: string): string {
		if (!url) return '';
		const trimmedUrl = url.trim();
		if (trimmedUrl.startsWith('http://') || trimmedUrl.startsWith('https://')) {
			return trimmedUrl;
		}
		// Handle LinkedIn specific case - if it's just a username, build the full URL
		if (trimmedUrl.includes('linkedin.com/in/') && !trimmedUrl.startsWith('www.')) {
			return `https://www.${trimmedUrl}`;
		}
		return `https://${trimmedUrl}`;
	}
</script>

<div class="mx-auto p-8">
	<!-- Profile Header -->
	<div class="flex items-start justify-between mb-6 border-b border-gray-200 pb-6">
		<div class="flex items-center gap-4">
			<div class="w-16 h-16 rounded-full overflow-hidden">
				{#if userData.avatar}
					<img src={userData.avatar} alt={userData.fullName} class="w-16 h-16 object-cover" />
				{:else}
					<div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center">
						<span class="text-xl font-semibold text-gray-600">
							{userData.fullName ? userData.fullName.charAt(0) : 'T'}
						</span>
						
					</div>
				{/if}
			</div>
			<div>
				<h1 class="text-2xl font-semibold text-gray-900 mb-1">{displayValue(userData.fullName || userData.name)}</h1>
				<div class="flex items-center gap-3 text-sm text-gray-600 mb-2">
					<span>{displayValue(userData.desiredRole)}</span>
					<div class="flex items-center gap-1">
						<MapPin class="w-4 h-4" />
						<span>{displayValue(userData.location)}</span>
					</div>
				</div>
				
				<!-- Application Info -->
				{#if showApplicationInfo}
					<div class="flex items-center gap-3 text-sm text-gray-600 mt-2">
						<div class="flex items-center gap-2">
							<BriefcaseBusiness class="w-4 h-4" />
							<span>Applied for {appliedJobTitle || 'Your Job'}</span>
						</div>
					</div>
				{/if}
			</div>
		</div>
		
		<!-- Action buttons (only in non-preview mode) -->
		{#if !isPreviewMode && showApplicationInfo}
			<div class="flex gap-2">
				<button 
					class="flex items-center px-3 py-1.5 text-sm bg-white font-medium rounded-md border border-gray-200 hover:bg-gray-50"
					onclick={handleNotesClick}
				>
					<StickyNote class="w-4 h-4 mr-2" />Notes
				</button>
				{#if canTakeAction}
					<button 
						class="flex items-center px-3 py-1.5 text-sm bg-green-600 font-medium text-white rounded-md hover:bg-green-700 disabled:opacity-50"
						onclick={() => onAccept(candidateId)}
						disabled={isUpdatingStatus}
					>
						<Check class="w-4 h-4 mr-2" />Accept
					</button>
					<button 
						class="flex items-center px-3 py-1.5 text-sm bg-red-500 font-medium text-white rounded-md hover:bg-red-600 disabled:opacity-50"
						onclick={() => onReject(candidateId)}
						disabled={isUpdatingStatus}
					>
						<X class="w-4 h-4 mr-2" />Reject
					</button>
				{/if}
			</div>
		{/if}
	</div>

	<!-- Contact Information -->
	<div class="mb-10">
		<h2 class="text-md font-medium text-gray-900 mb-4 flex items-center gap-2">
			Contact Information
		</h2>
		
		<div class="grid grid-cols-2 gap-x-8 gap-y-4">
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">EMAIL</p>
				<p class="text-gray-900 text-sm">{displayValue(userData.email)}</p>
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">PHONE</p>
				<p class="text-gray-900 text-sm">{displayValue(userData.phone)}</p>
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">LOCATION</p>
				<p class="text-gray-900 text-sm">{displayValue(userData.location)}</p>
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">AGE</p>
				<p class="text-gray-900 text-sm">{age ? `${age} years old` : '-'}</p>
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">LINKEDIN</p>
				{#if userData.linkedin && userData.linkedin.trim() && userData.linkedin.trim() !== '-'}
					<!-- Debug: LinkedIn value is: {userData.linkedin} -->
					<a href={ensureHttps(userData.linkedin)} target="_blank" rel="noopener noreferrer" class="flex items-center gap-1 hover:text-blue-600 transition-colors">
						<svg class="h-4 w-4" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
						<p class="text-gray-900 text-sm underline">{userData.linkedin}</p>
					</a>
				{:else}
					<div class="flex items-center gap-1">
						<svg class="h-4 w-4" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
						<p class="text-gray-900 text-sm">-</p>
					</div>
				{/if}
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">GITHUB</p>
				{#if userData.github && userData.github.trim()}
					<a href={ensureHttps(userData.github)} target="_blank" rel="noopener noreferrer" class="flex items-center gap-1 hover:text-gray-800 transition-colors">
						<svg class="h-4 w-4" fill="#181717" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.30.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
						<p class="text-gray-900 text-sm underline">{userData.github}</p>
					</a>
				{:else}
					<div class="flex items-center gap-1">
						<svg class="h-4 w-4" fill="#181717" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
						<p class="text-gray-900 text-sm">-</p>
					</div>
				{/if}
			</div>
			<div>
				<p class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">PORTFOLIO WEBSITE</p>
				{#if userData.portfolio && userData.portfolio.trim()}
					<a href={ensureHttps(userData.portfolio)} target="_blank" rel="noopener noreferrer" class="flex items-center gap-1 hover:text-blue-600 transition-colors">
						<Globe class="w-4 h-4 text-gray-500" />
						<p class="text-gray-900 text-sm underline">{userData.portfolio}</p>
					</a>
				{:else}
					<div class="flex items-center gap-1">
						<Globe class="w-4 h-4 text-gray-500" />
						<p class="text-gray-900 text-sm">-</p>
					</div>
				{/if}
			</div>
		</div>
	</div>

	<!-- About Me -->
	{#if userData.aboutMe && userData.aboutMe.trim()}
		<div class="mb-10">
			<h2 class="text-md font-medium text-gray-900 mb-4 flex items-center gap-2">
				About Me
			</h2>
			
			<div class="text-sm text-gray-700 leading-relaxed">
				<p>{userData.aboutMe}</p>
			</div>
		</div>
	{/if}

	<!-- Skills -->
	<div class="mb-10">
		<h2 class="text-md font-medium text-gray-900 mb-4">
			Skills
		</h2>
		
		{#if skills.length > 0}
			<div class="flex flex-wrap gap-2">
				{#each skills as skill (skill)}
					<SkillTag {skill} />
				{/each}
			</div>
		{:else}
			<div class="text-center py-8 text-gray-500">
				<p class="text-sm">No skills listed</p>
			</div>
		{/if}
	</div>

	<!-- Documents -->
	<div class="mb-8">
		<h2 class="text-md font-medium text-gray-900 mb-4 flex items-center gap-2">
			Documents
		</h2>
		
		{#if documents.length > 0}
			<div class="space-y-3">
				{#each documents as doc (doc.name)}
					<div class="flex items-center justify-between p-3 bg-white border border-gray-200 rounded">
						<div class="flex items-center gap-3">
							<div class="w-8 h-8 bg-red-500 rounded flex items-center justify-center">
								<FileText class="w-4 h-4 text-white" />
							</div>
							<div>
								<p class="font-medium text-gray-900">{doc.name || 'Document'}</p>
								<p class="text-sm text-gray-500">{doc.size || 'Unknown size'}</p>
							</div>
						</div>
						<button class="p-2 text-gray-400 hover:text-gray-600 transition-colors">
							<Download class="w-4 h-4" />
						</button>
					</div>
				{/each}
			</div>
		{:else}
			<div class="text-center py-8 text-gray-500">
				<p class="text-sm">No documents uploaded</p>
			</div>
		{/if}
	</div>
</div>

<!-- Notes Modal -->
<NotesModal 
	bind:isOpen={showNotesModal}
	onClose={() => showNotesModal = false}
	candidateId={candidateId}
	candidateName={userData.fullName || userData.name || 'Unknown Candidate'}
	candidateAvatar={userData.avatar}
/>