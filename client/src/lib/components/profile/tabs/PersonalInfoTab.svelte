<script>
	import TextInput from '$lib/components/forms/TextInput.svelte';
	// import Select from '$lib/components/forms/Select.svelte';
	import Textarea from '$lib/components/forms/Textarea.svelte';
	import { Globe } from 'lucide-svelte';
	
	let {
		profileData = $bindable({})
	} = $props();
	
	let skillInput = $state('');
	
	// Initialize skills array if it doesn't exist
	if (!profileData.skills) {
		profileData.skills = [];
	}
	
	function addSkill(event) {
		if (event.key === 'Enter' && skillInput.trim()) {
			event.preventDefault();
			if (!profileData.skills.includes(skillInput.trim())) {
				profileData.skills = [...profileData.skills, skillInput.trim()];
			}
			skillInput = '';
		}
	}
	
	function removeSkill(skill) {
		profileData.skills = profileData.skills.filter(s => s !== skill);
	}
</script>

<div class="divide-y divide-gray-200">
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Full name</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.fullName}
					placeholder="Enter your full name"
				/>
			</div>
		</div>
	</div>

	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Date of birth</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.dateOfBirth}
					type="date"
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Location</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.location}
					placeholder="City, State/Country"
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Desired role</label>
				<p class="text-xs text-gray-500 mt-1">The type of role you're looking for</p>
			</div>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.desiredRole}
					placeholder="e.g., Product Designer, UX/UI Designer"
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">About me</label>
				<p class="text-xs text-gray-500 mt-1">Write a short introduction</p>
			</div>
			<div class="col-span-1">
				<Textarea 
					bind:value={profileData.aboutMe}
					placeholder="Tell us about yourself..."
					maxLength={500}
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Phone number</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.phone}
					type="tel"
					placeholder="Enter your phone number"
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">LinkedIn profile</label>
				<p class="text-xs text-gray-500 mt-1">Your professional profile URL</p>
			</div>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.linkedin}
					type="url"
					placeholder="linkedin.com/in/your-profile"
					icon='<svg class="h-5 w-5" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>'
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Portfolio website</label>
				<p class="text-xs text-gray-500 mt-1">Showcase your work</p>
			</div>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.portfolio}
					type="url"
					placeholder="yourwebsite.com"
					iconComponent={Globe}
				/>
			</div>
		</div>
	</div>
	
	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">GitHub profile</label>
				<p class="text-xs text-gray-500 mt-1">Your code repositories</p>
			</div>
			<div class="col-span-1">
				<TextInput 
					bind:value={profileData.github}
					type="url"
					placeholder="github.com/username"
					icon='<svg class="h-5 w-5" fill="#181717" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.30.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>'
				/>
			</div>
		</div>
	</div>

	<div class="py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Skills</label>
				<p class="text-xs text-gray-500 mt-1">Your technical and professional skills</p>
			</div>
			<div class="col-span-1">
				<div class="relative">
					<input
						type="text"
						bind:value={skillInput}
						onkeydown={addSkill}
						placeholder="Type a skill and press Enter to add"
						class="w-full px-4 py-2 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all"
					/>
					
					<!-- Skill Tags -->
					{#if profileData.skills && profileData.skills.length > 0}
						<div class="flex flex-wrap gap-2 mt-3">
							{#each profileData.skills as skill, index (index)}
								<span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm bg-gray-100 text-gray-700">
									{skill}
									<button
										type="button"
										onclick={() => removeSkill(skill)}
										class="ml-1 text-gray-500 hover:text-gray-700 focus:outline-none"
									>
										<svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
										</svg>
									</button>
								</span>
							{/each}
						</div>
					{/if}
				</div>
				<p class="mt-2 text-xs text-gray-500">Press Enter to add each skill</p>
			</div>
		</div>
	</div>
</div>