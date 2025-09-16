<script>
  import { ArrowLeft } from 'lucide-svelte';
  import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
  import OrDivider from '$lib/components/auth/OrDivider.svelte';
  import { fly } from 'svelte/transition';
  
  let currentStep = 1;
  let email = '';
  let password = '';
  let confirmPassword = '';
  let companyName = '';
  let industry = '';
  let size = '';
  let website = '';
  let showPassword = false;
  let showConfirmPassword = false;
  
  function handleEmailStep() {
    if (email) {
      currentStep = 2;
    }
  }
  
  function goBackToEmailStep() {
    currentStep = 1;
  }
  
  function handleSignup() {
    console.log('Company signup:', {
      email,
      password,
      companyName,
      industry,
      size,
      website
    });
  }
</script>

<div class="min-h-screen flex">
  <!-- Left Side - Empty -->
  <div class="hidden lg:flex lg:w-2/3 bg-gradient-to-br from-green-50 to-blue-50">
  </div>

  <!-- Right Side - Signup Form -->
  <div class="w-full lg:w-1/2 flex items-center justify-center p-8 bg-white">
    <div class="w-full max-w-lg">
      <!-- Logo - Always on top -->
      <div class="text-center mb-12">
        <div class="flex items-baseline justify-center gap-1">
          <span class="text-2xl font-semibold text-gray-900">Job Applier</span>
          <span class="text-2xl font-semibold text-green-700">3000</span>
        </div>
      </div>
      {#if currentStep === 1}
        <!-- Step 1: Email -->
        <div class="mb-8" in:fly={{ x: -20, duration: 200 }}>
          <h1 class="text-3xl font-semibold text-gray-900 mb-2">Company Account</h1>
          <p class="text-sm text-gray-500">Enter your company email to get started</p>
        </div>
        
        <form onsubmit={e => { e.preventDefault(); handleEmailStep(); }} class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">Company Email</label>
            <input
              id="email"
              type="email"
              bind:value={email}
              placeholder="Enter company email..."
              required
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
            />
          </div>
          
          <button
            type="submit"
            class="w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-semibold rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
          >
            Continue with Email
          </button>
        </form>

        <OrDivider />
        <GoogleOAuthButton text="Continue with Google" />

        <p class="text-center text-sm text-gray-600 mt-8">
          Already have an account?
          <a href="/login" class="text-green-600 hover:text-green-700 font-medium">Sign in</a>
        </p>
        
      {:else}
        <!-- Step 2: Additional Information -->
        <div class="mb-6" in:fly={{ x: 20, duration: 200 }}>
          <button 
            type="button" 
            onclick={goBackToEmailStep}
            class="flex items-center text-sm text-gray-500 hover:text-gray-700 mb-4 transition-colors"
          >
            <ArrowLeft class="w-4 h-4 mr-1" />
            <span>{email}</span>
          </button>
          <h1 class="text-3xl font-semibold text-gray-900 mb-2">Company Details</h1>
          <p class="text-sm text-gray-500">Tell us more about your company</p>
        </div>
        
        <form onsubmit={e => { e.preventDefault(); handleSignup(); }} class="space-y-4">
          <div>
            <label for="companyName" class="block text-sm font-medium text-gray-700 mb-2">Company Name</label>
            <input
              id="companyName"
              type="text"
              bind:value={companyName}
              placeholder="Enter company name..."
              required
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
            />
          </div>
          
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="industry" class="block text-sm font-medium text-gray-700 mb-2">Industry</label>
              <select
                id="industry"
                bind:value={industry}
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
              >
                <option value="" disabled selected>Select industry</option>
                <option value="technology">Technology</option>
                <option value="finance">Finance</option>
                <option value="healthcare">Healthcare</option>
                <option value="education">Education</option>
                <option value="retail">Retail</option>
                <option value="manufacturing">Manufacturing</option>
                <option value="other">Other</option>
              </select>
            </div>
            <div>
              <label for="companySize" class="block text-sm font-medium text-gray-700 mb-2">Company Size</label>
              <select
                id="companySize"
                bind:value={size}
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
              >
                <option value="" disabled selected>Select size</option>
                <option value="1-10">1-10 employees</option>
                <option value="11-50">11-50 employees</option>
                <option value="51-200">51-200 employees</option>
                <option value="201-500">201-500 employees</option>
                <option value="500+">500+ employees</option>
              </select>
            </div>
          </div>
          
          <div>
            <label for="website" class="block text-sm font-medium text-gray-700 mb-2">Website (Optional)</label>
            <input
              id="website"
              type="url"
              bind:value={website}
              placeholder="Enter website URL..."
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
            />
          </div>
          
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <div class="relative">
              <input
                id="password"
                type={showPassword ? 'text' : 'password'}
                bind:value={password}
                placeholder="Enter password..."
                required
                class="w-full px-4 py-3 pr-12 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
              />
              <button
                type="button"
                onclick={() => showPassword = !showPassword}
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
                aria-label={showPassword ? 'Hide password' : 'Show password'}
              >
                {#if showPassword}
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                  </svg>
                {:else}
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                {/if}
              </button>
            </div>
          </div>
          
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">Confirm Password</label>
            <div class="relative">
              <input
                id="confirmPassword"
                type={showConfirmPassword ? 'text' : 'password'}
                bind:value={confirmPassword}
                placeholder="Enter password again..."
                required
                class="w-full px-4 py-3 pr-12 text-sm border border-gray-300 rounded-lg focus:ring-1 focus:ring-gray-500 focus:border-gray-400 outline-none transition-colors bg-gray-50"
              />
              <button
                type="button"
                onclick={() => showConfirmPassword = !showConfirmPassword}
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
                aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
              >
                {#if showConfirmPassword}
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
                  </svg>
                {:else}
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                {/if}
              </button>
            </div>
          </div>
          
          <button
            type="submit"
            class="w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-semibold rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 mt-6"
          >
            Create Account
          </button>
        </form>
      {/if}
    </div>
  </div>
</div>