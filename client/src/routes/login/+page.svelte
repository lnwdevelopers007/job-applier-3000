<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { isAuthenticated } from '$lib/utils/auth';
  import AuthLayout from '$lib/components/auth/AuthLayout.svelte';
  import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
  import FormInput from '$lib/components/auth/FormInput.svelte';
  import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
  import FormButton from '$lib/components/auth/FormButton.svelte';
  import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
  import OrDivider from '$lib/components/auth/OrDivider.svelte';
  
  let email = '';
  let password = '';
  let rememberMe = false;
  let showPassword = false;

  onMount(() => {
    if (isAuthenticated()) {
      goto('/app/jobs');
    }
  });

  function handleLogin() {
    console.log('Login:', { email, password, rememberMe });
  }

</script>

<AuthLayout>
  <AuthHeader showLogo={true} />
  
  <div class="mb-8">
    <h1 class="text-3xl font-semibold text-gray-900 mb-2">Log in</h1>
    <p class="text-sm text-gray-500">Welcome back! Please login to continue.</p>
  </div>

  <form onsubmit={e => { e.preventDefault(); handleLogin(); }} class="space-y-6">
    <FormInput
      id="email"
      type="email"
      label="Email"
      placeholder="Enter email..."
      bind:value={email}
      required
    />

    <PasswordInput
      id="password"
      label="Password"
      placeholder="Enter password..."
      bind:value={password}
      bind:showPassword
      required
    />

    <FormButton type="submit">Login</FormButton>
  </form>

  <OrDivider />
  
  <GoogleOAuthButton text="Continue with Google" userType="login" />
  
  <p class="text-center text-sm text-gray-600 mt-8">
    Don't have an account?
    <a href="/signup" class="text-green-600 hover:text-green-700 font-medium">
      Sign up now
    </a>
  </p>
</AuthLayout>