<script>
  import { ArrowLeft } from 'lucide-svelte';
  import AuthContainer from '$lib/components/auth/AuthContainer.svelte';
  import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
  import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
  import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
  import OrDivider from '$lib/components/auth/OrDivider.svelte';
  
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
  
  function handleGoogleOAuth() {
    console.log('Google OAuth for company');
  }
</script>

<AuthContainer maxWidth={currentStep === 2 ? 'max-w-2xl' : 'max-w-lg'}>
  <AuthHeader onBack={() => currentStep === 1 ? window.history.back() : goBackToEmailStep()} />
  
  {#if currentStep === 1}
    <!-- Step 1: Email -->
    <h1 class="text-2xl font-semibold text-gray-900 mb-2">Register Company</h1>
    <p class="text-sm text-gray-500 mb-6">Start hiring talented CPE students and alumni</p>
    
    <form on:submit|preventDefault={handleEmailStep} class="space-y-4">
      <div>
        <input
          type="email"
          bind:value={email}
          placeholder="Enter your company email"
          required
          class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
        />
      </div>
      
      <button
        type="submit"
        class="w-full py-2.5 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
      >
        Continue
      </button>
    </form>

    <OrDivider />
    <GoogleOAuthButton onClick={handleGoogleOAuth} />

    <p class="text-center text-sm text-gray-600 mt-6">
      Already have an account? 
      <a href="/auth/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
    </p>
    
  {:else}
    <!-- Step 2: Additional Information -->
    <h1 class="text-2xl font-semibold text-gray-900 mb-2">Complete Your Profile</h1>
    <p class="text-sm text-gray-500 mb-6">Tell us more about your company</p>
    
    <div class="flex items-center space-x-3 mb-6">
      <ArrowLeft class="w-4 h-4 text-gray-600 cursor-pointer" on:click={goBackToEmailStep} />
      <p class="text-sm font-light text-gray-600"><span class="font-medium">{email}</span></p>
    </div>
    
    <form on:submit|preventDefault={handleSignup} class="space-y-4">
      <div>
        <label for="companyName" class="block text-sm font-medium text-gray-700 mb-2">Company Name</label>
        <input
          id="companyName"
          type="text"
          bind:value={companyName}
          placeholder="Enter company name"
          required
          class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
        />
      </div>
      
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label for="industry" class="block text-sm font-medium text-gray-700 mb-2">Industry</label>
          <select
            id="industry"
            bind:value={industry}
            required
            class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
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
            class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
          >
            <option value="" disabled selected>Select company size</option>
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
          placeholder="https://yourcompany.com"
          class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
        />
      </div>
      
      <PasswordInput
        id="password"
        label="Password"
        placeholder="Create password"
        bind:value={password}
        bind:showPassword
      />
      
      <PasswordInput
        id="confirmPassword"
        label="Confirm Password"
        placeholder="Confirm password"
        bind:value={confirmPassword}
        bind:showPassword={showConfirmPassword}
      />
      
      <button
        type="submit"
        class="w-full py-3 px-4 bg-green-600 hover:bg-green-700 text-white mt-5 text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
      >
        Sign Up
      </button>
    </form>
  {/if}
</AuthContainer>