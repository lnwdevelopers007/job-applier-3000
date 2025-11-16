<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { ShieldAlert, Mail, LogOut } from 'lucide-svelte';

  async function handleLogout() {
    try {
      await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include'
      });
    } catch (error) {
      console.error('Logout failed:', error);
    }
    
    goto('/');
  }

  onMount(() => {
    // Clear any pending actions in session storage
    sessionStorage.clear();
  });
</script>

<svelte:head>
  <title>Account Banned - Job Applier 3000</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-red-50 via-orange-50 to-yellow-50 flex items-center justify-center p-4">
  <div class="max-w-md w-full">
    <!-- Card -->
    <div class="bg-white rounded-2xl shadow-xl p-8 text-center">
      <!-- Icon -->
      <div class="flex justify-center mb-6">
        <div class="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center">
          <ShieldAlert class="w-10 h-10 text-red-600" />
        </div>
      </div>

      <!-- Title -->
      <h1 class="text-2xl font-bold text-gray-900 mb-3">
        Account Suspended
      </h1>

      <!-- Message -->
      <p class="text-gray-600 mb-6 leading-relaxed">
        Your account has been suspended and you can no longer access Job Applier 3000 services.
      </p>

      <!-- Info Box -->
      <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 text-left">
        <h2 class="text-sm font-semibold text-red-900 mb-2">
          Why was my account suspended?
        </h2>
        <p class="text-sm text-red-800 mb-3">
          Your account may have been suspended for violating our Terms of Service or Community Guidelines.
        </p>
        <h2 class="text-sm font-semibold text-red-900 mb-2">
          What can I do?
        </h2>
        <p class="text-sm text-red-800">
          If you believe this is a mistake, please contact our support team for assistance.
        </p>
      </div>

      <!-- Actions -->
      <div class="space-y-3">
        <a
          href="mailto:support@jobapplier3000.com?subject=Account%20Ban%20Appeal"
          class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition-colors"
        >
          <Mail class="w-4 h-4" />
          Contact Support
        </a>

        <button
          onclick={handleLogout}
          class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gray-100 text-gray-700 text-sm font-medium rounded-lg hover:bg-gray-200 transition-colors"
        >
          <LogOut class="w-4 h-4" />
          Return to Home
        </button>
      </div>

      <!-- Footer note -->
      <p class="text-xs text-gray-500 mt-6">
        For urgent matters, you can reach us at support@jobapplier3000.com
      </p>
    </div>

    <!-- Additional info -->
    <div class="text-center mt-6">
      <p class="text-sm text-gray-600">
        Need help? Visit our <a href="/" class="text-green-600 hover:text-green-700 font-medium">Help Center</a>
      </p>
    </div>
  </div>
</div>