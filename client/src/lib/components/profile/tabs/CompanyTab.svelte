<script>
	import TextInput from '$lib/components/forms/TextInput.svelte';
	import Select from '$lib/components/forms/Select.svelte';
	import Textarea from '$lib/components/forms/Textarea.svelte';
	import FileUpload from '$lib/components/forms/FileUpload.svelte';
	import { Globe } from 'lucide-svelte';
	
	let {
		companyData = $bindable({}),
		onSave = () => {}
	} = $props();
</script>

<div class="divide-y divide-gray-200">
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Company name</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={companyData.companyName}
					placeholder="Enter your company name"
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Company logo</label>
				<p class="text-xs text-gray-500 mt-1">This will be displayed on job postings and your profile</p>
			</div>
			<div class="col-span-2">
				<FileUpload 
					currentImage={companyData.companyLogo}
					maxSize="2MB"
					onFileSelect={(file) => {
						companyData.companyLogoFile = file;
					}}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Industry</label>
			<div class="col-span-1">
				<Select 
					bind:value={companyData.industry}
					placeholder="Select industry..."
					options={[
						'Technology',
						'Healthcare',
						'Finance',
						'Education',
						'Manufacturing',
						'Retail',
						'Consulting',
						'Other'
					]}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Company size</label>
			<div class="col-span-1">
				<Select 
					bind:value={companyData.companySize}
					placeholder="Select size..."
					options={[
						'1-10 employees',
						'11-50 employees',
						'51-200 employees',
						'201-500 employees',
						'501-1000 employees',
						'1000+ employees'
					]}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Founded</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={companyData.foundedYear}
					type="number"
					placeholder="Enter year founded"
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">Headquarters</label>
				<p class="text-xs text-gray-500 mt-1">Primary company location</p>
			</div>
			<div class="col-span-1">
				<TextInput 
					bind:value={companyData.headquarters}
					placeholder="City, State/Country"
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-start">
			<div>
				<label class="text-sm font-medium text-gray-700">About us</label>
				<p class="text-xs text-gray-500 mt-1">Brief overview of your company and what you do</p>
			</div>
			<div class="col-span-1">
				<Textarea 
					bind:value={companyData.aboutCompany}
					placeholder="Tell candidates about your company..."
					maxLength={1000}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">Website</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={companyData.companyWebsite}
					type="url"
					placeholder="https://yourcompany.com"
					iconComponent={Globe}
				/>
			</div>
		</div>
	</div>
	
	<div class="px-8 py-5">
		<div class="grid grid-cols-3 gap-8 items-center">
			<label class="text-sm font-medium text-gray-700">LinkedIn</label>
			<div class="col-span-1">
				<TextInput 
					bind:value={companyData.companyLinkedin}
					type="url"
					placeholder="linkedin.com/company/yourcompany"
					icon='<svg class="h-5 w-5" fill="#0A66C2" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>'
				/>
			</div>
		</div>
	</div>
</div>