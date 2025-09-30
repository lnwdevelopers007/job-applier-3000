<script lang="ts">
	type Company = {
		id: string;
		name: string;
		logo?: string;
		logoStyle?: string;
		positions: number;
	};

	let { company }: { company: Company } = $props();

	function getCompanyInitial(name: string): string {
		return name.charAt(0).toUpperCase();
	}
</script>

<div
	class="bg-white border border-gray-200 rounded-xl p-6 hover:shadow-sm hover:border-green-600 transition-all duration-200 cursor-pointer flex flex-col items-center gap-3"
>
	<div
		class="w-14 h-14 rounded-lg flex items-center justify-center text-2xl font-semibold {company.logoStyle ||
			'bg-gradient-to-br from-gray-100 to-gray-200 text-gray-600'}"
	>
		{company.logo || getCompanyInitial(company.name)}
	</div>
	<h4 class="font-semibold text-gray-900">{company.name}</h4>
	<p class="text-sm text-gray-600 font-medium">{company.positions} positions</p>
</div>