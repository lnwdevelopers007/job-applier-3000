<script>
	let {
		benefits = $bindable([]),
		onSave = () => {}
	} = $props();
</script>

<div class="divide-y divide-gray-200">
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Health & wellness</label>
				<p class="text-xs text-gray-500 mt-1">Medical, dental, and wellness benefits</p>
			</div>
			<div class="col-span-2">
				<div class="space-y-3">
					{#each ['Health insurance', 'Dental insurance', 'Vision insurance', 'Gym membership'] as benefit}
						<label class="flex items-center">
							<input 
								type="checkbox" 
								checked={benefits?.includes(benefit)}
								onchange={(e) => {
									if (!benefits) benefits = [];
									if (e.target.checked) {
										benefits = [...benefits, benefit];
									} else {
										benefits = benefits.filter(b => b !== benefit);
									}
								}}
								class="mr-2 text-green-600 rounded"
							/>
							<span class="text-sm text-gray-700">{benefit}</span>
						</label>
					{/each}
				</div>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Work-life balance</label>
				<p class="text-xs text-gray-500 mt-1">Flexibility and time-off benefits</p>
			</div>
			<div class="col-span-2">
				<div class="space-y-3">
					{#each ['Remote work', 'Flexible hours', 'Unlimited PTO', 'Parental leave'] as benefit}
						<label class="flex items-center">
							<input 
								type="checkbox" 
								checked={benefits?.includes(benefit)}
								onchange={(e) => {
									if (!benefits) benefits = [];
									if (e.target.checked) {
										benefits = [...benefits, benefit];
									} else {
										benefits = benefits.filter(b => b !== benefit);
									}
								}}
								class="mr-2 text-green-600 rounded"
							/>
							<span class="text-sm text-gray-700">{benefit}</span>
						</label>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>