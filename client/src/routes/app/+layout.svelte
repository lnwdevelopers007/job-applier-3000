<div class="min-h-screen bg-slate-50 flex flex-col">

  <LayoutHeader />

  <div class="flex flex-1 py-6 px-12 gap-6 items-start">
    <SidebarNav {navItems} />

    <main class="flex-1 p-1 min-w-0">
      <slot />
    </main>
  </div>
</div>

<script>
  import LayoutHeader from '$lib/components/LayoutHeader.svelte';
  import SidebarNav from '$lib/components/SidebarNav.svelte';
  import { House, Cog, Mail, ChartColumn, Briefcase } from 'lucide-svelte';
  
  const navItems = [
    { href: '/app', label: 'Home', icon: House },
    { href: '/app/jobs', label: 'Jobs', icon: Briefcase },
    { href: '/app/applications', label: 'Applications', icon: ChartColumn },
    { href: '/app/settings', label: 'Settings', icon: Cog },
    { href: '/app/requests', label: 'Requests', icon: Mail }
  ];
</script>
