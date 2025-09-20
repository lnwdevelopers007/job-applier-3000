<script>
  import { ArrowLeft } from 'lucide-svelte';
  import { fly } from 'svelte/transition';
  import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
  import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
  import FormInput from '$lib/components/auth/FormInput.svelte';
  import FormSelect from '$lib/components/auth/FormSelect.svelte';
  import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
  import FormButton from '$lib/components/auth/FormButton.svelte';
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
</script>

<AuthLayout backHref="/signup">
  <AuthHeader />
  
  {#if currentStep === 1}
    <div class="mb-8" in:fly={{ x: -20, duration: 200 }}>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Sign up</h1>
      <p class="text-sm text-gray-500">Enter your email to start finding your dream job</p>
    </div>
    
    <form onsubmit={e => { e.preventDefault(); handleEmailStep(); }} class="space-y-6">
      <FormInput
        id="email"
        type="email"
        label="Email"
        placeholder="Enter email..."
        bind:value={email}
        required
      />
      
      <FormButton type="submit">Continue with Email</FormButton>
    </form>

    <OrDivider />
    <GoogleOAuthButton text="Continue with Google" />

    <p class="text-center text-sm text-gray-500 mt-8">
      Already have an account?
      <a href="/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
    </p>
    
  {:else}
    <div class="mb-6" in:fly={{ x: 20, duration: 200 }}>
      <button 
        type="button" 
        onclick={goBackToEmailStep}
        class="flex items-center text-sm text-gray-500 hover:text-gray-700 mb-4 transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 mr-1" />
        <span>{email}</span>
      </button>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Complete Profile</h1>
      <p class="text-sm text-gray-500">Tell us more about yourself</p>
    </div>
    
    <form onsubmit={e => { e.preventDefault(); handleSignup(); }} class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <FormInput
          id="firstName"
          label="First Name"
          placeholder="Enter first name..."
          bind:value={firstName}
          required
        />
        <FormInput
          id="lastName"
          label="Last Name"
          placeholder="Enter last name..."
          bind:value={lastName}
          required
        />
      </div>
      
      <div class="grid grid-cols-2 gap-4">
        <FormInput
          id="studentId"
          label="Student ID"
          placeholder="Enter student ID..."
          bind:value={studentId}
          required
        />
        <FormSelect
          id="year"
          label="Year"
          bind:value={year}
          placeholder="Select year"
          required
        >
          <option value="1">Year 1</option>
          <option value="2">Year 2</option>
          <option value="3">Year 3</option>
          <option value="4">Year 4</option>
          <option value="alumni">Alumni</option>
        </FormSelect>
      </div>
      
      <PasswordInput
        id="password"
        label="Password"
        placeholder="Enter password..."
        bind:value={password}
        bind:showPassword
        required
      />
      
      <PasswordInput
        id="confirmPassword"
        label="Confirm Password"
        placeholder="Enter password again..."
        bind:value={confirmPassword}
        bind:showPassword={showConfirmPassword}
        required
      />
      
      <FormButton type="submit">Create Account</FormButton>
    </form>
  {/if}
</AuthLayout>