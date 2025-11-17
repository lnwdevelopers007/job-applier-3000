<script lang="ts">
	interface Props {
		name: string;
		email: string;
		avatarURL?: string;
	}

	let { name, email, avatarURL }: Props = $props();
</script>

<div class="flex items-center">
	<div class="flex-shrink-0 h-10 w-10">
		{#if avatarURL}
			<img class="h-10 w-10 rounded-full object-cover" src={avatarURL} alt={name} />
		{:else}
			<div class="h-10 w-10 rounded-full bg-gray-300 flex items-center justify-center">
				<span class="text-sm font-medium text-gray-700">{name.charAt(0).toUpperCase()}</span>
			</div>
		{/if}
	</div>
	<div class="ml-4">
		<div class="text-sm font-medium text-gray-900">{name}</div>
		<div class="text-sm text-gray-500">{email}</div>
	</div>
</div>