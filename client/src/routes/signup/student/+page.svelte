<script lang="ts">
	import { get } from 'svelte/store';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-french-toast';
	import { getUserInfo } from '$lib/utils/auth';
	import { fileService, type FileMetadata } from '$lib/services/fileService';
	import { ArrowLeft } from 'lucide-svelte';
	import { fly } from 'svelte/transition';
	import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
	import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
	import FormInput from '$lib/components/auth/FormInput.svelte';
	import FormButton from '$lib/components/auth/FormButton.svelte';
	import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
	import FileItem from '$lib/components/files/FileItem.svelte';
	import FilePreviewModal from '$lib/components/files/FilePreviewModal.svelte';
	import FileUploadModal from '$lib/components/files/FileUploadModal.svelte';
	import DeleteConfirmModal from '$lib/components/files/DeleteConfirmModal.svelte';
	import { userService } from '$lib/services/userService';
	import PDPAModal from '$lib/components/modals/PDPAModal.svelte';
	import seekerDashboard from '$lib/assets/seeker-dashboard.png';

	let showPDPA = false;

	let currentStep = 1;
	let email = '';
	let fullName = '';
	let desiredRole = '';
	let linkedIn = '';
	let github = '';
	let portfolio = '';
	let phone = '';
	let location = '';
	let aboutMe = '';
	let dateOfBirth = '';
	let username = '';
	let avatar = '';
	let uid = '';

	let files: FileMetadata[] = [];
	let isUploadModalOpen = false;
	let selectedFile: FileMetadata | null = null;
	let isPreviewModalOpen = false;
	let fileToDelete: FileMetadata | null = null;
	let isDeleteModalOpen = false;
	let isDeleting = false;

	let userInfo = getUserInfo();
	let userID = userInfo?.userID || '';
	const userRole = userInfo?.role || 'jobSeeker';

	function goBackToEmailStep() {
		currentStep = 1;
	}

	function handleSignup() {
		console.log('Student signup:', {
			email,
			firstName: fullName,
			lastName: desiredRole,
			portfolio: portfolio,
			github: github
		});
	}

	async function loadJobseekerDetails() {
		try {
			const user = await userService.getCurrentUser();
			const frontendData = userService.transformToFrontendFormat(user);
			email = user.email || '';
			uid = frontendData.userID as string;
			username = frontendData.name as string;
			avatar = frontendData.avatarURL as string;
			fullName = frontendData.fullName as string;
			location = frontendData.location as string;
			phone = frontendData.phone as string;
			linkedIn = frontendData.linkedIn as string;
			desiredRole = frontendData.desiredRole as string;
			aboutMe = frontendData.aboutMe as string;
			dateOfBirth = frontendData.dateOfBirth as string;
			portfolio = frontendData.portfolio as string;
			github = frontendData.github as string;
			console.log(frontendData);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load jobseeker details');
		}
	}

	async function handleSubmit() {
		try {
			const payload = await userService.transformToBackendFormat(
				{
					fullName,
					location,
					phone,
					linkedIn,
					desiredRole,
					aboutMe,
					dateOfBirth,
					portfolio,
					github
				},
				'jobSeeker'
			);
			payload.role = 'jobSeeker';
			payload.email = email;
			payload.provider = 'google';
			payload.name = username;
			payload.avatarURL = avatar;
			payload.userID = uid;
			await userService.updateUser(userID!, payload);
			toast.success('Jobseeker details updated successfully');
			goto('/app/jobs');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to update jobseeker details');
		}
	}
	async function loadFiles() {
		if (!userID) return;
		try {
			files = await fileService.listUserFiles(userID);
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to load files');
		}
	}

	function handleUploadClick() {
		isUploadModalOpen = true;
	}

	async function handleUploadSuccess(file: FileMetadata) {
		await loadFiles();
		toast.success(`${fileService.getCategoryLabel(file.category)} uploaded successfully`);
	}

	function handlePreview(file: FileMetadata) {
		selectedFile = file;
		isPreviewModalOpen = true;
	}

	function handleDeleteClick(file: FileMetadata) {
		fileToDelete = file;
		isDeleteModalOpen = true;
	}

	async function handleDeleteConfirm() {
		if (!fileToDelete) return;

		isDeleting = true;
		try {
			await fileService.deleteFile(fileToDelete.id);
			files = files.filter((f) => f.id !== fileToDelete.id);
			toast.success('Document deleted successfully');
			isDeleteModalOpen = false;
			fileToDelete = null;
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to delete document');
		} finally {
			isDeleting = false;
		}
	}

	function handleDeleteCancel() {
		fileToDelete = null;
	}

	onMount(async () => {
		const urlStep = Number(get(page).url.searchParams.get('currentStep'));
		if (urlStep === 2) {
			currentStep = 2;
		}
		if (currentStep === 2) {
			userInfo = getUserInfo();
			userID = userInfo?.userID;
			let retries = 0;
			while (!userInfo && retries < 10) {
				await new Promise((resolve) => setTimeout(resolve, 50));
				userInfo = getUserInfo();
				userID = userInfo?.userID;
				retries++;
			}
			if (!userID) {
				console.error('Failed to load userID after login!');
				toast.error('Unable to load user session. Try reloading the page.');
				return;
			}
		}
		if (currentStep === 2) {
			await loadFiles();
			await loadJobseekerDetails();
		}
	});
</script>

<AuthLayout backHref="/signup">
  <div slot="left-side" class="h-full flex flex-col justify-between p-12 overflow-hidden">
    <div class="max-w-2xl pt-20 px-20">
      <p class="text-gray-700 text-2xl leading-relaxed mb-6">
        Launch your career with Job Applier 3000! Connect with top companies looking for talented CPSK students and alumni like you
      </p>
      <p class="text-2xl font-semibold text-gray-900">Your dream job is just a click away!</p>
    </div>

    <div class="flex justify-start">
      <div class="w-full max-w-4xl -ml-18 -mb-18">
        <div class="bg-white rounded-3xl border-5 border-black overflow-hidden">
          <!-- Jobs listing image -->
          <img src={seekerDashboard} alt="Seeker dashboard preview" class="w-full h-160 object-cover" />
        </div>
      </div>
    </div>
  </div>

  <AuthHeader />
  
  {#if currentStep === 1}
    <div class="mb-8" in:fly={{ x: -20, duration: 200 }}>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Sign up</h1>
      <p class="text-sm text-gray-500">Register as a job seekeer to start finding your dream job</p>
    </div>
    
    <GoogleOAuthButton text="Continue with Google" userType="jobSeeker" />

		<p class="mt-8 text-center text-sm text-gray-500">
			By using our service, you consent to the processing of your Personal Data as described in our

			<button
				type="button"
				on:click={() => (showPDPA = true)}
				class="font-medium text-blue-600 hover:text-blue-700"
			>
				Privacy Notice
			</button>
		</p>

		<p class="mt-8 text-center text-sm text-gray-500">
			Already have an account?
			<a href="/login" class="font-medium text-green-600 hover:text-green-700">Log in</a>
		</p>
	{:else}
		<div class="mb-6" in:fly={{ x: 20, duration: 200 }}>
			<button
				type="button"
				on:click={goBackToEmailStep}
				class="mb-4 flex cursor-pointer items-center text-sm text-gray-500 transition-colors hover:text-gray-700"
			>
				<ArrowLeft class="mr-1 h-4 w-4" />
				<span>{email}</span>
			</button>
			<h1 class="mb-2 text-3xl font-semibold text-gray-900">Complete Profile</h1>
			<p class="text-sm text-gray-500">Tell us more about yourself</p>
		</div>

		<form
			on:submit|preventDefault={handleSignup}
			class="space-y-4"
		>
			<div class="grid grid-cols-2 gap-4">
				<FormInput
					id="fullname"
					label="Full Name"
					placeholder="Enter your name..."
					bind:value={fullName}
					required
				/>
				<FormInput
					id="desiredrole"
					label="Desired Role"
					placeholder="Enter your preferred role..."
					bind:value={desiredRole}
					required
				/>
			</div>

			<div class="grid grid-cols-2 gap-4">
				<FormInput
					id="portfolio"
					type="url"
					label="Portfolio"
					placeholder="Your portfolio website (Optional)"
					bind:value={portfolio}
				/>
				<FormInput
					id="github"
					type="url"
					label="Github"
					bind:value={github}
					placeholder="Your github link (Optional)"
				></FormInput>
			</div>

			<div class="mt-4">
				<button
					type="button"
					on:click={handleUploadClick}
					class="rounded bg-green-600 px-4 py-2 text-white"
				>
					Upload Files
				</button>

				{#if files.length > 0}
					<div class="mt-2 space-y-2">
						{#each files as file (file.id)}
							<FileItem
								{file}
								onPreview={() => handlePreview(file)}
								onDelete={() => handleDeleteClick(file)}
							/>
						{/each}
					</div>
				{/if}

				<FileUploadModal
					bind:isOpen={isUploadModalOpen}
					{userRole}
					onUploadSuccess={handleUploadSuccess}
				/>

				{#if selectedFile}
					<FilePreviewModal
						bind:isOpen={isPreviewModalOpen}
						fileId={selectedFile.id}
						filename={selectedFile.filename}
					/>
				{/if}

				{#if fileToDelete}
					<DeleteConfirmModal
						bind:isOpen={isDeleteModalOpen}
						filename={fileToDelete.filename}
						{isDeleting}
						onConfirm={handleDeleteConfirm}
						onCancel={handleDeleteCancel}
					/>
				{/if}
			</div>

			<FormButton type="submit" on:click={handleSubmit}>Create Account</FormButton>
		</form>
	{/if}
</AuthLayout>
<PDPAModal bind:isVisible={showPDPA} />
