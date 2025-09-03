<script lang="ts">
  import { page } from '$app/stores';
  import type { Component } from 'svelte';
  
  interface NavItem {
    href: string;
    label: string;
    icon: Component;
    badge?: number;
    separatorBefore?: boolean;
  }
  
  interface Props {
    navItems: NavItem[];
  }
  
  let { navItems }: Props = $props();
</script>

<aside class="w-72 bg-white border border-gray-200 shadow-sm rounded-xl p-3">
  <nav class="space-y-1">
    {#each navItems as item (item.href)}
      {#if item.separatorBefore}
        <div class="border-t border-gray-200 my-4 mx-3"></div>
      {/if}
      <a 
        href={item.href} 
        class="flex items-center justify-between px-4 py-3 rounded-lg transition-all duration-200 {$page.url.pathname === item.href ? 'bg-green-50 text-green-700' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
      >
        <div class="flex items-center gap-3">
          <item.icon class="w-5 h-5 {$page.url.pathname === item.href ? 'text-green-700' : 'text-gray-500'}" />
          <span class="font-medium text-sm">{item.label}</span>
        </div>
        {#if item.badge}
          <span class="px-2.5 py-1 text-xs font-semibold rounded-full bg-green-600 text-white">
            {item.badge}
          </span>
        {/if}
      </a>
    {/each}
  </nav>
</aside>