<script lang="ts">
	import { onMount } from 'svelte';
	import HeroSection from '$lib/components/home/HeroSection.svelte';
	import JobCard from '$lib/components/home/JobCard.svelte';
	import CompanyCard from '$lib/components/home/CompanyCard.svelte';
	import CareerPathCard from '$lib/components/home/CareerPathCard.svelte';
	import CTASection from '$lib/components/home/CTASection.svelte';
	import Footer from '$lib/components/home/Footer.svelte';
	import LayoutHeader from '$lib/components/LayoutHeader.svelte';
	import { ArrowRight, Code, ChartLine, Brush, Shield, Smartphone, Cloud, Bot, Gamepad2, ChevronLeft, ChevronRight, Briefcase } from 'lucide-svelte';
	import { isAuthenticated } from '$lib/utils/auth';
	import AuthModal from '$lib/components/ui/AuthModal.svelte';
	import { goto } from '$app/navigation';

	let showAuthModal = $state(false);
	let authModalTitle = $state('Sign in required');
	let authModalDescription = $state('Please log in or sign up to continue');

	function handleViewAllJobs() {
		if (isAuthenticated()) {
			goto('/app/jobs');
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

	let recentJobs = $state([]);
	let isLoadingJobs = $state(true);

	async function fetchRecentJobs() {
		try {
			const res = await fetch('/jobs/query');
			if (res.ok) {
				const data = await res.json();
				// Get the 3 most recent jobs and map to home page format
				recentJobs = data.slice(0, 3).map((job: any) => ({
					id: job.id,
					company: "Unknown Company",
					title: job.title || "Untitled Position",
					location: job.location || "N/A",
					locationType: job.workArrangement?.toLowerCase() === 'remote' ? 'remote' : 
								  job.workArrangement?.toLowerCase() === 'hybrid' ? 'hybrid' : 'on-site',
					minSalary: job.minSalary || 0,
					maxSalary: job.maxSalary || 0,
					currency: job.currency || 'THB',
					type: job.workType ? job.workType.charAt(0).toUpperCase() + job.workType.slice(1) : 'Full-time',
					tags: job.requiredSkills 
						? job.requiredSkills.split(',').map((skill: string) => skill.trim()).slice(0, 3)
						: [],
					postedAt: job.postOpenDate && job.postOpenDate !== "0001-01-01T00:00:00Z"
						? `Posted ${Math.floor((Date.now() - new Date(job.postOpenDate).getTime()) / (1000 * 60 * 60 * 24))} days ago`
						: 'Recently posted',
					badge: (() => {
						// Show "Remote" badge for remote jobs
						if (job.workArrangement?.toLowerCase() === 'remote') {
							return { text: 'Remote', type: 'remote' };
						}
						// Show "Internship" badge for internship/casual positions
						if (job.workType?.toLowerCase().includes('intern') || job.workType?.toLowerCase() === 'casual') {
							return { text: 'Internship', type: 'internship' };
						}
						// Show "New" badge for jobs posted within last 3 days
						if (job.postOpenDate && job.postOpenDate !== "0001-01-01T00:00:00Z") {
							const daysAgo = Math.floor((Date.now() - new Date(job.postOpenDate).getTime()) / (1000 * 60 * 60 * 24));
							if (daysAgo <= 3) {
								return { text: 'New', type: 'new' };
							}
						}
						return null;
					})(),
					logoStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600'
				}));
			}
		} catch (error) {
			console.error('Failed to fetch recent jobs:', error);
			// Keep empty array on error
		} finally {
			isLoadingJobs = false;
		}
	}

	onMount(() => {
		fetchRecentJobs();
	});

	// TODO: Replace with actual API fetch from backend
	const topCompanies = [
		{ id: '1', name: 'Google', positions: 12, logoStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600', logo: 'G' },
		{ id: '2', name: 'Apple', positions: 8, logoStyle: 'bg-gradient-to-br from-gray-800 to-gray-900 text-white', logo: 'A' },
		{ id: '3', name: 'Meta', positions: 15, logoStyle: 'bg-gradient-to-br from-blue-100 to-blue-200 text-blue-600', logo: 'M' },
		{ id: '4', name: 'Amazon', positions: 20, logoStyle: 'bg-gradient-to-br from-orange-100 to-orange-200 text-orange-600', logo: 'A' },
		{ id: '5', name: 'Netflix', positions: 5, logoStyle: 'bg-gradient-to-br from-red-100 to-red-200 text-red-600', logo: 'N' },
		{ id: '6', name: 'Spotify', positions: 7, logoStyle: 'bg-gradient-to-br from-green-100 to-green-200 text-green-600', logo: 'S' },
		{ id: '7', name: 'Microsoft', positions: 18, logoStyle: 'bg-gradient-to-br from-purple-100 to-purple-200 text-purple-600', logo: 'M' },
		{ id: '8', name: 'Tesla', positions: 9, logoStyle: 'bg-gradient-to-br from-red-100 to-red-200 text-red-600', logo: 'T' },
		{ id: '9', name: 'Stripe', positions: 6, logoStyle: 'bg-gradient-to-br from-indigo-100 to-indigo-200 text-indigo-600', logo: 'S' },
		{ id: '10', name: 'Airbnb', positions: 11, logoStyle: 'bg-gradient-to-br from-pink-100 to-pink-200 text-pink-600', logo: 'A' }
	];

	let currentSlide = $state(0);
	let previousSlide = $state(0);
	let screenWidth = $state(0);
	
	// Reactive companies per slide based on screen width
	$effect(() => {
		if (typeof window !== 'undefined') {
			screenWidth = window.innerWidth;
			const updateScreenWidth = () => {
				screenWidth = window.innerWidth;
			};
			window.addEventListener('resize', updateScreenWidth);
			return () => window.removeEventListener('resize', updateScreenWidth);
		}
	});
	
	const companiesPerSlide = $derived(() => {
		if (screenWidth < 640) return 2; // mobile: 2 per slide
		if (screenWidth < 1024) return 3; // tablet: 3 per slide  
		return 5; // desktop: 5 per slide
	});
	
	const totalSlides = $derived(() => Math.ceil(topCompanies.length / companiesPerSlide()));

	function nextSlide() {
		if (currentSlide < totalSlides() - 1) {
			previousSlide = currentSlide;
			currentSlide = currentSlide + 1;
		}
	}

	function prevSlide() {
		if (currentSlide > 0) {
			previousSlide = currentSlide;
			currentSlide = currentSlide - 1;
		}
	}

	function goToSlide(index: number) {
		if (index !== currentSlide) {
			previousSlide = currentSlide;
			currentSlide = index;
		}
	}

	function getVisibleCompanies(slideIndex: number) {
		const start = slideIndex * companiesPerSlide();
		return topCompanies.slice(start, start + companiesPerSlide());
	}
	
	// Reset slide when screen size changes
	$effect(() => {
		if (currentSlide >= totalSlides()) {
			currentSlide = 0;
			previousSlide = 0;
		}
	});

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

<div class="min-h-screen bg-gray-50">
	<!-- Layout Header with transparent background for home page -->
	<LayoutHeader transparent={true} absolute={true} />
	
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
						{#each Array(3) as _}
							<JobCard skeleton={true} />
						{/each}
					{:else if recentJobs.length > 0}
						{#each recentJobs as job}
							<JobCard {job} />
						{/each}
					{:else}
						<div class="col-span-full text-center py-8">
							<div class="max-w-md mx-auto">
								<div class="text-gray-400 mb-4">
									<Briefcase class="w-16 h-16 mx-auto" />
								</div>
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

		<!-- Top Companies -->
		<section class="py-16 px-4 sm:px-6 lg:px-8 bg-slate-50">
			<div class="container mx-auto max-w-7xl">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900">Top Companies</h2>
					<p class="text-gray-600 mt-2">Leading technology companies actively hiring KU students</p>
				</div>

				<!-- Carousel Container -->
				<div class="relative px-4 sm:px-8 lg:px-16">
					<!-- Navigation Buttons positioned outside overflow container -->
					{#if currentSlide > 0}
						<button
							onclick={prevSlide}
							class="absolute left-0 z-10 p-1.5 sm:p-2 bg-white rounded-full shadow-sm cursor-pointer transition-all border border-gray-200 hidden sm:block"
							style="top: 84px; transform: translateY(-50%);"
							aria-label="Previous companies"
						>
							<ChevronLeft class="w-4 h-4 sm:w-5 sm:h-5 text-gray-600" />
						</button>
					{/if}
					
					{#if currentSlide < totalSlides() - 1}
						<button
							onclick={nextSlide}
							class="absolute right-0 z-10 p-1.5 sm:p-2 bg-white rounded-full shadow-sm cursor-pointer transition-all border border-gray-200 hidden sm:block"
							style="top: 84px; transform: translateY(-50%);"
							aria-label="Next companies"
						>
							<ChevronRight class="w-4 h-4 sm:w-5 sm:h-5 text-gray-600" />
						</button>
					{/if}

					<!-- Companies Grid -->
					<div class="overflow-hidden">
						<div 
							class="flex transition-transform duration-500 ease-in-out"
							style="transform: translateX(-{currentSlide * (100 / totalSlides())}%); width: {totalSlides() * 100}%;"
						>
							{#each Array(totalSlides()) as _, slideIndex}
								<div class="flex gap-2 sm:gap-3 lg:gap-4 flex-shrink-0 px-1 sm:px-2" style="width: {100 / totalSlides()}%;">
									{#each getVisibleCompanies(slideIndex) as company (company.id)}
										<div class="flex-1">
											<CompanyCard {company} />
										</div>
									{/each}
								</div>
							{/each}
						</div>
					</div>

					<!-- Pagination Dots -->
					<div class="flex justify-center mt-8 gap-2">
						{#each Array(totalSlides()) as _, index}
							<button
								onclick={() => goToSlide(index)}
								class="w-2 h-2 rounded-full transition-colors {currentSlide === index ? 'bg-green-600' : 'bg-gray-300'} hover:cursor-pointer"
								aria-label="Go to slide {index + 1}"
							></button>
						{/each}
					</div>
				</div>
			</div>
		</section>

		<!-- Career Paths -->
		<section class="py-16 px-4 sm:px-6 lg:px-8 bg-white">
			<div class="container mx-auto max-w-7xl">
				<div class="text-center mb-12">
					<h2 class="text-3xl font-bold text-gray-900">Explore Career Paths</h2>
					<p class="text-gray-600 mt-2 max-w-2xl mx-auto">
						Discover opportunities across different technology domains and find the perfect match for your expertise
					</p>
				</div>

				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
					{#each careerPaths as path}
						<CareerPathCard {path} />
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