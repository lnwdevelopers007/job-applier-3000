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

<AuthLayout backHref="/signup">
  <AuthHeader />
  
  {#if currentStep === 1}
    <div class="mb-8" in:fly={{ x: -20, duration: 200 }}>
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Sign up</h1>
      <p class="text-sm text-gray-500">Enter your email to start hiring talents</p>
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
      <h1 class="text-3xl font-semibold text-gray-900 mb-2">Company Details</h1>
      <p class="text-sm text-gray-500">Tell us more about your company</p>
    </div>
    
    <form onsubmit={e => { e.preventDefault(); handleSignup(); }} class="space-y-4">
      <FormInput
        id="companyName"
        label="Company Name"
        placeholder="Enter company name..."
        bind:value={companyName}
        required
      />
      
      <div class="grid grid-cols-2 gap-4">
        <FormSelect
          id="industry"
          label="Industry"
          bind:value={industry}
          placeholder="Select industry"
          required
        >
          <option value="technology">Technology</option>
          <option value="finance">Finance</option>
          <option value="healthcare">Healthcare</option>
          <option value="education">Education</option>
          <option value="retail">Retail</option>
          <option value="manufacturing">Manufacturing</option>
          <option value="other">Other</option>
        </FormSelect>
        
        <FormSelect
          id="companySize"
          label="Company Size"
          bind:value={size}
          placeholder="Select size"
          required
        >
          <option value="1-10">1-10 employees</option>
          <option value="11-50">11-50 employees</option>
          <option value="51-200">51-200 employees</option>
          <option value="201-500">201-500 employees</option>
          <option value="500+">500+ employees</option>
        </FormSelect>
      </div>
      
      <FormInput
        id="website"
        type="url"
        label="Website (Optional)"
        placeholder="Enter website URL..."
        bind:value={website}
      />
      
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