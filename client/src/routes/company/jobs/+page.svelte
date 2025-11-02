<script lang="ts">
  import { authStore } from '$lib/stores/auth.svelte';
  import JobManagement from '$lib/components/JobManagement.svelte';
  
  // Get the current user's company ID from auth store
  const companyID = $derived(authStore.user?.id);
</script>

<JobManagement
  title="Job Postings"
  description="Manage and monitor your company's job postings"
  companyID={companyID}
  showCompanyColumn={true}
  editUrlPrefix="/company/edit"
/>