<script>
  import AuthContainer from '$lib/components/auth/AuthContainer.svelte';
  import AuthHeader from '$lib/components/auth/AuthHeader.svelte';
  import PasswordInput from '$lib/components/auth/PasswordInput.svelte';
  import GoogleOAuthButton from '$lib/components/auth/GoogleOAuthButton.svelte';
  import OrDivider from '$lib/components/auth/OrDivider.svelte';
  
  let email = '';
  let password = '';
  let rememberMe = false;
  let showPassword = false;

  function handleLogin() {
    console.log('Login:', { email, password, rememberMe });
  }
  
  function handleGoogleOAuth() {
    console.log('Google OAuth login');
  }
</script>

<AuthContainer>
  <AuthHeader />
  <h1 class="text-2xl font-semibold text-gray-900 mb-5">Sign in</h1>

  <form on:submit|preventDefault={handleLogin} class="space-y-6">
    <div>
      <input
        id="email"
        type="email"
        bind:value={email}
        placeholder="Enter email..."
        required
        class="w-full px-4 py-3 text-sm border border-gray-300 rounded-md focus:ring-1 focus:ring-gray-500 outline-none transition-colors"
      />
    </div>

    <PasswordInput
      id="password"
      label=""
      placeholder="Enter password..."
      bind:value={password}
      bind:showPassword
    />

    <div class="flex items-center justify-between">
      <label class="flex items-center cursor-pointer">
        <input
          type="checkbox"
          bind:checked={rememberMe}
          class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-0 focus:ring-offset-0"
        />
        <span class="ml-2 text-sm text-gray-700">Remember me</span>
      </label>
      <a href="/auth/forgot-password" class="text-sm text-green-600 hover:text-green-700 font-medium">
        Forgot password
      </a>
    </div>

    <button
      type="submit"
      class="w-full py-2 px-4 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-md transition-colors focus:outline-none focus:ring-1 focus:ring-green-500"
    >
      Log in
    </button>
  </form>

  <OrDivider />
  <GoogleOAuthButton onClick={handleGoogleOAuth} />
</AuthContainer>