<script lang="ts">
  import { onMount } from 'svelte';
  import { jwtDecode } from 'jwt-decode';
  import apiFetch from './apifetch';

  type TokenPayload = {
    email: string;
    name: string;
    exp: number;
  };

  // Using `any` since we don't care about the specific structure for raw display
  let user: TokenPayload | null = null;
  let jobs: any[] = [];

  async function loadJobs(): Promise<any[]> {
    try {
      const res = await apiFetch("http://localhost:8080/jobs/");
      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      return await res.json();
    } catch (err) {
      console.error("Failed to fetch jobs:", err);
      return [{ error: "Failed to load jobs." }]; // Show an error object
    }
  }

  onMount(async () => {
    // Get token from query params
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');

    if (token) {
      try {
        user = jwtDecode<TokenPayload>(token);
        localStorage.setItem('access_token', token);
      } catch (err) {
        console.error('Invalid token', err);
      }
    }

    jobs = await loadJobs();
  });
</script>

{#if user}
  <h1>Welcome, {user.name}!</h1>
  <p>Email: {user.email}</p>

  <hr />
  
  <pre>{JSON.stringify(jobs, null, 2)}</pre>

{:else}
  <p>Loading user...</p>
{/if}