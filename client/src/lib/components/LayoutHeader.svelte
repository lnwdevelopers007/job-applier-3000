<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { jwtDecode } from 'jwt-decode';
  import { LogOut, User } from 'lucide-svelte';
  import { scale } from 'svelte/transition';
  import userAvatar from '$lib/assets/user.png';
  
  type TokenPayload = {
    email: string;
    name: string;
    exp: number;
    role?: string;
    avatarURL?: string;
    userID?: string;
    provider?: string;
  };

  type UserInfo = {
    name: string;
    email: string;
    role: string;
    avatarUrl?: string;
  };

  let userInfo: UserInfo = $state({
    name: 'Guest',
    email: '',
    role: 'User',
    avatarUrl: undefined
  });

  let isLoggedIn = $state(false);
  let isDropdownOpen = $state(false);
  let dropdownRef = $state<HTMLDivElement>();

  onMount(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      try {
        const parsed = JSON.parse(storedUser);
        userInfo = {
          name: parsed.name || 'Guest',
          email: parsed.email || '',
          role: parsed.role || 'User',
          avatarUrl: parsed.avatarUrl
        };
        isLoggedIn = true;
        return;
      } catch (err) {
        console.error('Failed to parse stored user info:', err);
      }
    }

    // Fallback: decode token directly
    const token = localStorage.getItem('access_token');
    if (token) {
      try {
        const decoded = jwtDecode<TokenPayload>(token);
        
        // Check if token is expired
        if (decoded.exp && decoded.exp * 1000 < Date.now()) {
          console.warn('Token has expired');
          localStorage.removeItem('access_token');
          localStorage.removeItem('user');
          isLoggedIn = false;
        } else {
          // Get avatar URL from the JWT (matches MongoDB schema field)
          
          userInfo = {
            name: decoded.name || decoded.email.split('@')[0],
            email: decoded.email,
            role: decoded.role || 'User',
            avatarUrl: decoded.avatarURL // Using avatarURL from MongoDB schema
          };
          
          localStorage.setItem('user', JSON.stringify(userInfo));
          isLoggedIn = true;
        }
      } catch (err) {
        console.error('Failed to decode token:', err);
        isLoggedIn = false;
      }
    } else {
      isLoggedIn = false;
    }
  });

  // Format role display
  function formatRole(role: string): string {
    if (!role) return 'User';
    
    // Capitalize first letter and handle common roles
    const roleMap: Record<string, string> = {
      'student': 'Student',
      'company': 'Recruiter',
      'admin': 'Administrator',
      'recruiter': 'Recruiter'
    };
    
    return roleMap[role.toLowerCase()] || role.charAt(0).toUpperCase() + role.slice(1);
  }

  function toggleDropdown() {
    isDropdownOpen = !isDropdownOpen;
  }

  function handleLogout() {
    // Clear all auth-related data from localStorage
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    
    goto('/');
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

<header class="flex items-center justify-between bg-white border-b border-gray-200 px-6 py-3">
  <div class="flex items-start space-x-1">
      <h2 class="text-lg font-semibold text-black">Job Applier </h2>
      <h2 class="text-lg font-semibold text-green-700">3000</h2>
  </div>

  {#if isLoggedIn}
    <div class="relative" bind:this={dropdownRef}>
      <button
        onclick={toggleDropdown}
        class="flex items-center space-x-3 hover:cursor-pointer focus:outline-gray-300 focus:outline  focus:transition-colors rounded-full"
      >
        <img src={userInfo.avatarUrl || userAvatar} alt="Avatar" class="w-9 h-9 rounded-full object-cover" />
      </button>

      {#if isDropdownOpen}
      <div 
        class="absolute right-0 mt-1 w-64 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50"
        in:scale={{ duration: 100, start: 0.95, opacity: 0 }}
        out:scale={{ duration: 75, start: 0.98, opacity: 0 }}>
        <!-- User info section -->
        <div class="px-4 py-3 border-b border-gray-100 mb-1">
          <div class="flex items-center space-x-3">
            <img src={userInfo.avatarUrl || userAvatar} alt="Avatar" class="w-10 h-10 rounded-full object-cover" />
            <div>
              <p class="text-sm font-semibold text-gray-900">{userInfo.name}</p>
              <p class="text-xs text-gray-500">{userInfo.email}</p>
              <p class="text-xs text-gray-500 mt-0.5">{formatRole(userInfo.role)}</p>
            </div>
          </div>
        </div>

        <!-- Menu items -->
        <div class="px-2">
          <button
            onclick={() => {
              isDropdownOpen = false;
              goto(userInfo.role === 'company' ? '/company/profile' : '/profile');
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