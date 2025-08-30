<div class="p-2 bg-gray-100">
  <h1 class="text-2xl font-bold text-gray-900">
    All Applicants
  </h1>
  <p class="my-3 text-base text-gray-600">
    Manage candidates across all job postings
  </p>
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search candidates by name, email, skills or position..."
          class="flex-1 ml-2 outline-none border-none text-sm"
        />
      </div>

      <div class="flex gap-2 my-2">      
      <select class="px-7 py-1 bg-white border border-gray-200 rounded text-sm">
        <option>All Jobs</option>
        <option>AI Researcher</option>
        <option>Developer</option>
        <option>Tester</option>
      </select>

      <select class="px-7 py-1 bg-white border border-gray-200 rounded text-sm">
        <option>All Statuses</option>
        <option>Accepted</option>
        <option>Pending</option>
        <option>Rejected</option>
      </select>
      <div class="border-l border-gray-300 h-6 py-1"></div>
        <button class="px-3 py-1 bg-white border-1 border-gray-200 rounded-full text-sm">All</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">New (34)</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">High Match</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Recent Grads</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Experienced</button>
      </div>
</div>

<script>
  import { Search } from 'lucide-svelte';
</script>