<script lang="ts">
	type Company = {
		id: string; // Used for keying in each loops
		name: string;
		logo?: string;
		logoStyle?: string;
		positions: number;
	};

	// eslint-disable-next-line svelte/no-unused-props
	let { company }: { company: Company } = $props();

	function getCompanyInitial(name: string): string {
		return name.charAt(0).toUpperCase();
	}
</script>

<div
	class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-sm hover:border-green-600 transition-all duration-200 cursor-pointer flex flex-col gap-2"
>
	<div
		class="w-12 h-12 rounded-md flex items-center justify-center text-xl font-semibold {company.logoStyle ||
			'bg-gradient-to-br from-gray-100 to-gray-200 text-gray-600'}"
	>
		{company.logo || getCompanyInitial(company.name)}
	</div>
	<h4 class="font-medium text-gray-900 text-sm mb-5">{company.name}</h4>
	<div class="flex">
		<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
			{company.positions} positions
		</span>
	</div>
</div>