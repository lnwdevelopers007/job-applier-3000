<script lang="ts">
	import SideDrawer from '$lib/components/ui/SideDrawer.svelte';
	import JobDetailCard from '$lib/components/job/JobDetailCard.svelte';
	import JobCard from '../job/JobCard.svelte';
	import { SettingsService, type CompanyUserData } from '$lib/services/settingsService';
	
	let { 
		isOpen = $bindable(false),
		formData = {}
	} = $props();
	
	let activeTab = $state('detail'); // 'detail' or 'search'
	let companyData = $state<CompanyUserData | null>(null);
	// let loadingCompany = $state(false);
	
	function formatSalary() {
		if (!formData.minSalary || !formData.maxSalary) return 'Not specified';
		const currency = formData.currency || 'THB';
		return `${currency} ${formData.minSalary.toLocaleString()}-${formData.maxSalary.toLocaleString()}`;
	}

	// Load company data when drawer opens
	$effect(() => {
		if (isOpen && !companyData) {
			loadCompanyData();
		}
	});

	async function loadCompanyData() {
		try {
			// loadingCompany = true;
			companyData = await SettingsService.loadUserData<CompanyUserData>('company');
		} catch (error) {
			console.error('Failed to load company data:', error);
			// Use default data if loading fails
			companyData = null;
		} finally {
			// loadingCompany = false;
		}
	}

	// Transform formData to match Job interface for JobDetailCard
	const previewJob = $derived({
		id: 'preview',
		title: formData.jobTitle || 'Job Title',
		company: companyData?.companyName || 'Your Company Name',
		logo: companyData?.avatar || companyData?.companyLogo || '',
		location: formData.location || companyData?.headquarters || 'Bangkok, Thailand',
		workType: formData.workArrangement || 'on-site',
		workArrangement: formData.workType || 'full-time',
		salary: formatSalary(),
		posted: formData.postingOpenDate || new Date().toISOString(),
		closeDate: formData.postingCloseDate || '',
		description: formData.jobDescription || 'No job description provided',
		tags: formData.requiredSkills || []
	});

	// Transform company data for JobDetailCard
	const previewCompany = $derived({
		name: companyData?.companyName || 'Your Company Name',
		logo: companyData?.avatar || companyData?.companyLogo || '',
		location: companyData?.headquarters || formData.location || 'Bangkok, Thailand',
		size: companyData?.companySize || '1-10 employees',
		industry: companyData?.industry || 'Technology',
		aboutUs: companyData?.aboutCompany || 'This is where your company description will appear. You can update your company information in your company profile settings.',
		website: companyData?.companyWebsite || '',
		linkedin: companyData?.companyLinkedin || '',
		foundedYear: companyData?.foundedYear || ''
	});

	function handlePreviewApply() {
		// Disabled for preview
	}
</script>

<SideDrawer 
	bind:isOpen 
	title="Preview" 
	subtitle="This is a preview of what your job post will look like to appliers"
	width="820px"
>
	<!-- Tab Navigation -->
	<div class="flex gap-4 mb-3 pl-6">
		<div class="border-b border-gray-200">
			<button 
				onclick={() => activeTab = 'detail'}
				class="pb-3 px-1 text-sm font-medium transition-colors {activeTab === 'detail' ? 'border-b-2 border-green-600 text-gray-900' : 'text-gray-500 hover:text-gray-700'} hover:cursor-pointer"
			>
				Job detail
			</button>
			<button 
				onclick={() => activeTab = 'search'}
				class="pb-3 px-1 text-sm font-medium transition-colors {activeTab === 'search' ? 'border-b-2 border-green-600 text-gray-900' : 'text-gray-500 hover:text-gray-700'} hover:cursor-pointer"
			>
				Search result
			</button>
		</div>
	</div>
	<div class="p-6">
	{#if activeTab === 'detail'}
		<!-- Job Detail View using JobDetailCard -->
		<div class="rounded-lg bg-white border border-gray-200 h-full overflow-hidden preview-mode">
			<JobDetailCard 
				job={previewJob} 
				companyInfo={previewCompany}
				onApply={handlePreviewApply}
			/>
		</div>
	{:else}
		<!-- Search Result View -->
		<div class="space-y-4 max-w-sm mx-auto">
			<div class="p-4 border border-gray-200 rounded-lg animate-pulse">
				<div class="flex items-start gap-3">
					<div class="w-10 h-10 bg-gray-200 rounded"></div>
					<div class="flex-1">
						<div class="h-5 bg-gray-200 rounded w-3/4 mb-2"></div>
						<div class="h-4 bg-gray-200 rounded w-1/2 mb-3"></div>
						<div class="flex gap-2">
							<div class="h-6 bg-gray-200 rounded-full w-20"></div>
							<div class="h-6 bg-gray-200 rounded-full w-16"></div>
						</div>
						<div class="h-3 bg-gray-200 rounded w-full mt-3"></div>
						<div class="h-3 bg-gray-200 rounded w-5/6 mt-2"></div>
					</div>
				</div>
			</div>

			<!-- Active Job Card -->
			<div class="ring-2 ring-green-600 rounded-lg">
				<JobCard
					job={previewJob}
				/>
			</div>

			<!-- Skeleton Card 2 -->
			<div class="p-4 border border-gray-200 rounded-lg animate-pulse">
				<div class="flex items-start gap-3">
					<div class="w-10 h-10 bg-gray-200 rounded"></div>
					<div class="flex-1">
						<div class="h-5 bg-gray-200 rounded w-2/3 mb-2"></div>
						<div class="h-4 bg-gray-200 rounded w-1/3 mb-3"></div>
						<div class="flex gap-2">
							<div class="h-6 bg-gray-200 rounded-full w-24"></div>
						</div>
						<div class="h-3 bg-gray-200 rounded w-full mt-3"></div>
					</div>
				</div>
			</div>

			<!-- Skeleton Card 3 -->
			<div class="p-4 border border-gray-200 rounded-lg animate-pulse">
				<div class="flex items-start gap-3">
					<div class="w-10 h-10 bg-gray-200 rounded"></div>
					<div class="flex-1">
						<div class="h-5 bg-gray-200 rounded w-4/5 mb-2"></div>
						<div class="h-4 bg-gray-200 rounded w-2/5 mb-3"></div>
						<div class="flex gap-2">
							<div class="h-6 bg-gray-200 rounded-full w-20"></div>
							<div class="h-6 bg-gray-200 rounded-full w-20"></div>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
	</div>
</SideDrawer>

<style>
	:global(.preview-mode a),
	:global(.preview-mode button) {
		pointer-events: none;
		opacity: 0.6;
		cursor: default;
	}
	
	:global(.preview-mode a:hover),
	:global(.preview-mode button:hover) {
		background-color: initial;
	}
</style>