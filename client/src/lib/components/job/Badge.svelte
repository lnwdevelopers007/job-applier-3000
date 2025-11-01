<script lang="ts">
  type BadgeVariant = 'primary' | 'secondary' | 'warning' | 'danger' | 'success' | 'info' | 'purple' | 'indigo' | 'pink' | 'teal' | 'orange' | 'cyan' | 'lime' | 'rose';
  
  let {
    variant = 'secondary',
    size = 'sm',
    text,
    children
  }: {
    variant?: BadgeVariant;
    size?: 'xs' | 'sm' | 'md';
    text?: string;
    children?: import('svelte').Snippet;
  } = $props();
  
  // Format text to title case
  function formatLabel(text: string): string {
    if (!text) return '';
    return text
      .toLowerCase()
      .split('-')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join('-');
  }

  const variantClasses = {
    primary: 'bg-green-50 text-green-700 ring-green-600/15',
    secondary: 'bg-gray-100 text-gray-700 ring-gray-600/15',
    warning: 'bg-yellow-50 text-yellow-700 ring-yellow-600/15',
    danger: 'bg-red-50 text-red-700 ring-red-600/15',
    success: 'bg-green-50 text-green-700 ring-green-600/15',
    info: 'bg-sky-50 text-sky-700 ring-sky-600/15',
    purple: 'bg-purple-50 text-purple-700 ring-purple-600/15',
    indigo: 'bg-indigo-50 text-indigo-700 ring-indigo-600/15',
    pink: 'bg-pink-50 text-pink-700 ring-pink-600/15',
    teal: 'bg-teal-50 text-teal-700 ring-teal-600/15',
    orange: 'bg-orange-50 text-orange-700 ring-orange-600/15',
    cyan: 'bg-cyan-50 text-cyan-700 ring-cyan-600/15',
    lime: 'bg-lime-50 text-lime-700 ring-lime-600/15',
    rose: 'bg-rose-50 text-rose-700 ring-rose-600/15'
  };

  const sizeClasses = {
    xs: 'px-2 py-0.5 text-xs',
    sm: 'px-2.5 py-1 text-xs',
    md: 'px-2.5 py-1 text-sm'
  };
</script>

<span class="inline-flex items-center rounded-md font-medium {variantClasses[variant]} {sizeClasses[size]} {variant !== 'secondary' ? 'ring-1 ring-inset' : ''}">
  {#if text}
    {formatLabel(text)}
  {:else if children}
    {@render children()}
  {/if}
</span>