<script>
  import { ChevronLeft } from 'lucide-svelte';
  
  let { onBack = () => window.history.back() } = $props();
</script>

<div class="flex items-center space-x-1 mb-6">
  <ChevronLeft 
    class="w-5 h-5 text-gray-600 cursor-pointer hover:text-gray-800 transition-colors" 
    onclick={onBack}
  />
  <h2 class="text-lg font-semibold text-black">Job Applier</h2>
  <h2 class="text-lg font-semibold text-green-700">3000</h2>
</div>