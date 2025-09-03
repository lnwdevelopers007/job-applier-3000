<div class="min-h-screen bg-slate-50 flex flex-col mx-auto">

  <LayoutHeader userName="Company name" userRole="Recruiter" />

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
  import { Cog, Building, ChartColumn, ChartLine, Users, CirclePlus } from 'lucide-svelte';
  
  const navItems = [
    { href: '/company/dashboard', label: 'Dashboard', icon: ChartColumn },
    { href: '/company/applicants', label: 'Applicants', icon: Users, badge: 67 },
    { href: '/company/analytics', label: 'Analytics', icon: ChartLine },
    { href: '/company/create', label: 'Post New Job', icon: CirclePlus, separatorBefore: true },
    { href: '/company/profile', label: 'Company Profile', icon: Building },
    { href: '/company/settings', label: 'Settings', icon: Cog }
  ];
</script>
