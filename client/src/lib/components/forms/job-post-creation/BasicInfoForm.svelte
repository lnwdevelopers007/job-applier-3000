<script>
  let { formData = $bindable() } = $props();
</script>

<div class="space-y-6">
  <h2 class="text-lg font-medium text-gray-900 mb-6">Basic Job Information</h2>
  
  <!-- Job Title -->
  <div>
    <label for="jobTitle" class="block text-sm font-medium text-gray-700 mb-2">
      Job title <span class="text-red-500">*</span>
    </label>
    <input
      id="jobTitle"
      type="text"
      bind:value={formData.jobTitle}
      placeholder="Enter job title..."
      class="text-sm min-w-lg px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
    />
  </div>

  <!-- Location -->
  <div>
    <label for="location" class="block text-sm font-medium text-gray-700 mb-2">
      Location <span class="text-red-500">*</span>
    </label>
    <input
      id="location"
      type="text"
      bind:value={formData.location}
      placeholder="Enter location..."
      class="text-sm min-w-lg px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
    />
  </div>

  <!-- Work Type -->
  <fieldset>
    <legend class="block text-sm font-medium text-gray-700 mb-3">
      Work type <span class="text-red-500">*</span>
    </legend>
    <div class="space-y-3">
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workType}
          value="full-time"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Full-time</span>
      </label>
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workType}
          value="part-time"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Part-time</span>
      </label>
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workType}
          value="contract"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Contract</span>
      </label>
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workType}
          value="casual"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Casual</span>
      </label>
    </div>
  </fieldset>

  <!-- Work Arrangement -->
  <fieldset>
    <legend class="block text-sm font-medium text-gray-700 mb-3">
      Work arrangement <span class="text-red-500">*</span>
    </legend>
    <div class="space-y-3">
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workArrangement}
          value="on-site"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">On-site</span>
      </label>
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workArrangement}
          value="hybrid"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Hybrid</span>
      </label>
      <label class="flex items-center">
        <input
          type="radio"
          bind:group={formData.workArrangement}
          value="remote"
          class="w-4 h-4 text-green-500 border-gray-300 focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Remote</span>
      </label>
    </div>
  </fieldset>

  <!-- Pay Range -->
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">Pay range</label>
    <p class="text-sm text-gray-500 mb-3">Enter a range to offer candidates.</p>
    
    <div class="flex items-center gap-2">
      <!-- Currency -->
      <div class="w-20">
        <label for="currency" class="block text-xs text-gray-500 mb-1">Currency</label>
        <div class="relative">
          <select
            id="currency"
            bind:value={formData.currency}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 bg-white appearance-none pr-6 text-sm"
          >
            <option value="THB">THB</option>
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
          </select>
        </div>
      </div>
      
      <!-- From -->
      <div class="flex-shrink-0">
        <label for="minSalary" class="block text-xs text-gray-500 mb-1">From</label>
        <input
          id="minSalary"
          type="number"
          bind:value={formData.minSalary}
          placeholder="Enter minimum salary"
          class="text-sm min-w-xs px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
        />
      </div>
      
      <!-- To -->
      <div class="flex-shrink-0">
        <label for="maxSalary" class="block text-xs text-gray-500 mb-1">To</label>
        <input
          id="maxSalary"
          type="number"
          bind:value={formData.maxSalary}
          placeholder="Enter maximum salary"
          class="text-sm min-w-xs px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
        />
      </div>
    </div>
  </div>
</div>