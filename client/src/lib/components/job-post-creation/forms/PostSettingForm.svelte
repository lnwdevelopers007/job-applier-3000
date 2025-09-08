<script>
  let { formData = $bindable() } = $props();
  
  if (!formData.applicationRequirements) {
    formData.applicationRequirements = {
      resume: false,
      coverLetter: false,
      portfolio: false,
      linkedin: false
    };
  }
  
  if (!formData.notificationSettings) {
    formData.notificationSettings = {
      notifyOnApply: false,
      weeklyAppSummary: false,
      autoAcknowledge: false
    };
  }
</script>

<div class="space-y-6">
  <h2 class="text-lg font-medium text-gray-900 mb-6">Post Setting</h2>
  
  <!-- Posting Timeline -->
  <div>
    <h3 class="text-sm font-medium text-gray-900 mb-4">Posting Timeline</h3>
    <div class="grid grid-cols-2 gap-4">
      <div>
        <label for="postingOpenDate" class="block text-sm text-gray-700 mb-2">
          Posting Open Date <span class="text-red-500">*</span>
        </label>
        <input
          id="postingOpenDate"
          type="date"
          bind:value={formData.postingOpenDate}
          placeholder="DD/MM/YYYY"
          class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
        />
        <p class="mt-1 text-xs text-gray-500">When should this job posting go live?</p>
      </div>
      
      <div>
        <label for="postingCloseDate" class="block text-sm text-gray-700 mb-2">
          Posting Close Date
        </label>
        <input
          id="postingCloseDate"
          type="date"
          bind:value={formData.postingCloseDate}
          placeholder="DD/MM/YYYY"
          class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
        />
        <p class="mt-1 text-xs text-gray-500">When should applications close? (optional)</p>
      </div>
    </div>
  </div>

  <!-- Application Requirements -->
  <div>
    <h3 class="text-sm font-medium text-gray-900 mb-4">Application Requirements</h3>
    <div class="space-y-3">
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.applicationRequirements.resume}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Resume/CV required</span>
      </label>
      
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.applicationRequirements.coverLetter}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Cover letter required</span>
      </label>
      
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.applicationRequirements.portfolio}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Portfolio/work samples required</span>
      </label>
      
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.applicationRequirements.linkedin}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">LinkedIn profile required</span>
      </label>
    </div>
  </div>

  <!-- Screening Questions -->
  <div>
    <h3 class="text-sm font-medium text-gray-900 mb-2">Screening Questions</h3>
    <textarea
      bind:value={formData.screeningQuestions}
      placeholder="Add questions to help filter candidates."
      rows="3"
      class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 resize-none"
    ></textarea>
    <p class="mt-1 text-xs text-gray-500">Optional questions to help screen applicants before interviews</p>
  </div>

  <!-- Notification Settings -->
  <div>
    <h3 class="text-sm font-medium text-gray-900 mb-4">Notification Settings</h3>
    <div class="space-y-3">
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.notificationSettings.notifyOnApply}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Notify me when candidates apply</span>
      </label>
      
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.notificationSettings.weeklyAppSummary}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Send weekly application summary</span>
      </label>
      
      <label class="flex items-center">
        <input
          type="checkbox"
          bind:checked={formData.notificationSettings.autoAcknowledge}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
        />
        <span class="ml-2 text-sm text-gray-700">Auto-acknowledge applications to candidates</span>
      </label>
    </div>
  </div>
</div>