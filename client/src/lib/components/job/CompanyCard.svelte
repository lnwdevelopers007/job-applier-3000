<script lang="ts">
  import { Users, Building } from 'lucide-svelte';
  
  type Company = {
    name: string;
    logo: string;
    industry: string;
    employees: string;
    description: string;
  };

  let {
    company
  }: {
    company: Company;
  } = $props();

  let expanded = $state(false);
</script>

<div class="flex items-start gap-3 bg-gray-100 rounded-lg p-6">
  <img src={company.logo} alt={company.name} class="w-10 h-10 rounded">
  <div class="flex-1">
    <h3 class="text-sm font-medium text-gray-900">{company.name}</h3>
    <p class="text-sm text-gray-600 mb-2">{company.industry}</p>
    <div class="flex items-center gap-2 text-sm text-gray-600 mb-3">
      <div class="flex items-center gap-1">
        <Users class="w-3.5 h-3.5" />
        <span>{company.employees}</span>
      </div>
      <span class="text-gray-400">•</span>
      <div class="flex items-center gap-1">
        <Building class="w-3.5 h-3.5" />
        <span>Enterprise</span>
      </div>
    </div>
    <p class="text-sm text-gray-700 leading-relaxed mb-2">
      {expanded ? company.description : company.description.slice(0, 120) + '...'}
    </p>
    <button 
      onclick={() => expanded = !expanded}
      class="text-sm text-green-600 hover:text-green-700 transition-colors font-medium"
    >
      {expanded ? 'Show less' : 'Show more'}
    </button>
  </div>
</div>