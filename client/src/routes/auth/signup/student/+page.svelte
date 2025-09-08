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
  let firstName = '';
  let lastName = '';
  let studentId = '';
  let year = '';
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
    console.log('Student signup:', {
      email,
      password,
      firstName,
      lastName,
      studentId,
      year
    });
  }
  
  function handleGoogleOAuth() {
    console.log('Google OAuth for student');
  }
</script>

<AuthContainer maxWidth={currentStep === 2 ? 'max-w-2xl' : 'max-w-lg'}>
  <AuthHeader onBack={() => currentStep === 1 ? window.history.back() : goBackToEmailStep()} />
  
  {#if currentStep === 1}
    <!-- Step 1: Email -->
    <h1 class="text-2xl font-semibold text-gray-900 mb-2">Create Student Account</h1>
    <p class="text-sm text-gray-500 mb-6">Join to find your dream job opportunities</p>
    
    <form on:submit|preventDefault={handleEmailStep} class="space-y-4">
      <div>
        <input
          type="email"
          bind:value={email}
          placeholder="Enter your university email"
          required
          class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
        />
      </div>
      
      <button
        type="submit"
        class="w-full py-2 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
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
    <p class="text-sm text-gray-500 mb-6">Tell us more about yourself</p>
    
    <div class="flex items-center space-x-3 mb-6">
      <ArrowLeft class="w-4 h-4 text-gray-600 cursor-pointer" on:click={goBackToEmailStep} />
      <p class="text-sm font-light text-gray-600"><span class="font-medium">{email}</span></p>
    </div>
    
    <form on:submit|preventDefault={handleSignup} class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label for="firstName" class="block text-sm font-medium text-gray-700 mb-2">First Name</label>
          <input
            id="firstName"
            type="text"
            bind:value={firstName}
            placeholder="Enter first name"
            required
            class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
          />
        </div>
        <div>
          <label for="lastName" class="block text-sm font-medium text-gray-700 mb-2">Last Name</label>
          <input
            id="lastName"
            type="text"
            bind:value={lastName}
            placeholder="Enter last name"
            required
            class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
          />
        </div>
      </div>
      
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label for="studentId" class="block text-sm font-medium text-gray-700 mb-2">Student ID</label>
          <input
            id="studentId"
            type="text"
            bind:value={studentId}
            placeholder="Enter student ID"
            required
            class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
          />
        </div>
        <div>
          <label for="year" class="block text-sm font-medium text-gray-700 mb-2">Year</label>
          <select
            id="year"
            bind:value={year}
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