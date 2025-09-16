<script>
  let { 
    type = 'button', 
    variant = 'primary', 
    fullWidth = true, 
    disabled = false, 
    onclick = undefined 
  } = $props();
  
  const variants = {
    primary: 'bg-green-600 hover:bg-green-700 text-white focus:ring-2 focus:ring-green-500 focus:ring-offset-2',
    secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2',
    ghost: 'bg-transparent hover:bg-gray-100 text-gray-700'
  };
</script>

<button
  {type}
  {disabled}
  {onclick}
  class="{fullWidth ? 'w-full' : ''} py-3 px-4 {variants[variant]} text-sm font-semibold rounded-lg transition-colors hover:cursor-pointer focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
>
  <slot />
</button>