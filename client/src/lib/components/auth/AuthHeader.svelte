<script>
  import { ChevronLeft } from 'lucide-svelte';
  
  export let onBack = () => window.history.back();
</script>

<div class="flex items-center space-x-1 mb-6">
  <ChevronLeft 
    class="w-5 h-5 text-gray-600 cursor-pointer" 
    on:click={onBack}
  />
  <h2 class="text-lg font-semibold text-black">Job Applier</h2>
  <h2 class="text-lg font-semibold text-green-700">3000</h2>
</div>