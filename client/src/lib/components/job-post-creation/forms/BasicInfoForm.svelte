<script>
  import TextInput from '$lib/components/forms/TextInput.svelte';
  import Select from '$lib/components/forms/Select.svelte';
  
  let { 
    formData = $bindable(),
    validationErrors = {},
    showValidationErrors = false
  } = $props();
  
  const currencyOptions = [
    { value: 'THB', label: 'THB' },
    { value: 'USD', label: 'USD' },
    { value: 'EUR', label: 'EUR' }
  ];
</script>

<div class="space-y-6">
  <h2 class="text-lg font-medium text-gray-900 mb-6">Basic Job Information</h2>
  
  <!-- Job Title -->
  <TextInput
    label="Job title"
    bind:value={formData.jobTitle}
    placeholder="Enter job title..."
    required={true}
    error={validationErrors.jobTitle}
    showError={showValidationErrors}
    class="max-w-lg"
  />

  <!-- Location -->
  <TextInput
    label="Location"
    bind:value={formData.location}
    placeholder="Enter location..."
    required={true}
    error={validationErrors.location}
    showError={showValidationErrors}
    class="max-w-lg"
  />

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
      <div class="w-22">
        <Select
          label="Currency"
          bind:value={formData.currency}
          options={currencyOptions}
          class="text-xs"
        />
      </div>
      
      <!-- From -->
      <div class="flex-shrink-0">
        <TextInput
          label="From"
          type="number"
          bind:value={formData.minSalary}
          placeholder="Enter minimum salary"
          required={true}
          error={validationErrors.minSalary}
          showError={showValidationErrors}
          class="max-w-xs"
        />
      </div>
      
      <!-- To -->
      <div class="flex-shrink-0">
        <TextInput
          label="To"
          type="number"
          bind:value={formData.maxSalary}
          placeholder="Enter maximum salary"
          required={true}
          error={validationErrors.maxSalary}
          showError={showValidationErrors}
          class="max-w-xs"
        />
      </div>
    </div>
  </div>
</div>