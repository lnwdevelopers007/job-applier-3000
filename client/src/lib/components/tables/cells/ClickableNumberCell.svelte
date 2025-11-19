<script lang="ts">
	interface Props {
		value: number;
		onClick: () => void;
	}

	let { value, onClick }: Props = $props();
</script>

{#if value === 0}
	<span class="text-gray-500">0</span>
{:else}
	<button
		onclick={(e) => {
			e.stopPropagation();
			onClick();
		}}
		class="text-gray-600 hover:text-gray-800 hover:underline font-medium cursor-pointer transition-colors"
	>
		{value}
	</button>
{/if}