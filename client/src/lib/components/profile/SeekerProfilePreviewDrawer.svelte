<script lang="ts">
	import { FileText, Download, Eye, BriefcaseBusiness, Globe, MapPin } from 'lucide-svelte';
	import SideDrawer from '$lib/components/ui/SideDrawer.svelte';
	import { fileService } from '$lib/services/fileService';

	let { isOpen = $bindable(false), userData = {} } = $props();

	// Skills from userData or fallback to default skills if empty
	const skills = $derived(() => {
		if (userData.skills && Array.isArray(userData.skills) && userData.skills.length > 0) {
			return userData.skills;
		}
		return [
			'JavaScript',
			'TypeScript',
			'React',
			'Node.js',
			'Python',
			'PostgreSQL',
			'MongoDB',
			'Docker',
			'AWS',
			'Git'
		];
	});

	// Documents from userData
	const documents = $derived(() => {
		if (userData.documents && Array.isArray(userData.documents) && userData.documents.length > 0) {
			return userData.documents;
		}
		return [];
	});

	// Helper function to display value or dash
	const displayValue = (value?: string) => (value && value.trim() ? value : '-');

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

	// Get category color and label
	const getCategoryColor = (category: string) => fileService.getCategoryColor(category);
	const getCategoryLabel = (category: string) => fileService.getCategoryLabel(category);
	const formatFileSize = (size: number | string) => {
		if (typeof size === 'number') {
			return fileService.formatFileSize(size);
		}
		return size || 'Unknown size';
	};
</script>

<SideDrawer
	bind:isOpen
	title="Profile Preview"
	subtitle="This is how your profile appears to employers"
	width="820px"
>
	<div class="shadow-xs mx-auto max-w-3xl rounded-lg border border-gray-200 bg-white p-8">
		<!-- Profile Header -->
		<div class="mb-6 flex items-start justify-between border-b border-gray-200 pb-6">
			<div class="flex items-center gap-4">
				<div class="h-16 w-16 overflow-hidden rounded-full">
					{#if userData.avatar}
						<img src={userData.avatar} alt={userData.fullName} class="h-16 w-16 object-cover" />
					{:else}
						<div class="flex h-16 w-16 items-center justify-center rounded-full bg-gray-200">
							<span class="text-xl font-semibold text-gray-600">
								{userData.fullName ? userData.fullName.charAt(0) : 'T'}
							</span>
						</div>
					{/if}
				</div>
				<div>
					<h1 class="mb-1 text-2xl font-semibold text-gray-900">
						{displayValue(userData.fullName || userData.name)}
					</h1>
					<div class="mb-2 flex items-center gap-3 text-sm text-gray-600">
						<span>{displayValue(userData.desiredRole)}</span>
						<div class="flex items-center gap-1">
							<MapPin class="h-3 w-3" />
							<span>{displayValue(userData.location)}</span>
						</div>
					</div>

					<!-- Application Info -->
					<div class="mt-2 flex items-center gap-2 text-sm text-gray-600">
						<BriefcaseBusiness class="h-4 w-4" />
						<span>Applied for Your Job 4 hours ago</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Contact Information -->
		<div class="mb-10">
			<h2 class="text-md mb-4 flex items-center gap-2 font-medium text-gray-900">
				Contact Information
			</h2>

			<div class="grid grid-cols-2 gap-x-8 gap-y-4">
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">EMAIL</p>
					<p class="text-sm text-gray-900">{displayValue(userData.email)}</p>
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">PHONE</p>
					<p class="text-sm text-gray-900">{displayValue(userData.phone)}</p>
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">LOCATION</p>
					<p class="text-sm text-gray-900">{displayValue(userData.location)}</p>
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">AGE</p>
					<p class="text-sm text-gray-900">{age ? `${age} years old` : '-'}</p>
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">LINKEDIN</p>
					{#if userData.linkedin && userData.linkedin.trim()}
						<div class="flex items-center gap-1">
							<svg class="h-4 w-4" fill="#0A66C2" viewBox="0 0 24 24"
								><path
									d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"
								/></svg
							>
							<p class="text-sm text-gray-900">{userData.linkedin}</p>
						</div>
					{:else}
						<p class="text-sm text-gray-900">-</p>
					{/if}
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">GITHUB</p>
					{#if userData.github && userData.github.trim()}
						<div class="flex items-center gap-1">
							<svg class="h-4 w-4" fill="#181717" viewBox="0 0 24 24"
								><path
									d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.30.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
								/></svg
							>
							<p class="text-sm text-gray-900">{userData.github}</p>
						</div>
					{:else}
						<p class="text-sm text-gray-900">-</p>
					{/if}
				</div>
				<div>
					<p class="mb-1 text-xs font-medium uppercase tracking-wide text-gray-500">
						PORTFOLIO WEBSITE
					</p>
					{#if userData.portfolio && userData.portfolio.trim()}
						<div class="flex items-center gap-1">
							<Globe class="h-4 w-4 text-gray-500" />
							<p class="text-sm text-gray-900">{userData.portfolio}</p>
						</div>
					{:else}
						<p class="text-sm text-gray-900">-</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- About Me -->
		{#if userData.aboutMe && userData.aboutMe.trim()}
			<div class="mb-10">
				<h2 class="text-md mb-4 flex items-center gap-2 font-medium text-gray-900">About Me</h2>

				<div class="text-sm leading-relaxed text-gray-700">
					<p>{userData.aboutMe}</p>
				</div>
			</div>
		{/if}

		<!-- Skills -->
		<div class="mb-10">
			<h2 class="text-md mb-4 font-medium text-gray-900">Skills</h2>

			<div class="flex flex-wrap gap-2">
				{#each skills() as skill (skill)}
					<span class="rounded border border-gray-200 bg-gray-100 px-3 py-1 text-sm text-gray-700">
						{skill}
					</span>
				{/each}
			</div>
		</div>

		<!-- Documents -->
		<div class="mb-8">
			<h2 class="text-md mb-4 flex items-center gap-2 font-medium text-gray-900">Documents</h2>

			{#if documents().length > 0}
				<div class="space-y-3">
					{#each documents() as doc (doc.id || doc.name)}
						<div
							class="flex items-center justify-between rounded border border-gray-200 bg-white p-3"
						>
							<div class="flex min-w-0 flex-1 items-center gap-3">
								<!-- Icon with category color -->
								<div
									class="h-10 w-10 bg-{getCategoryColor(
										doc.category
									)}-100 flex flex-shrink-0 items-center justify-center rounded"
								>
									<FileText class="h-5 w-5 text-{getCategoryColor(doc.category)}-600" />
								</div>

								<div class="min-w-0 flex-1">
									<div class="mb-0.5 flex items-center gap-2">
										<p class="truncate font-medium text-gray-900">
											{doc.filename || doc.name || 'Document'}
										</p>
										<!-- Category badge -->
										<span
											class="px-2 py-0.5 text-xs font-medium bg-{getCategoryColor(
												doc.category
											)}-100 text-{getCategoryColor(doc.category)}-700 flex-shrink-0 rounded"
										>
											{getCategoryLabel(doc.category)}
										</span>
									</div>
									<p class="text-sm text-gray-500">{formatFileSize(doc.size)}</p>
								</div>
							</div>

							<!-- Action buttons (disabled for preview) -->
							<div class="flex items-center gap-2">
								<button
									class="cursor-not-allowed p-2 text-gray-300"
									disabled
									title="Preview not available in preview mode"
								>
									<Eye class="h-4 w-4" />
								</button>
								<button
									class="cursor-not-allowed p-2 text-gray-300"
									disabled
									title="Download not available in preview mode"
								>
									<Download class="h-4 w-4" />
								</button>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<div class="py-8 text-center text-gray-500">
					<p class="text-sm">No documents uploaded</p>
				</div>
			{/if}
		</div>
	</div>
</SideDrawer>
