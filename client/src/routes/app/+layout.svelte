<div class="min-h-screen bg-gray-100 flex flex-col">

  <header class="flex items-center justify-between bg-white shadow px-6 py-4">
    <div class="flex items-center space-x-3">
      <img src={database} alt="Logo" class="w-10 h-10 rounded-full" />
      <div>
        <h2 class="text-lg font-semibold">Job Applier 3000</h2>
        <p class="text-sm text-gray-500">KU Computer Engineering</p>
      </div>
    </div>

    <div class="flex items-center space-x-3">
      <span class="text-right">
        <span class="block font-medium text-gray-700">Kanye West</span>
        <span class="block text-sm text-gray-500">Student</span>
      </span>
      <img src={user} alt="Avatar" class="w-10 h-10 rounded-full" />
    </div>
  </header>

  <div class="flex flex-1 p-6 gap-6 items-start">
    <aside class="w-72 bg-white shadow rounded-lg p-4">
      <nav class="space-y-2">
        <a href="/app" class="nav-link">
          <House class="w-5 h-5" />
          <span>Home</span>
        </a>
        <a href="/app/jobs" class="nav-link">
          <CircleCheck class="w-5 h-5" />
          <span>Jobs</span>
        </a>
        <a href="/app/applications" class="nav-link">
          <ChartColumn class="w-5 h-5" />
          <span>Applications</span>
        </a>
        <a href="/app/settings" class="nav-link">
          <Cog class="w-5 h-5" />
          <span>Settings</span>
        </a>
        <a href="/app/requests" class="nav-link">
          <Mail class="w-5 h-5" />
          <span>Requests</span>
        </a>
      </nav>
    </aside>

    <main class="flex-1 bg-white shadow rounded-lg p-6">
      <slot />
    </main>
  </div>
</div>

<script>
  import user from '$lib/assets/user.png';
  import database from '$lib/assets/database.png';
  import { House, Cog, Mail, ChartColumn, CircleCheck } from 'lucide-svelte';
</script>
