<script lang="ts">
	import { goto } from '$app/navigation';
	const currentYear = new Date().getFullYear();

	const footerLinks = {
		jobSeekers: [
			{ text: 'Browse Jobs', href: '#' },
			{ text: 'Companies', href: '#' },
			{ text: 'Career Advice', href: '#' },
			{ text: 'Salary Guide', href: '#' }
		],
		employers: [
			{ text: 'Post a Job', href: '#' },
			{ text: 'Pricing', href: '#' },
			{ text: 'Employer Resources', href: '#' }
		],
		connect: [
			{ text: 'About Us', href: '#' },
			{ text: 'Contact', href: '#' },
			{ text: 'Privacy Policy', href: '#' },
			{ text: 'Terms of Service', href: '#' }
		]
	};
</script>

<footer class="bg-white py-12 px-4 sm:px-6 lg:px-8 border-t border-gray-200">
	<div class="container mx-auto max-w-7xl">
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
			<div>
				<button 
					class="flex items-start space-x-1 mb-4 group cursor-pointer hover:opacity-80 transition-opacity focus:outline-none focus:opacity-80"
					onclick={() => goto('/')}
					onkeydown={(e) => e.key === 'Enter' && goto('/')}
					aria-label="Go to home page"
				>
					<h2 class="text-lg font-semibold text-black">Job Applier </h2>
					<h2 class="text-lg font-semibold text-green-700">3000</h2>
				</button>
				<p class="text-sm text-gray-600">
					Your gateway to amazing tech opportunities at KU Computer Engineering
				</p>
			</div>

			<div>
				<h4 class="text-gray-900 font-semibold mb-4">For Job Seekers</h4>
				<ul class="space-y-2">
					{#each footerLinks.jobSeekers as link}
						<li>
							<a href={link.href} class="text-sm text-gray-500 hover:text-gray-900 transition-colors">
								{link.text}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<div>
				<h4 class="text-gray-900 font-semibold mb-4">For Employers</h4>
				<ul class="space-y-2">
					{#each footerLinks.employers as link}
						<li>
							<a href={link.href} class="text-sm text-gray-500 hover:text-gray-900 transition-colors">
								{link.text}
							</a>
						</li>
					{/each}
				</ul>
			</div>

			<div>
				<h4 class="text-gray-900 font-semibold mb-4">Connect</h4>
				<ul class="space-y-2">
					{#each footerLinks.connect as link}
						<li>
							<a href={link.href} class="text-sm text-gray-500 hover:text-gray-900 transition-colors">
								{link.text}
							</a>
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<div class="border-t border-gray-200 pt-8 text-center">
			<p class="text-sm text-gray-500">© {currentYear} Job Applier 3000. All rights reserved.</p>
		</div>
	</div>
</footer>