<script>
  import { Check } from 'lucide-svelte';

  let { 
    steps, 
    currentStep = $bindable()
  } = $props();

  export function handleProgress(stepIncrement) {
    if (stepIncrement === 1) {
      currentStep++;
      if (currentStep > steps.length) {
        currentStep = steps.length;
      }
    } else {
      currentStep--;
      if (currentStep < 1) {
        currentStep = 1;
      }
    }
  }
</script>

<div class="px-8 py-6 mb-3">
  <div class="max-w-4xl relative">
    <!-- Steps with integrated progress lines -->
    <div class="flex items-start">
      {#each steps as step, i}
        <div class="flex items-start">
          <!-- Step circle and text -->
          <div class="flex flex-col items-start">
            <div class="w-6 h-6 {currentStep > i + 1 ? 'bg-green-600' : currentStep === i + 1 ? 'bg-green-600' : 'bg-gray-300'} rounded-full flex items-center justify-center mb-2 transition-colors duration-300">
              {#if currentStep > i + 1}
                <Check class="w-3 h-3 text-white" strokeWidth={4}/>
              {:else if currentStep === i + 1}
                <div class="w-4 h-4 bg-green-600 rounded-full border-3 border-white"></div>
              {:else}
                <div class="w-4 h-4 bg-white rounded-full"></div>
              {/if}
            </div>
            <span class="text-sm {currentStep >= i + 1 ? 'font-medium text-gray-900' : 'text-gray-500'} transition-colors duration-300 w-24">{step.name}</span>
          </div>
          
          <!-- Connecting line immediately after each step (except the last) -->
          {#if i < steps.length - 1}
            <div class="w-48 h-1 rounded-xl {currentStep > i + 1 ? 'bg-green-600' : 'bg-gray-300'} mt-3 -ml-15 mr-3 flex-shrink-0 transition-colors duration-300"></div>
          {/if}
        </div>
      {/each}
    </div>
  </div>
</div>