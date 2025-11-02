<script lang="ts">
  import JobManagement from '$lib/components/JobManagement.svelte';
</script>

<JobManagement
  title="Job Management"
  description="Manage and monitor all job postings across your platform"
  showCompanyColumn={true}
  editUrlPrefix="/company/edit"
/>