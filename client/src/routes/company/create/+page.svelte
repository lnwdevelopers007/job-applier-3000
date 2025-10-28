<script>
	import { goto } from '$app/navigation';
	import { fly } from 'svelte/transition';
	import { ArrowLeft, ArrowRight, ChevronLeft } from 'lucide-svelte';
	import Stepper from '$lib/components/job-post-creation/Stepper.svelte';
	import BasicInfoForm from '$lib/components/job-post-creation/forms/BasicInfoForm.svelte';
	import DescriptionForm from '$lib/components/job-post-creation/forms/DescriptionForm.svelte';
	import RequirementForm from '$lib/components/job-post-creation/forms/RequirementForm.svelte';
	import PostSettingForm from '$lib/components/job-post-creation/forms/PostSettingForm.svelte';
	import JobPreviewDrawer from '$lib/components/job-post-creation/JobPreviewDrawer.svelte';
	import { getUserInfo } from '$lib/utils/auth';

	let currentStep = $state(1); // Current step (1-4)
	let isPreviewOpen = $state(false);
	let stepper;
	const userInfo = getUserInfo();
	const steps = [
		{ id: 1, name: 'Basic Info' },
		{ id: 2, name: 'Description' },
		{ id: 3, name: 'Requirement' },
		{ id: 4, name: 'Post Settings' }
	];

	let formData = $state({
		// Basic Info
		jobTitle: '',
		companyID: userInfo.userID,
		location: '',
		workType: 'full-time',
		workArrangement: 'on-site',
		currency: 'THB',
		minSalary: 1,
		maxSalary: 1,

		// Description
		jobDescription: '',
		jobSummary: '',

		// Requirements
		requiredSkills: '',
		experienceLevel: '',
		education: '',
		niceToHave: '',

		// Post Settings
		applicationDeadline: '',
		numberOfPositions: 1,
		visibility: 'public',
		emailNotifications: false,
		autoReject: false
	});

	function handleBackNavigation() {
		goto('/company/dashboard');
	}

	function handleSaveDraft() {
		// Handle save draft
	}

	function handlePreview() {
		isPreviewOpen = true;
	}

	function buildPayload(formData) {
		function withDeadlineTime(dateStr) {
			const d = new Date(dateStr);
			d.setUTCHours(16, 59, 0, 0);
			return d.toISOString();
		}
		return {
			// Basic Info
			title: formData.jobTitle || 'Test Job Title',
			companyID: String(formData.companyID || '64f0c44a27b1c27f4d92e9a2'),
			location: formData.location || 'Bangkok, Thailand',
			workType: formData.workType,
			workArrangement: formData.workArrangement,
			currency: formData.currency,
			minSalary: Number(formData.minSalary || 0),
			maxSalary: Number(formData.maxSalary || 0),

			// Description
			jobDescription: formData.jobDescription || 'Test description',
			jobSummary: formData.jobSummary || 'Test summary',

			// Requirements
			requiredSkills:
				Array.isArray(formData.requiredSkills) && formData.requiredSkills.length
					? formData.requiredSkills.join(', ')
					: 'JS, Node',
			experienceLevel: formData.yearsOfExperience || 'Mid-Level',
			education: formData.educationLevel || 'Bachelor',
			niceToHave: formData.niceToHave || '',
			questions: formData.screeningQuestions || 'What is your expected salary?',

			// Post Settings
			postOpenDate: formData.postingOpenDate
				? new Date(formData.postingOpenDate).toISOString()
				: new Date().toISOString(),
			applicationDeadline: formData.postingCloseDate
				? withDeadlineTime(formData.postingCloseDate)
				: withDeadlineTime(new Date()),
			numberOfPositions: Number(formData.numberOfPositions || 1),
			visibility: formData.visibility || 'public',
			emailNotifications: Boolean(formData.emailNotifications),
			autoReject: Boolean(formData.autoReject)
		};
	}

	const handleProgress = async (stepIncrement) => {
		if (currentStep === 4 && stepIncrement === 1) {
			try {
				const payload = buildPayload(formData);
				console.log('Sending payload:', payload);

				const res = await fetch('/jobs/', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify(payload)
				});

				if (!res.ok) {
					const text = await res.text();
					console.error('Server error body:', text);
					throw new Error(`Failed: ${res.status} ${text}`);
				}

				const data = await res.json();
				console.log('Job created:', data);
				goto('/company/dashboard');
			} catch (err) {
				console.error('Error creating job:', err);
			}
		} else {
			stepper.handleProgress(stepIncrement);
		}
	};
</script>

<div class="mx-auto max-w-4xl py-8">
	<div class="mb-2 px-8">
		<div class="flex items-center gap-4">
			<button onclick={handleBackNavigation} class="rounded p-1 hover:bg-gray-200">
				<ChevronLeft class="h-5 w-5 text-gray-600" />
			</button>
			<h1 class="text-xl font-medium text-gray-900">Create Job Post</h1>
		</div>
	</div>

	<Stepper {steps} bind:currentStep bind:this={stepper} />

	<form class="px-8 pb-8">
		{#if currentStep === 1}
			<div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
				<BasicInfoForm bind:formData />
			</div>
		{:else if currentStep === 2}
			<div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
				<DescriptionForm bind:formData />
			</div>
		{:else if currentStep === 3}
			<div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
				<RequirementForm bind:formData />
			</div>
		{:else if currentStep === 4}
			<div in:fly={{ x: 50, duration: 300, delay: 200 }} out:fly={{ x: -50, duration: 200 }}>
				<PostSettingForm bind:formData />
			</div>
		{/if}

		<div class="mt-8 flex justify-between pt-6">
			<div>
				{#if currentStep > 1}
					<button
						type="button"
						onclick={() => handleProgress(-1)}
						class="flex items-center px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800"
					>
						<ArrowLeft class="mr-1 h-3 w-3" />
						Back
					</button>
				{/if}
			</div>
			<div class="flex gap-3">
				<button
					type="button"
					onclick={handleSaveDraft}
					class="px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800"
				>
					Save Draft
				</button>
				<button
					type="button"
					onclick={handlePreview}
					class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
				>
					Preview
				</button>
				<button
					type="button"
					onclick={() => handleProgress(1)}
					class="flex items-center rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700"
				>
					{currentStep < 4 ? 'Continue' : 'Publish'}
					<ArrowRight class="ml-1 h-3 w-3" />
				</button>
			</div>
		</div>
	</form>
</div>

<JobPreviewDrawer bind:isOpen={isPreviewOpen} {formData} />
