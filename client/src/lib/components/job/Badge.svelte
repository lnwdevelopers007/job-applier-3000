<script lang="ts">
  type BadgeVariant = 'primary' | 'secondary' | 'warning' | 'success' | 'info';
  
  let {
    variant = 'secondary',
    size = 'sm',
    children
  }: {
    variant?: BadgeVariant;
    size?: 'xs' | 'sm' | 'md';
    children: import('svelte').Snippet;
  } = $props();

  const variantClasses = {
    primary: 'bg-green-50 text-green-700 ring-green-600/15',
    secondary: 'bg-gray-100 text-gray-700',
    warning: 'bg-amber-50 text-amber-700 ring-amber-600/15',
    success: 'bg-green-50 text-green-700 ring-green-600/15',
    info: 'bg-sky-50 text-sky-700 ring-sky-600/15'
  };

  const sizeClasses = {
    xs: 'px-2 py-0.5 text-xs',
    sm: 'px-2.5 py-1 text-xs',
    md: 'px-2.5 py-1 text-sm'
  };
</script>

<span class="inline-flex items-center rounded-md font-medium {variantClasses[variant]} {sizeClasses[size]} {variant !== 'secondary' ? 'ring-1 ring-inset' : ''}">
  {@render children()}
</span>