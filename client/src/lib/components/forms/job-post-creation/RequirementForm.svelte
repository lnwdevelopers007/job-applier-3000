<script>
  let { formData = $bindable() } = $props();
  
  let skillInput = $state('');
  
  if (!formData.requiredSkills) {
    formData.requiredSkills = [];
  }
  
  function addSkill(event) {
    if (event.key === 'Enter' && skillInput.trim()) {
      event.preventDefault();
      if (!formData.requiredSkills.includes(skillInput.trim())) {
        formData.requiredSkills = [...formData.requiredSkills, skillInput.trim()];
      }
      skillInput = '';
    }
  }
  
  function removeSkill(skill) {
    formData.requiredSkills = formData.requiredSkills.filter(s => s !== skill);
  }
</script>

<div class="space-y-6">
  <h2 class="text-lg font-medium text-gray-900 mb-6">Requirement</h2>
  
  <!-- Years of Experience and Education Level (side by side) -->
  <div class="grid grid-cols-2 gap-4">
    <div>
      <label for="yearsOfExperience" class="block text-sm font-medium text-gray-700 mb-2">
        Years of Experience
      </label>
      <select
        id="yearsOfExperience"
        bind:value={formData.yearsOfExperience}
        class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 bg-white"
      >
        <option value="">Select Experience</option>
        <option value="entry">Entry Level (0-2 years)</option>
        <option value="junior">Junior (2-4 years)</option>
        <option value="mid">Mid-Level (4-6 years)</option>
        <option value="senior">Senior (6-10 years)</option>
        <option value="expert">Expert (10+ years)</option>
      </select>
    </div>
    
    <div>
      <label for="educationLevel" class="block text-sm font-medium text-gray-700 mb-2">
        Education Level
      </label>
      <select
        id="educationLevel"
        bind:value={formData.educationLevel}
        class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400 bg-white"
      >
        <option value="">Select Education</option>
        <option value="highschool">High School</option>
        <option value="associate">Associate Degree</option>
        <option value="bachelor">Bachelor's Degree</option>
        <option value="master">Master's Degree</option>
        <option value="phd">PhD</option>
        <option value="other">Other</option>
      </select>
    </div>
  </div>
  
  <!-- Required Skills & Technologies -->
  <div>
    <label for="requiredSkills" class="block text-sm font-medium text-gray-700 mb-2">
      Required Skills & Technologies
    </label>
    <div class="relative">
      <input
        id="requiredSkills"
        type="text"
        bind:value={skillInput}
        onkeydown={addSkill}
        placeholder="Type a skill and press Enter to add"
        class="w-full text-sm px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-gray-500 focus:border-gray-400"
      />
      
      <!-- Skill Tags -->
      {#if formData.requiredSkills && formData.requiredSkills.length > 0}
        <div class="flex flex-wrap gap-2 mt-3">
          {#each formData.requiredSkills as skill, index (index)}
            <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm bg-gray-100 text-gray-700">
              {skill}
              <button
                type="button"
                onclick={() => removeSkill(skill)}
                class="ml-1 text-gray-500 hover:text-gray-700 focus:outline-none"
              >
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </span>
          {/each}
        </div>
      {/if}
    </div>
    <p class="mt-2 text-xs text-gray-500">Add individual skills as tags. These help with job matching and filtering.</p>
  </div>
</div>