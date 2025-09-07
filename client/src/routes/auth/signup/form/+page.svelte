<script>
  import { page } from '$app/stores';
  import { ChevronLeft, Eye, EyeOff, ArrowLeft, } from 'lucide-svelte';
  
  $: role = $page.url.searchParams.get('role');
  $: isStudent = role === 'student';
  $: isCompany = role === 'company';
  
  // Step management
  let currentStep = 2; // 1: email, 2: info
  
  // Common data
  let email = 'iampattapon@gmail.com';
  let password = '';
  let confirmPassword = '';
  let agreeToTerms = false;
  
  // Student form data
  let studentFirstName = '';
  let studentLastName = '';
  let studentId = '';
  let studentYear = '';
  
  // Company form data
  let companyName = '';
  let companyWebsite = '';
  let companyIndustry = '';
  let companySize = '';
  
  let showPassword = false;
  let showConfirmPassword = false;
  
  function handleEmailStep() {
    if (email) {
      currentStep = 2;
    }
  }
  
  function handleStudentSignup() {
    console.log('Student signup:', {
      email,
      password,
      firstName: studentFirstName,
      lastName: studentLastName,
      studentId,
      year: studentYear
    });
  }
  
  function handleCompanySignup() {
    console.log('Company signup:', {
      name: companyName,
      email,
      password,
      website: companyWebsite,
      industry: companyIndustry,
      size: companySize
    });
  }
  
  function goBackToEmailStep() {
    currentStep = 1;
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-white py-12">
  {#if isStudent}
    <!-- Student Signup Form -->
    <div class="w-full max-w-lg p-10 bg-white rounded-sm border border-gray-200 shadow-sm">
      <!-- Header -->
      <div class="flex items-center space-x-1 mb-6">
        <ChevronLeft 
          class="w-5 h-5 text-gray-600 cursor-pointer" 
          on:click={() => currentStep === 1 ? window.history.back() : goBackToEmailStep()} 
        />
        <h2 class="text-lg font-semibold text-black">Job Applier </h2>
        <h2 class="text-lg font-semibold text-green-700">3000</h2>
      </div>
      
      {#if currentStep === 1}
        <!-- Step 1: Email -->
        <h1 class="text-2xl font-semibold text-gray-900 mb-2">Create Student Account</h1>
        <p class="text-sm text-gray-500 mb-6">Join to find your dream job opportunities</p>
        
        <form on:submit|preventDefault={handleEmailStep} class="space-y-4">
          <!-- Email -->
          <div>
            <input
              type="email"
              bind:value={email}
              placeholder="Enter your university email"
              required
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
          </div>
          
          <!-- Continue button -->
          <button
            type="submit"
            class="w-full py-2 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
          >
            Continue
          </button>
        </form>

        <!-- Divider -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-300"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-white text-gray-500">or</span>
          </div>
        </div>

        <!-- Google OAuth button -->
        <button
          type="button"
          class="w-full flex items-center justify-center gap-3 py-2 px-4 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors focus:outline-none focus:ring-1 focus:ring-gray-500"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          <span class="text-sm font-medium text-gray-700">Continue with Google</span>
        </button>

        <!-- Login link -->
        <p class="text-center text-sm text-gray-600 mt-6">
          Already have an account? 
          <a href="/auth/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
        </p>
        
      {:else}
        <!-- Step 2: Additional Information -->
        <h1 class="text-2xl font-semibold text-gray-900 mb-2">Complete Your Profile</h1>
        <p class="text-sm text-gray-500 mb-6">Tell us more about yourself</p>
        
        <!-- Email display -->
        <div class="mb-4 p-3 bg-gray-50 rounded-md">
          <p class="text-sm text-gray-600">Email: <span class="font-medium">{email}</span></p>
        </div>
        
        <form on:submit|preventDefault={handleStudentSignup} class="space-y-4">
          <!-- Name fields row -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <input
                type="text"
                bind:value={studentFirstName}
                placeholder="First name"
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              />
            </div>
            <div>
              <input
                type="text"
                bind:value={studentLastName}
                placeholder="Last name"
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              />
            </div>
          </div>
          
          <!-- Student ID and Year -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <input
                type="text"
                bind:value={studentId}
                placeholder="Student ID"
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              />
            </div>
            <div>
              <select
                bind:value={studentYear}
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              >
                <option value="" disabled selected>Select year</option>
                <option value="1">Year 1</option>
                <option value="2">Year 2</option>
                <option value="3">Year 3</option>
                <option value="4">Year 4</option>
                <option value="alumni">Alumni</option>
              </select>
            </div>
          </div>
          
          <!-- Password -->
          <div class="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              bind:value={password}
              placeholder="Create password"
              required
              class="w-full px-4 py-3 text-sm pr-12 border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
            <button
              type="button"
              on:click={() => showPassword = !showPassword}
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
            >
              {#if showPassword}
                <EyeOff class="w-4 h-4" />
              {:else}
                <Eye class="w-4 h-4" />
              {/if}
            </button>
          </div>
          
          <!-- Confirm Password -->
          <div class="relative">
            <input
              type={showConfirmPassword ? 'text' : 'password'}
              bind:value={confirmPassword}
              placeholder="Confirm password"
              required
              class="w-full px-4 py-3 text-sm pr-12 border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
            <button
              type="button"
              on:click={() => showConfirmPassword = !showConfirmPassword}
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
              aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
            >
              {#if showConfirmPassword}
                <EyeOff class="w-4 h-4" />
              {:else}
                <Eye class="w-4 h-4" />
              {/if}
            </button>
          </div>
          
          <!-- Terms checkbox -->
          <label class="flex items-start cursor-pointer">
            <input
              type="checkbox"
              bind:checked={agreeToTerms}
              required
              class="w-4 h-4 mt-0.5 text-green-600 border-gray-300 rounded focus:ring-0 focus:ring-offset-0"
            />
            <span class="ml-2 text-sm text-gray-700">
              I agree to the <a href="/terms" class="text-green-600 hover:text-green-700">Terms and Conditions</a> and 
              <a href="/privacy" class="text-green-600 hover:text-green-700">Privacy Policy</a>
            </span>
          </label>
          
          <!-- Submit button -->
          <button
            type="submit"
            class="w-full py-2 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
          >
            Create Account
          </button>
        </form>
      {/if}
    </div>
    
  {:else if isCompany}
    <!-- Company Signup Form -->
    <div class="w-full max-w-lg p-10 bg-white rounded-sm">
      <!-- Header -->
      <div class="flex items-center space-x-1 mb-6">
        <ChevronLeft 
          class="w-5 h-5 text-gray-600 cursor-pointer" 
          on:click={() => currentStep === 1 ? window.history.back() : goBackToEmailStep()} 
        />
        <h2 class="text-lg font-semibold text-black">Job Applier </h2>
        <h2 class="text-lg font-semibold text-green-700">3000</h2>
      </div>
      
      {#if currentStep === 1}
        <!-- Step 1: Email -->
        <h1 class="text-2xl font-semibold text-gray-900 mb-2">Register Company</h1>
        <p class="text-sm text-gray-500 mb-6">Start hiring talented CPE students and alumni</p>
        
        <form on:submit|preventDefault={handleEmailStep} class="space-y-4">
          <!-- Email -->
          <div>
            <input
              type="email"
              bind:value={email}
              placeholder="Enter your company email"
              required
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
          </div>
          
          <!-- Continue button -->
          <button
            type="submit"
            class="w-full py-2.5 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
          >
            Continue
          </button>
        </form>

        <!-- Divider -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-300"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-white text-xs text-gray-500">or</span>
          </div>
        </div>

        <!-- Google OAuth button -->
        <button
          type="button"
          class="w-full flex items-center justify-center gap-3 py-2.5 px-4 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors focus:outline-none focus:ring-1 focus:ring-gray-500"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          <span class="text-sm font-medium text-gray-700">Continue with Google</span>
        </button>

        <!-- Login link -->
        <p class="text-center text-sm text-gray-600 mt-6">
          Already have an account? 
          <a href="/auth/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
        </p>
        
      {:else}
        <!-- Step 2: Additional Information -->
        <h1 class="text-2xl font-semibold text-gray-900 mb-2">Complete Your Profile</h1>
        <p class="text-sm text-gray-500 mb-6">Tell us more about your company</p>
        
        <!-- Email display -->
        <div class="flex items-center space-x-3 mb-6">
          <ArrowLeft class="w-4 h-4 text-gray-600 cursor-pointer" on:click={goBackToEmailStep} />
          <p class="text-sm font-light text-gray-600"><span class="font-medium">{email}</span></p>
        </div>
        
        <form on:submit|preventDefault={handleCompanySignup} class="space-y-4">
          <!-- Company Name -->
          <div>
            <input
              type="text"
              bind:value={companyName}
              placeholder="Company name"
              required
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
          </div>
          
          <!-- Industry and Size -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <select
                bind:value={companyIndustry}
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              >
                <option value="" disabled selected>Industry</option>
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
              <select
                bind:value={companySize}
                required
                class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
              >
                <option value="" disabled selected>Company size</option>
                <option value="1-10">1-10 employees</option>
                <option value="11-50">11-50 employees</option>
                <option value="51-200">51-200 employees</option>
                <option value="201-500">201-500 employees</option>
                <option value="500+">500+ employees</option>
              </select>
            </div>
          </div>
          
          <!-- Company Website -->
          <div>
            <input
              type="url"
              bind:value={companyWebsite}
              placeholder="Company website (optional)"
              class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
          </div>
          
          <!-- Password -->
          <div class="relative">
            <input
              type={showPassword ? 'text' : 'password'}
              bind:value={password}
              placeholder="Create password"
              required
              class="w-full px-4 py-3 text-sm pr-12 border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
            <button
              type="button"
              on:click={() => showPassword = !showPassword}
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
            >
              {#if showPassword}
                <EyeOff class="w-4 h-4" />
              {:else}
                <Eye class="w-4 h-4" />
              {/if}
            </button>
          </div>
          
          <!-- Confirm Password -->
          <div class="relative">
            <input
              type={showConfirmPassword ? 'text' : 'password'}
              bind:value={confirmPassword}
              placeholder="Confirm password"
              required
              class="w-full px-4 py-3 text-sm pr-12 border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
            />
            <button
              type="button"
              on:click={() => showConfirmPassword = !showConfirmPassword}
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
              aria-label={showConfirmPassword ? 'Hide password' : 'Show password'}
            >
              {#if showConfirmPassword}
                <EyeOff class="w-4 h-4" />
              {:else}
                <Eye class="w-4 h-4" />
              {/if}
            </button>
          </div>
          
          <!-- Terms checkbox -->
          <label class="flex items-start cursor-pointer">
            <input
              type="checkbox"
              bind:checked={agreeToTerms}
              required
              class="w-4 h-4 mt-0.5 text-green-600 border-gray-300 rounded focus:ring-0 focus:ring-offset-0"
            />
            <span class="ml-2 text-sm text-gray-700">
              I agree to the <a href="/terms" class="text-green-600 hover:text-green-700">Terms and Conditions</a> and 
              <a href="/privacy" class="text-green-600 hover:text-green-700">Privacy Policy</a>
            </span>
          </label>
          
          <!-- Submit button -->
          <button
            type="submit"
            class="w-full py-2 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
          >
            Create Company Account
          </button>
        </form>
      {/if}
    </div>
    
  {:else}
    <!-- No role selected -->
    <div class="w-full max-w-lg p-10 bg-white rounded-sm border border-gray-200 shadow-sm">
      <h1 class="text-2xl font-semibold text-gray-900 mb-4">Please select a role</h1>
      <a href="/auth/signup" class="text-green-600 hover:text-green-700 font-medium">Go back to role selection</a>
    </div>
  {/if}
</div>