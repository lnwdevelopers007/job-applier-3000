<script lang="ts">
	interface Props {
		title: string;
		company?: string;
		companyLogo?: string;
	}

	let { title, company, companyLogo }: Props = $props();
</script>

<div class="flex items-center gap-3">
	{#if companyLogo}
		<img src={companyLogo} alt={company || ''} class="h-12 w-12 rounded-lg p-1 object-contain" />
	{/if}
	<div>
		{#if company}
			<div class="text-xs text-gray-500 truncate">{company}</div>
		{/if}
		<div class="text-sm font-medium text-gray-900 truncate">{title}</div>
	</div>
</div>