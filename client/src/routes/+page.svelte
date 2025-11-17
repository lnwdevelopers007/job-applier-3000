<script lang="ts">
	import { onMount } from 'svelte';
	import HeroSection from '$lib/components/home/HeroSection.svelte';
	import JobCard from '$lib/components/home/JobCard.svelte';
	import CareerPathCard from '$lib/components/home/CareerPathCard.svelte';
	import CTASection from '$lib/components/home/CTASection.svelte';
	import Footer from '$lib/components/home/Footer.svelte';
	import { ArrowRight, Code, ChartLine, Brush, Shield, Smartphone, Cloud, Bot, Gamepad2, Users, FileText, Bell, GraduationCap } from 'lucide-svelte';
	import { isAuthenticated, navigateWithAuth } from '$lib/utils/auth';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { fetchCompanyPublicInfo, DEFAULT_COMPANY_NAME, DEFAULT_COMPANY_LOGO } from '$lib/utils/fetcher';

	let showAuthModal = $state(false);
	let authModalTitle = $state('Sign in required');
	let authModalDescription = $state('Please log in or sign up to continue');

	function handleViewAllJobs() {
		if (isAuthenticated()) {
			navigateWithAuth('/app/jobs');
		} else {
			authModalTitle = 'Sign in to view all jobs';
			authModalDescription = 'Please log in or sign up to browse all available job opportunities';
			showAuthModal = true;
			sessionStorage.setItem('pendingNavigation', '/app/jobs');
		}
	}

	function handleAuthModalClose() {
		showAuthModal = false;
		sessionStorage.removeItem('pendingNavigation');
	}

	let recentJobs = $state<Array<{
		id: string;
		company: string;
		companyLogo: string;
		title: string;
		location: string;
		locationType: 'on-site' | 'remote' | 'hybrid' | null;
		minSalary: number;
		maxSalary: number;
		currency: string;
		type: string | null;
		tags: string[];
		postedAt: string;
		badge: { text: string; type: string } | null;
		logoStyle: string;
	}>>([]);
	let isLoadingJobs = $state(true);

	async function fetchRecentJobs() {
		try {
			const backendUrl = import.meta.env.VITE_BACKEND || 'http://localhost:8080';
			const res = await fetch(`${backendUrl}/jobs/public/latest`);
			if (res.ok) {
				const data = await res.json();

				const jobPromises = data.slice(0, 3).map(async (job: {
					id: string;
					title?: string;
					location?: string;
					workArrangement?: string;
					minSalary?: number;
					maxSalary?: number;
					currency?: string;
					workType?: string;
					requiredSkills?: string;
					postOpenDate?: string;
					companyID?: string;
				}) => {
					// Fetch minimal public company info (with custom logo support)
					let companyName = DEFAULT_COMPANY_NAME;
					let companyLogo = DEFAULT_COMPANY_LOGO;
					
					if (job.companyID) {
						try {
							const companyInfo = await fetchCompanyPublicInfo(job.companyID);
							companyName = companyInfo.name;
							companyLogo = companyInfo.profileImage;
						} catch (err) {
							console.warn('Failed to fetch company info for job', job.id, err);
						}
					}

					return {
						id: job.id,
						company: companyName,
						companyLogo: companyLogo,
						title: job.title || "Untitled Position",
						location: job.location || "N/A",
						locationType: job.workArrangement?.toLowerCase() === 'remote' ? 'remote' : 
									  job.workArrangement?.toLowerCase() === 'hybrid' ? 'hybrid' : 
									  job.workArrangement?.toLowerCase() === 'on-site' ? 'on-site' : null,
						minSalary: job.minSalary || 0,
						maxSalary: job.maxSalary || 0,
						currency: job.currency || 'THB',
						type: job.workType ? job.workType.charAt(0).toUpperCase() + job.workType.slice(1) : '',
						tags: job.requiredSkills 
							? job.requiredSkills.split(',').map((skill: string) => skill.trim()).slice(0, 3)
							: [],
						postedAt: job.postOpenDate && job.postOpenDate !== "0001-01-01T00:00:00Z"
							? `Posted ${Math.floor((Date.now() - new Date(job.postOpenDate).getTime()) / (1000 * 60 * 60 * 24))} days ago`
							: 'Recently posted',
						badge: (() => {
							if (job.workArrangement?.toLowerCase() === 'remote') {
								return { text: 'Remote', type: 'remote' };
							}
							if (job.workType?.toLowerCase().includes('intern') || job.workType?.toLowerCase() === 'casual') {
								return { text: 'Internship', type: 'internship' };
							}
							if (job.postOpenDate && job.postOpenDate !== "0001-01-01T00:00:00Z") {
								const daysAgo = Math.floor((Date.now() - new Date(job.postOpenDate).getTime()) / (1000 * 60 * 60 * 24));
								if (daysAgo <= 3) {
									return { text: 'New', type: 'new' };
								}
							}
							return null;
						})(),
						logoStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600'
					};
				});

				recentJobs = await Promise.all(jobPromises);
			}
		} catch (error) {
			console.error('Failed to fetch recent jobs:', error);
			recentJobs = [];
		} finally {
			isLoadingJobs = false;
		}
	}

	onMount(() => {
		fetchRecentJobs();
	});


	// App features
	const appFeatures = [
		{
			id: '1',
			title: 'KU-Focused Platform',
			description: 'Exclusively designed for CPSK students with opportunities tailored to your Computer Engineering background and skills.',
			icon: GraduationCap,
			iconStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600'
		},
		{
			id: '2', 
			title: 'One-Click Applications',
			description: 'Apply to multiple positions effortlessly with your pre-filled profile and documents.',
			icon: FileText,
			iconStyle: 'bg-gradient-to-br from-green-100 to-green-200 text-green-600'
		},
		{
			id: '3',
			title: 'Real-time Notifications',
			description: 'Stay updated on application status, new opportunities, and company responses instantly.',
			icon: Bell,
			iconStyle: 'bg-gradient-to-br from-purple-100 to-purple-200 text-purple-600'
		},
		{
			id: '4',
			title: 'Direct Company Connect',
			description: 'Connect directly with hiring managers and companies actively seeking KU talent.',
			icon: Users,
			iconStyle: 'bg-gradient-to-br from-orange-100 to-orange-200 text-orange-600'
		}
	];

	// TODO: Replace with actual API fetch from backend
	const careerPaths = [
		{ id: '1', title: 'Software Dev', icon: Code, iconStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600', positions: 1 },
		{ id: '2', title: 'Data Science', icon: ChartLine, iconStyle: 'bg-gradient-to-br from-green-100 to-green-200 text-green-600', positions: 2 },
		{ id: '3', title: 'Design & UX', icon: Brush, iconStyle: 'bg-gradient-to-br from-purple-100 to-purple-200 text-purple-600', positions: 3 },
		{ id: '4', title: 'Cybersecurity', icon: Shield, iconStyle: 'bg-gradient-to-br from-red-100 to-red-200 text-red-600', positions: 4 },
		{ id: '5', title: 'Mobile Dev', icon: Smartphone, iconStyle: 'bg-gradient-to-br from-orange-100 to-orange-200 text-orange-600', positions: 5 },
		{ id: '6', title: 'Cloud & DevOps', icon: Cloud, iconStyle: 'bg-gradient-to-br from-cyan-100 to-cyan-200 text-cyan-600', positions: 6 },
		{ id: '7', title: 'AI & ML', icon: Bot, iconStyle: 'bg-gradient-to-br from-yellow-100 to-yellow-200 text-yellow-600', positions: 7 },
		{ id: '8', title: 'Game Dev', icon: Gamepad2, iconStyle: 'bg-gradient-to-br from-indigo-100 to-indigo-200 text-indigo-600', positions: 8 }
	];
</script>

<svelte:head>
	<title>Job Applier 3000</title>
	<meta name="description" content="Connect with leading tech companies and discover opportunities tailored for KU Computer Engineering students" />
</svelte:head>

<div class="min-h-screen">
	<main>
		<HeroSection />

		<!-- Recent Opportunities -->
		<section class="py-16 px-4 sm:px-6 lg:px-8 bg-white">
			<div class="container mx-auto max-w-7xl">
				<div class="flex items-center justify-between mb-8">
					<div>
						<h2 class="text-3xl font-bold text-gray-900">Recent Opportunities</h2>
						<p class="text-gray-600 mt-2">Latest job postings from industry leaders</p>
					</div>
					<button onclick={handleViewAllJobs} class="flex items-center gap-2 text-green-600 font-medium hover:text-green-700 cursor-pointer transition-colors">
						View all jobs
						<ArrowRight class="w-4 h-4" />
					</button>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#if isLoadingJobs}
						<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
						{#each Array(3) as _, index (index)}
							<JobCard skeleton={true} />
						{/each}
					{:else if recentJobs.length > 0}
						{#each recentJobs as job (job.id)}
							<JobCard {job} />
						{/each}
					{:else}
						<div class="col-span-full text-center py-8">
							<div class="max-w-md mx-auto">
								<h3 class="text-lg font-medium text-gray-900 mb-2">No recent opportunities</h3>
								<p class="text-gray-500 mb-4">Check back later for new job postings, or explore all available positions.</p>
								<button 
									onclick={handleViewAllJobs}
									class="inline-flex items-center gap-2 px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
								>
									Browse all jobs
									<ArrowRight class="w-4 h-4" />
								</button>
							</div>
						</div>
					{/if}
				</div>
			</div>
		</section>

		
		<!-- Career Paths -->
		<section class="py-16 px-4 sm:px-6 lg:px-8 bg-slate-50">
			<div class="container mx-auto max-w-7xl">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900">Explore Career Paths</h2>
					<p class="text-gray-600 mt-2 max-w-2xl mx-auto">
						Discover opportunities across different technology domains and find the perfect match for your expertise
					</p>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
					{#each careerPaths as path (path.id)}
						<CareerPathCard {path} />
					{/each}
				</div>
			</div>
		</section>

		<!-- App Features -->
		<section class="py-16 px-4 sm:px-6 lg:px-8 bg-white">
			<div class="container mx-auto max-w-7xl">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900">Why Choose Job Applier 3000?</h2>
					<p class="text-gray-600 mt-2 max-w-2xl mx-auto">
						Streamline your job search with powerful features designed specifically for CPSK students
					</p>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
					{#each appFeatures as feature (feature.id)}
						<div class="bg-white rounded-xl p-6 border border-gray-200 hover:shadow-sm transition-shadow duration-200">
							<div class="flex items-center justify-center w-12 h-12 rounded-lg {feature.iconStyle} mb-4">
								<feature.icon class="w-6 h-6" />
							</div>
							<h3 class="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
							<p class="text-gray-600 text-sm leading-relaxed">{feature.description}</p>
						</div>
					{/each}
				</div>
			</div>
		</section>


		<CTASection />
	</main>

	<Footer />
</div>

<!-- Auth Modal -->
<AuthModal
	bind:isOpen={showAuthModal}
	onClose={handleAuthModalClose}
	title={authModalTitle}
	description={authModalDescription}
/>