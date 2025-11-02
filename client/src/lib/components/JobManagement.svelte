<script lang="ts">
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth.svelte';
  import { JobService } from '$lib/services/jobService';
  import { Search } from 'lucide-svelte';
  import DeleteModal from '$lib/components/ui/DeleteModal.svelte';
  import Table from '$lib/components/ui/Table.svelte';
  import type { JobDisplay, JobFilters } from '$lib/types';

  interface Props {
    title?: string;
    description?: string;
    companyID?: string; // If provided, only show jobs for this company
    showCompanyColumn?: boolean; // Whether to show company column
    editUrlPrefix?: string; // URL prefix for edit action
  }

  let {
    title = "Job Management",
    description = "Manage and monitor all job postings across your platform",
    companyID,
    showCompanyColumn = true,
    editUrlPrefix = "/company/edit"
  }: Props = $props();

  // State Management
  let jobs = $state<JobDisplay[]>([]);
  let loading = $state(true);
  let searchQuery = $state('');
  let statusFilter = $state('all');
  let currentPage = $state(1);
  const itemsPerPage = 15;

  // Modal Management
  let showDeleteModal = $state(false);
  let jobToDelete = $state<JobDisplay | null>(null);
  let deleting = $state(false);

  const statusOptions = [
    { value: 'all', label: 'All Jobs' },
    { value: 'active', label: 'Active' },
    { value: 'closed', label: 'Closed' },
    { value: 'draft', label: 'Draft' }
  ];

  type ActionType = { label: string };
  
  async function handleAction(action: ActionType, job: JobDisplay) {
    if (action.label === 'Edit') {
      goto(`${editUrlPrefix}/${job.id}`);
    }
    if (action.label === 'View') {
      goto(`/app/jobs/${job.id}`);
    }
    if (action.label === 'Delete') {
      jobToDelete = job;
      showDeleteModal = true;
    }
  }

  async function confirmDelete(reason: string) {
    if (!jobToDelete) return;
    deleting = true;
    try {
      await JobService.deleteJobWithReason(jobToDelete.id, reason);
      jobs = jobs.filter(j => j.id !== jobToDelete?.id);
      showDeleteModal = false;
      jobToDelete = null;
    } catch (err) {
      console.error('Error deleting job:', err);
      throw err;
    } finally {
      deleting = false;
    }
  }

  // Filter jobs based on search and status
  const filteredJobs = $derived(
    jobs.filter(job => {
      const matchesSearch = !searchQuery || 
        job.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        job.company.toLowerCase().includes(searchQuery.toLowerCase());
      
      const matchesStatus = statusFilter === 'all' || job.status.toLowerCase() === statusFilter;
      
      return matchesSearch && matchesStatus;
    })
  );

  // Pagination
  const totalPages = $derived(Math.ceil(filteredJobs.length / itemsPerPage));

  // Reset page when filters change
  $effect(() => {
    searchQuery;
    statusFilter;
    currentPage = 1;
  });

  // Wait for auth to be ready before loading data
  $effect(() => {
    if (authStore.isAuthenticated && authStore.user) {
      loadAllJobs();
    }
  });

  async function loadAllJobs() {
    try {
      loading = true;
      const queryParams: JobFilters | undefined = companyID ? { companyID } : undefined;
      jobs = await JobService.loadJobsForDisplay(queryParams);
    } catch (err) {
      console.error('Error loading jobs:', err);
      jobs = [];
    } finally {
      loading = false;
    }
  }

  // Event handlers
  function handleSearch(event: Event) {
    searchQuery = (event.target as HTMLInputElement).value;
  }

  function handleStatusFilter(event: Event) {
    statusFilter = (event.target as HTMLSelectElement).value;
  }

  // Define columns based on whether to show company
  const columns = $derived([
    ...(showCompanyColumn ? [{ 
      key: 'job', 
      label: 'Job', 
      width: '400px', 
      align: 'left' as const, 
      cellClass: 'px-6 py-2.5 whitespace-nowrap',
      type: 'job' as const,
      jobFields: { logo: 'companyLogo', company: 'company', title: 'title' }
    }] : [{ 
      key: 'title', 
      label: 'Job Title', 
      width: '400px', 
      align: 'left' as const,
      type: 'text' as const
    }]),
    { key: 'location', label: 'Location', width: '200px', align: 'left' as const, type: 'text' as const },
    { key: 'posted', label: 'Posted', width: '150px', align: 'left' as const, type: 'text' as const },
    { key: 'expires', label: 'Deadline', width: '150px', align: 'left' as const, type: 'text' as const },
    { key: 'applicants', label: 'Applied', width: '100px', align: 'left' as const, type: 'text' as const },
    { 
      key: 'status', 
      label: 'Status', 
      width: '120px', 
      align: 'left' as const,
      type: 'badge' as const,
      badgeType: 'status' as const
    },
    { 
      key: 'actions', 
      label: '', 
      width: '80px', 
      align: 'center' as const, 
      cellClass: 'px-6 py-4 whitespace-nowrap text-center relative',
      type: 'actions' as const,
      actions: {
        items: [
          { label: 'View' },
          { label: 'Edit' },
          { label: 'Delete', variant: 'danger' }
        ]
      }
    }
  ]);
</script>

<div>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-semibold text-gray-900 mb-1">{title}</h1>
      <p class="text-base text-gray-600 mb-6">{description}</p>
    </div>

    <!-- Filters Card -->
    <div class="mb-6">
      <div class="flex flex-col lg:flex-row gap-4">
        <!-- Search -->
        <div class="flex-1">
          <div class="relative">
            <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search jobs{showCompanyColumn ? ' or companies' : ''}..."
              class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg max-w-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              bind:value={searchQuery}
              oninput={handleSearch}
            />
          </div>
        </div>

        <!-- Status Filter -->
        <select
          class="px-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent bg-white min-w-[150px]"
          bind:value={statusFilter}
          onchange={handleStatusFilter}
        >
          {#each statusOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>
    </div>

    <!-- Jobs Table -->
    <Table
      {columns}
      data={filteredJobs}
      {loading}
      {currentPage}
      {totalPages}
      {itemsPerPage}
      onPageChange={(page) => currentPage = page}
      emptyMessage={searchQuery || statusFilter !== 'all' ? "No jobs found" : "No job postings available"}
      emptyDescription={searchQuery || statusFilter !== 'all' ? "Try adjusting your filters" : "No job postings available"}
      onAction={(action, job) => handleAction({ label: action }, job)}
    />
  </div>
</div>

<!-- Delete Confirmation Modal -->
<DeleteModal
  bind:isOpen={showDeleteModal}
  onClose={() => {
    showDeleteModal = false;
    jobToDelete = null;
  }}
  onConfirm={confirmDelete}
  title="Delete Job Posting"
  itemName={jobToDelete?.title || ''}
  description="You're about to delete this job posting. This action cannot be undone and will remove the job posting along with all associated applications."
  reasonPlaceholder="Please provide a detailed reason for deleting this job posting..."
  confirmButtonText="Delete Job"
  isDeleting={deleting}
/>