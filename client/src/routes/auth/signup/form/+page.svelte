<script>
  import { page } from '$app/stores';
  $: role = $page.url.searchParams.get('role');
  $: isStudent = role === 'student';
  $: isCompany = role === 'company';
</script>

<div>
  {#if isStudent}
    <h1>Student Signup Form</h1>
    <!-- Student signup form -->
  {:else if isCompany}
    <h1>Company Signup Form</h1>
    <!-- Company signup form -->
  {:else}
    <h1>Please select a role</h1>
    <a href="/auth/signup">Go back</a>
  {/if}
</div>