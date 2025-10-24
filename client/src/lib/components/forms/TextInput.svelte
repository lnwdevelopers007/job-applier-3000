<script lang="ts">
	import type { ComponentType } from 'svelte';
	import { AlertCircle } from 'lucide-svelte';

	interface Props {
		label?: string;
		value: string;
		placeholder?: string;
		type?: string;
		readonly?: boolean;
		required?: boolean;
		icon?: string | null;
		iconComponent?: ComponentType | null;
		helpText?: string;
		error?: string;
		showError?: boolean;
		class?: string;
		onchange?: () => void;
	}

	let {
		label = '',
		value = $bindable(''),
		placeholder = '',
		type = 'text',
		readonly = false,
		required = false,
		icon = null,
		iconComponent = null,
		helpText = '',
		error = '',
		showError = false,
		class: className = '',
		onchange = () => {}
	}: Props = $props();

	// Check if we should show an icon (either SVG string or component)
	const hasIcon = $derived(icon || iconComponent);
	
	// Track if user has started typing to hide error
	let hideError = $state(false);
	
	// Hide error when user starts typing
	function handleInput() {
		hideError = true;
	}
	
	// Reset hideError when new validation occurs
	$effect(() => {
		if (showError && error) {
			hideError = false;
		}
	});
</script>

<div class={className}>
	{#if label}
		<label class="block text-sm font-medium text-gray-700 mb-1">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}
	
	<div class="relative">
		{#if hasIcon}
			<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
				{#if iconComponent}
					<svelte:component this={iconComponent} class="h-5 w-5 text-gray-400" />
				{:else if icon}
					<!-- eslint-disable-next-line svelte/no-at-html-tags -->
					{@html icon}
				{/if}
			</div>
		{/if}
		
		<input
			{type}
			bind:value
			{placeholder}
			{readonly}
			{required}
			onchange={onchange}
			oninput={handleInput}
			class="w-full {hasIcon ? 'pl-10' : 'px-4'} {showError && error && !hideError ? 'pr-10' : 'pr-4'} py-2 text-sm border rounded-lg focus:ring-1 focus:ring-gray-400 focus:border-transparent transition-all {readonly ? 'bg-gray-50 text-gray-500 cursor-not-allowed' : ''} {showError && error && !hideError ? 'border-red-500' : 'border-gray-300'}"
		/>
		
		{#if showError && error && !hideError}
			<!-- Error icon -->
			<div class="absolute inset-y-0 right-0 flex items-center pr-3">
				<AlertCircle class="w-5 h-5 text-white" fill="red" />
			</div>
			
			<!-- Floating error tooltip positioned below input -->
			<div class="absolute z-50 right-1 top-full mt-1 px-3 py-2 bg-red-500 text-white text-xs rounded-md shadow-sm whitespace-nowrap">
				<!-- Arrow pointing up to the alert icon -->
				<div class="absolute -top-1 right-4 w-2 h-2 bg-red-500 transform rotate-45"></div>
				{error}
			</div>
		{/if}
	</div>
	
	{#if helpText && (!showError || !error)}
		<p class="text-xs text-gray-500 mt-1">{helpText}</p>
	{/if}
	
	{#if error && !showError}
		<p class="text-xs text-red-500 mt-1">{error}</p>
	{/if}
</div>