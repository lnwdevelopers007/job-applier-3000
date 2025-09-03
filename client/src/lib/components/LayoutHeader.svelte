<script>
  import user from '$lib/assets/user.png';
  
  let { companyName = "Company name", userRole = "Recruiter" } = $props();
</script>

<header class="sticky top-0 z-50 flex items-center justify-between bg-white border-b border-gray-200 px-6 py-3">
  <div class="flex items-start space-x-1">
      <h2 class="text-lg font-semibold text-black">Job Applier </h2>
      <h2 class="text-lg font-semibold text-green-600">3000</h2>
  </div>

  <div class="flex items-center space-x-3">
    <img src={user} alt="Avatar" class="w-8 h-8 rounded-full" />
    <span class="text-left">
      <span class="text-sm block font-semibold text-gray-700">{companyName}</span>
      <span class="text-xs block text-gray-500">{userRole}</span>
    </span>
  </div>
</header>