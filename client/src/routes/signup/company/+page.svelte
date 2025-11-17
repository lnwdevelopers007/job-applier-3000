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
  import FormInput from '$lib/components/auth/FormInput.svelte';
  import FormSelect from '$lib/components/auth/FormSelect.svelte';
  import FormButton from '$lib/components/auth/FormButton.svelte';
  import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
  import FileItem from '$lib/components/files/FileItem.svelte';
  import FilePreviewModal from '$lib/components/files/FilePreviewModal.svelte';
  import FileUploadModal from '$lib/components/files/FileUploadModal.svelte';
	import DeleteConfirmModal from '$lib/components/files/DeleteConfirmModal.svelte';
	import { userService } from '$lib/services/userService';
  import companyApplicantsImage from '$lib/assets/company-applicants.png';
  import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
  import AuthHeader from '$lib/components/auth/AuthHeader.svelte';

  let currentStep = 1;
  let email = '';
  let companyName = '';
  let aboutCompany = '';
  let industry = '';
  let size = '';
  let website = '';
  let companyLogo = '';
  let foundedYear = '';
  let headquarters = '';
  let companyLinkedin = '';
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
  const userRole = userInfo?.role || 'company';

  function goBackToEmailStep() {
    currentStep = 1;
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
      files = files.filter(f => f.id !== fileToDelete.id);
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

  async function loadCompanyDetails() {
    try {
      const user = await userService.getCurrentUser();
      const frontendData = userService.transformToFrontendFormat(user);
      email = user.email || '';
      uid = frontendData.userID as string;
      username = frontendData.name as string;
      avatar = frontendData.avatarURL as string;
      companyName = frontendData.companyName as string;
      aboutCompany = frontendData.aboutCompany as string;
      industry = frontendData.industry as string;
      size = frontendData.companySize as string;
      website = frontendData.companyWebsite as string;
      companyLogo = frontendData.companyLogo as string;
      foundedYear = frontendData.foundedYear as string;
      headquarters = frontendData.headquarters as string;
      companyLinkedin = frontendData.companyLinkedin as string;
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to load company details');
    }
  }

  async function handleSubmit() {
    try {
      const payload = await userService.transformToBackendFormat({
        companyName,
        aboutUs: aboutCompany,
        industry,
        companySize: size,
        companyWebsite: website,
        logo: companyLogo,
        foundedYear,
        headquarters,
        linkedIn: companyLinkedin,
      }, 'company');
      payload.role = 'company';
      payload.email = email;
      payload.provider = 'google'
      payload.name = username;
      payload.avatarURL = avatar;
      payload.userID = uid;
      await userService.updateUser(userID!, payload);
      toast.success('Company details updated successfully');
      goto('/company/dashboard');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update company details');
    }
  }

  onMount(async () => {
    const stepFromUrl = Number(get(page).url.searchParams.get('currentStep'));
    if (stepFromUrl === 2) currentStep = 2;

    if (currentStep===2) {
      userInfo = getUserInfo();
      userID = userInfo?.userID;
      let retries = 0;
      while (!userInfo && retries < 10) {
        await new Promise(resolve => setTimeout(resolve, 50));
        userInfo = getUserInfo();
        userID = userInfo?.userID;
        retries++;
      }
    if (!userID) {
      console.error("Failed to load userID after login!");
      toast.error("Unable to load user session. Try reloading the page.");
      return;
    }
      // if (user?.role === 'company') goto('/company/dashboard');
      // else goto('/app/jobs');
    }

    // Load files when Step 2 is active
    if (currentStep === 2) {
      await loadFiles();
      await loadCompanyDetails();
    }
  });
</script>

<AuthLayout backHref="/signup">
  <div slot="left-side" class="h-full flex flex-col justify-between p-12 overflow-hidden">
    <div class="max-w-2xl pt-20 px-20">
      <p class="text-gray-700 text-2xl leading-relaxed mb-6">
        Join Job Applier 3000 as a recruiter and connect with talented CPSK students and alumni ready to make their mark in the tech industry
      </p>
      <p class="text-2xl font-semibold text-gray-900">Start hiring top talent today!</p>
    </div>

    <div class="flex justify-start">
      <div class="w-full max-w-4xl -ml-18 -mb-18">
        <div class="bg-white rounded-3xl border-5 border-black overflow-hidden">
          <!-- Jobs listing image -->
          <img src={companyApplicantsImage} alt="Jobs listing dashboard preview" class="w-full h-160 object-cover" />
        </div>
      </div>
    </div>
  </div>

  <AuthHeader />

  {#if currentStep === 1}
    <div class="mb-8" in:fly={{ x: -20, duration: 200 }}>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Sign up</h1>
      <p class="text-sm text-gray-500">Register as a recruiter to start hiring talents</p>
    </div>

    <GoogleOAuthButton text="Continue with Google" userType="company" />

    <p class="text-center text-sm text-gray-500 mt-8">
      Already have an account?
      <a href="/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
    </p>

  {:else}
    <div class="mb-6" in:fly={{ x: 20, duration: 200 }}>
      <button 
        type="button" 
        on:click={goBackToEmailStep}
        class="flex items-center text-sm text-gray-500 hover:text-gray-700 mb-4 transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 mr-1" />
        <span>{email}</span>
      </button>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Company Details</h1>
      <p class="text-sm text-gray-500">Tell us more about your company</p>
    </div>

    <form on:submit|preventDefault={() => console.log('Signup submitted')} class="space-y-4">
      <FormInput
        id="companyName"
        label="Company Name"
        placeholder="Enter company name..."
        bind:value={companyName}
        required
      />

      <div class="grid grid-cols-2 gap-4">
        <FormSelect id="industry" label="Industry" bind:value={industry} placeholder="Select industry" required>
          <option value="Technology">Technology</option>
          <option value="Finance">Finance</option>
          <option value="Healthcare">Healthcare</option>
          <option value="Education">Education</option>
          <option value="Retail">Retail</option>
          <option value="Manufacturing">Manufacturing</option>
          <option value="Other">Other</option>
        </FormSelect>

        <FormSelect id="companySize" label="Company Size" bind:value={size} placeholder="Select size" required>
          <option value="1-10 employees">1-10 employees</option>
          <option value="11-50 employees">11-50 employees</option>
          <option value="51-200 employees">51-200 employees</option>
          <option value="201-500 employees">201-500 employees</option>
          <option value="501-1000 employees">501-1000 employees</option>
          <option value="1000+ employees">1000+ employees</option>
        </FormSelect>
      </div>

      <FormInput
        id="website"
        type="url"
        label="Website (Optional)"
        placeholder="Enter website URL..."
        bind:value={website}
      />

      <div class="mt-4">
        <button type="button" on:click={handleUploadClick} class="px-4 py-2 bg-green-600 text-white rounded">
          Upload Files
        </button>

        {#if files.length > 0}
          <div class="mt-2 space-y-2">
            {#each files as file (file.id)}
              <FileItem {file} 
                onPreview={() => handlePreview(file)}
                onDelete={() => handleDeleteClick(file)}
              />
            {/each}
          </div>
        {/if}

        <FileUploadModal 
          bind:isOpen={isUploadModalOpen} {userRole} 
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

      <FormButton type="submit" onclick={handleSubmit}>Create Account</FormButton>
    </form>
  {/if}
</AuthLayout>
