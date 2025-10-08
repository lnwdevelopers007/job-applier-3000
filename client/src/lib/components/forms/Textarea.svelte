<script>
	let {
		label = '',
		value = $bindable(''),
		placeholder = '',
		rows = 4,
		maxLength = null,
		required = false,
		helpText = '',
		error = '',
		class: className = '',
		onchange = () => {}
	} = $props();
	
	let charactersLeft = $derived(maxLength ? maxLength - value.length : null);
</script>

<div class={className}>
	{#if label}
		<label class="block text-sm font-medium text-gray-700 mb-1">{label}</label>
	{/if}
	
	<textarea
		bind:value
		{placeholder}
		{rows}
		{required}
		maxlength={maxLength}
		onchange={onchange}
		class="w-full px-4 py-2 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all resize-none {error ? 'border-red-300' : ''}"
	></textarea>
	
	{#if charactersLeft !== null}
		<p class="text-xs text-gray-500 mt-1">{charactersLeft} characters remaining</p>
	{:else if helpText}
		<p class="text-xs text-gray-500 mt-1">{helpText}</p>
	{/if}
	
	{#if error}
		<p class="text-xs text-red-500 mt-1">{error}</p>
	{/if}
</div>