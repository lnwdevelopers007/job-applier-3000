<script>
	let {
		label = '',
		value = $bindable(''),
		options = [],
		placeholder = 'Select...',
		required = false,
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
	
	<select
		bind:value
		{required}
		onchange={onchange}
		class="w-full px-4 py-2 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all {error ? 'border-red-300' : ''}"
	>
		<option value="" disabled selected>{placeholder}</option>
		{#each options as option}
			{#if typeof option === 'string'}
				<option value={option}>{option}</option>
			{:else}
				<option value={option.value}>{option.label}</option>
			{/if}
		{/each}
	</select>
	
	{#if helpText}
		<p class="text-xs text-gray-500 mt-1">{helpText}</p>
	{/if}
	
	{#if error}
		<p class="text-xs text-red-500 mt-1">{error}</p>
	{/if}
</div>