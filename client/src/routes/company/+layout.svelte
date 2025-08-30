<div class="min-h-screen bg-gray-100 flex flex-col">

  <header class="flex items-center justify-between bg-white shadow px-6 py-4">
    <div class="flex items-start space-x-1">
        <h2 class="text-xl font-bold text-black">Job Applier </h2>
        <h2 class="text-xl font-bold text-green-600">3000</h2>
    </div>

    <div class="flex items-center space-x-3">
      <span class="text-right">
        <span class="block font-medium text-gray-700">Company name</span>
        <span class="block text-sm text-gray-500">Recruiter</span>
      </span>
      <img src={user} alt="Avatar" class="w-10 h-10 rounded-full" />
    </div>
  </header>

  <div class="flex flex-1 p-6 gap-6 items-start">
    <aside class="w-72 bg-white shadow rounded-lg p-4">
      <nav class="space-y-2">
        <a href="/" class="nav-link">
          <House class="w-5 h-5" />
          <span>Home</span>
        </a>
        <a href="/company/dashboard" class="nav-link">
          <ChartColumn class="w-5 h-5" />
          <span>Dashboard</span>
        </a>
        <a href="/company/applicants" class="nav-link">
          <Users class="w-5 h-5" />
          <span>Applicants</span>
        </a>
        <a href="/company/analytics" class="nav-link">
          <ChartLine class="w-5 h-5" />
          <span>Analytics</span>
        </a>
        <div class="border-t border-gray-200 my-2"></div>
        <a href="/company/create" class="nav-link">
          <CirclePlus class="w-5 h-5" />
          <span>Post New Job</span>
        </a>
        <a href="/company/profile" class="nav-link">
          <Building class="w-5 h-5" />
          <span>Company Profile</span>
        </a>
        <a href="/company/setting" class="nav-link">
          <Cog class="w-5 h-5" />
          <span>Settings</span>
        </a>
      </nav>
    </aside>

    <main class="flex-1 p-1 min-w-0">
      <slot />
    </main>
  </div>
</div>

<script>
  import user from '$lib/assets/user.png';
  import database from '$lib/assets/database.png';
  import { House, Cog, Building, ChartColumn, ChartLine, Users, CirclePlus } from 'lucide-svelte';
</script>
