<script>
	let {
		label = '',
		value = $bindable(''),
		placeholder = '',
		type = 'text',
		readonly = false,
		required = false,
		icon = null,
		helpText = '',
		error = '',
		class: className = '',
		onchange = () => {}
	} = $props();
</script>

<div class={className}>
	{#if label}
		<label class="block text-sm font-medium text-gray-700 mb-1">{label}</label>
	{/if}
	
	<div class="relative">
		{#if icon}
			<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
				{@html icon}
			</div>
		{/if}
		
		<input
			{type}
			bind:value
			{placeholder}
			{readonly}
			{required}
			onchange={onchange}
			class="w-full {icon ? 'pl-10' : 'px-4'} pr-4 py-2 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all {readonly ? 'bg-gray-50 text-gray-500 cursor-not-allowed' : ''} {error ? 'border-red-300' : ''}"
		/>
	</div>
	
	{#if helpText}
		<p class="text-xs text-gray-500 mt-1">{helpText}</p>
	{/if}
	
	{#if error}
		<p class="text-xs text-red-500 mt-1">{error}</p>
	{/if}
</div>