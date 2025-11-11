<script lang="ts">
  import { CircleCheck, CircleX, FileText, CircleDot } from 'lucide-svelte';

  interface Props {
    status: 'Active' | 'Closed' | 'Draft' | string;
    size?: 'xs' | 'sm' | 'md';
    showIcon?: boolean;
  }

  let { 
    status, 
    size = 'sm',
    showIcon = true
  }: Props = $props();

  // Status color configurations
  const statusConfig = {
    Active: {
      bg: 'bg-green-50',
      text: 'text-green-700',
      ring: 'ring-green-600/20'
    },
    Closed: {
      bg: 'bg-red-50',
      text: 'text-red-700',
      ring: 'ring-red-600/20'
    },
    Draft: {
      bg: 'bg-gray-50',
      text: 'text-gray-700',
      ring: 'ring-gray-600/20'
    },
    default: {
      bg: 'bg-gray-50',
      text: 'text-gray-700',
      ring: 'ring-gray-600/20'
    }
  };

  // Size configurations
  const sizeConfig = {
    xs: 'px-2 py-0.5 text-xs gap-1',
    sm: 'px-2.5 py-1.5 text-xs gap-1',
    md: 'px-3 py-1.5 text-sm gap-1.5'
  };

  // Icon sizes
  const iconSizes = {
    xs: 'w-3 h-3',
    sm: 'w-3.5 h-3.5',
    md: 'w-4 h-4'
  };

  const config = $derived(statusConfig[status as keyof typeof statusConfig] || statusConfig.default);
  const sizeClasses = $derived(sizeConfig[size]);
  const iconSize = $derived(iconSizes[size]);
</script>

<span class="inline-flex items-center rounded-full font-medium ring-1 ring-inset {config.bg} {config.text} {config.ring} {sizeClasses}">
  {#if showIcon}
    {#if status === 'Active'}
      <CircleCheck class="{iconSize} fill-green-700 stroke-white" />
    {:else if status === 'Closed'}
      <CircleX class="{iconSize} fill-red-700 stroke-white" />
    {:else if status === 'Draft'}
      <FileText class="{iconSize} fill-gray-700 stroke-white" />
    {:else}
      <CircleDot class="{iconSize} fill-gray-700 stroke-white" />
    {/if}
  {/if}
  {status}
</span>