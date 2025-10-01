<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { LogOut, User } from 'lucide-svelte';
  import { scale } from 'svelte/transition';
  import userAvatar from '$lib/assets/user.png';
  import { isAuthenticated, getUserInfo, logout, type UserInfo } from '$lib/utils/auth';
  
  let { transparent = false, absolute = false }: { transparent?: boolean; absolute?: boolean } = $props();

  let userInfo: UserInfo | null = $state(null);
  let isLoggedIn = $state(false);
  let isDropdownOpen = $state(false);
  let dropdownRef = $state<HTMLDivElement>();

  onMount(() => {
    isLoggedIn = isAuthenticated();
    if (isLoggedIn) {
      userInfo = getUserInfo();
    }
  });

  // Format role display
  function formatRole(role: string): string {
    if (!role) return 'User';
    
    const roleMap: Record<string, string> = {
      'student': 'Student',
      'company': 'Recruiter',
      'admin': 'Administrator',
    };
    
    return roleMap[role.toLowerCase()] || role.charAt(0).toUpperCase() + role.slice(1);
  }

  function toggleDropdown() {
    isDropdownOpen = !isDropdownOpen;
  }

  function handleLogout() {
    logout();
    
    if (window.location.pathname === '/') {
      window.location.reload();
    } else {
      goto('/');
    }
  }

  // Close dropdown when clicking outside
  function handleClickOutside(event: MouseEvent) {
    if (dropdownRef && !dropdownRef.contains(event.target as Node)) {
      isDropdownOpen = false;
    }
  }

  $effect(() => {
    if (isDropdownOpen) {
      document.addEventListener('click', handleClickOutside);
      return () => {
        document.removeEventListener('click', handleClickOutside);
      };
    }
  });
</script>

<header class="{absolute ? 'absolute top-0 left-0 right-0 z-20' : ''} flex items-center justify-between {transparent ? '' : 'bg-white border-b border-gray-200'} px-6 py-3">
  <button 
    class="flex items-start space-x-1 group cursor-pointer"
    onclick={() => goto('/')}
    onkeydown={(e) => e.key === 'Enter' && goto('/')}
    aria-label="Go to home page"
  >
    <h2 class="text-lg font-semibold text-black">Job Applier </h2>
    <h2 class="text-lg font-semibold text-green-700">3000</h2>
  </button>

  {#if isLoggedIn}
    <div class="relative" bind:this={dropdownRef}>
      <button
        onclick={toggleDropdown}
        class="flex items-center space-x-3 hover:cursor-pointer focus:outline-gray-300 focus:outline  focus:transition-colors rounded-full"
      >
        <img src={userInfo?.avatarUrl || userAvatar} alt="Avatar" class="w-9 h-9 rounded-full object-cover" />
      </button>

      {#if isDropdownOpen}
      <div 
        class="absolute right-0 mt-1 w-64 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50"
        in:scale={{ duration: 100, start: 0.95, opacity: 0 }}
        out:scale={{ duration: 75, start: 0.98, opacity: 0 }}>
        <!-- User info section -->
        <div class="px-4 py-3 border-b border-gray-100 mb-1">
          <div class="flex items-center space-x-3">
            <img src={userInfo?.avatarUrl || userAvatar} alt="Avatar" class="w-10 h-10 rounded-full object-cover" />
            <div>
              <p class="text-sm font-semibold text-gray-900">{userInfo?.name || 'Guest'}</p>
              <p class="text-xs text-gray-500">{userInfo?.email || ''}</p>
              <p class="text-xs text-gray-500 mt-0.5">{formatRole(userInfo?.role || 'User')}</p>
            </div>
          </div>
        </div>

        <!-- Menu items -->
        <div class="px-2">
          <button
            onclick={() => {
              isDropdownOpen = false;
              goto(userInfo?.role === 'company' ? '/company/profile' : '/profile');
            }}
            class="w-full flex items-center space-x-3 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 hover:cursor-pointer transition-colors"
          >
            <User class="w-4 h-4" />
            <span>Profile</span>
          </button>

          <hr class="my-1 border-gray-100" />

          <button
            onclick={handleLogout}
            class="w-full flex items-center space-x-3 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50 hover:cursor-pointer transition-colors"
          >
            <LogOut class="w-4 h-4" />
            <span>Sign out</span>
          </button>
        </div>
      </div>
      {/if}
    </div>
  {:else}
    <div class="flex items-center space-x-3">
      <a
        href="/login"
        class="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 transition-colors"
      >
        Log in
      </a>
      <a
        href="/signup"
        class="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
      >
        Sign up
      </a>
    </div>
  {/if}
</header>