<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { writable } from 'svelte/store';
  import { Chart, registerables } from 'chart.js';
  import { getUserInfo, isAuthenticated } from '$lib/utils/auth';
  import { goto } from '$app/navigation';

  Chart.register(...registerables);

  const viewOptions = [
      { label: 'Applicants', value: 'applicants' },
      { label: 'Accepted & Rejected', value: 'accepted_rejected' },
      { label: 'Job Postings', value: 'job_postings' },
      { label: 'Applied Jobs', value: 'applied_jobs' }
  ];
  const selectedView = writable('applicants');

  let company: any = null;
  let chart: Chart | null = null;
  let canvas: HTMLCanvasElement;

  function getChartTitle(view: string): string {
    switch (view) {
      case 'applicants': return 'Monthly Applicants Overview';
      case 'accepted_rejected': return 'Accepted vs Rejected Applicants (Last 12 Months)';
      case 'job_postings': return 'Ongoing Job Postings Over Time';
      case 'applied_jobs': return 'Top 10 Most Applied Jobs';
      default: return '';
    }
  }

  function getAxisLabels(view: string) {
    switch (view) {
      case 'applicants':
        return { x: 'Month', y: 'Number of Applicants' };
      case 'accepted_rejected':
        return { x: 'Month', y: 'Number of Applicants' };
      case 'job_postings':
        return { x: 'Month', y: 'Ongoing Job Postings' };
      case 'applied_jobs':
        return { x: 'Number of Applicants', y: 'Job title' };
      default:
        return { x: '', y: '' };
    }
  }

  function getLast12Months(): string[] {
    const now = new Date();
    const months: string[] = [];
    for (let i = 11; i >= 0; i--) {
      const dt = new Date(now.getFullYear(), now.getMonth() - i, 1);
      months.push(`${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}`);
    }
    return months;
  }

  function getNestedValue(obj: any, path: string): any {
    return path.split('.').reduce((acc, key) => acc && acc[key], obj);
  }

  function groupByMonth(data: any[], dateKey: string) {
    const result: Record<string, number> = {};
    for (const d of data) {
      const rawDate = getNestedValue(d, dateKey);
      // Not include corrupted records
      if (!rawDate) continue;

      const dt = new Date(rawDate);
      if (isNaN(dt.getTime())) continue;

      const key = `${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}`;
      result[key] = (result[key] || 0) + 1;
    }
    console.log('Grouped months:', result);
    return result;
  }

  async function fetchData(status: string = '', mode: 'applications' | 'jobs' = 'applications') {
    if (!isAuthenticated()) {
      goto('/login');
      return [];
    }
    if (!company?.userID) {
      console.error('Company not found');
      goto('/login');
      return [];
    }

    try {
      const jobsRes = await fetch(`/jobs/query?companyID=${company.userID}`);
      if (!jobsRes.ok) throw new Error('Failed to fetch company jobs');

      const jobsData = await jobsRes.json();
      if (!Array.isArray(jobsData) || jobsData.length === 0) return [];

      if (mode === 'jobs') return jobsData;

      const allApplications: any[] = [];

      for (const job of jobsData) {
        const jobID = job.id || job._id;
        if (!jobID) continue;

        const applyRes = await fetch(`/apply?jobID=${jobID}`);
        if (!applyRes.ok) continue;

        const applyData = await applyRes.json();
        const filtered = status
          ? applyData.filter(
              (a: any) => a.jobApplication?.status?.toUpperCase() === status.toUpperCase()
            )
          : applyData;

        for (const app of filtered) {
          allApplications.push({
            ...app,
            jobTitle: job.title || 'Unknown Job',
            closeAt: job.applicationDeadline || null,
            openAt: job.postOpenDate || null,
            jobID
          });
        }
      }

      return allApplications;
    } catch (err) {
      console.error('Error fetching data:', err);
      return [];
    }
  }

  async function loadChartData(view: string) {
    let chartData: any = { labels: [], datasets: [] };
    const axis = getAxisLabels(view);
    let chartOptions: any = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { position: 'bottom' }, tooltip: { mode: 'index', intersect: false },
        title: {
          display: true,
          text: getChartTitle(view),
          font: { size: 18, weight: 'bold' },
          color: '#333',
          padding: { top: 10, bottom: 20 }
      }},
      scales: {
        x: {
          stacked: false,
          title: {
            display: true,
            text: axis.x,
            font: { size: 14, weight: '600' },
            color: '#555'
          }
        },
        y: {
          stacked: false,
          beginAtZero: true,
          title: {
            display: true,
            text: axis.y,
            font: { size: 14, weight: '600' },
            color: '#555'
          }
        }
      },
    };
    let type: 'bar' | 'line' = 'bar';

    switch (view) {
     case 'applicants': {
        type = 'bar';
        const allApps = await fetchData();
        const grouped = groupByMonth(allApps, 'jobApplication.createdAt');
        const labels = getLast12Months();
        const data = labels.map((l) => grouped[l] || 0);

        chartData = {
          labels,
          datasets: [{ label: 'Applicants', data, backgroundColor: '#36A2EB' }]
        };
        break;
      }

      case 'accepted_rejected': {
        type = 'bar';
        const accepted = await fetchData('accepted');
        const rejected = await fetchData('rejected');

        const acceptGroup = groupByMonth(accepted, 'jobApplication.createdAt');
        const rejectGroup = groupByMonth(rejected, 'jobApplication.createdAt');
        const labels = getLast12Months();

        chartOptions.scales.x.stacked = true;
        chartOptions.scales.y.stacked = true;

        chartData = {
          labels,
          datasets: [
            { label: 'Accepted', data: labels.map((l) => acceptGroup[l] || 0), backgroundColor: '#4CAF50' },
            { label: 'Rejected', data: labels.map((l) => rejectGroup[l] || 0), backgroundColor: '#F44336' }
          ]
        };
        break;
      }

      case 'job_postings': {
        type = 'line';
        const jobs = await fetchData('', 'jobs');
        const now = new Date();
        const months: string[] = [];
        for (let i = 11; i >= 0; i--) {
          const dt = new Date(now.getFullYear(), now.getMonth() - i, 1);
          months.push(`${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}`);
        }

        const ongoingCount: Record<string, number> = {};
        months.forEach((m) => (ongoingCount[m] = 0));

        for (const job of jobs) {
          const openAt = new Date(job.postOpenDate || job.createdAt);
          const closeAt = job.applicationDeadline ? new Date(job.applicationDeadline) : null;

          for (const month of months) {
            const [y, m] = month.split('-').map(Number);
            const firstDay = new Date(y, m - 1, 1);
            const lastDay = new Date(y, m, 0);
            if (openAt <= lastDay && (!closeAt || closeAt >= firstDay)) ongoingCount[month]++;
          }
        }

        chartData = {
          labels: months,
          datasets: [
            {
              label: 'Ongoing Job Postings',
              data: months.map((m) => ongoingCount[m]),
              borderColor: '#FF9800',
              backgroundColor: 'rgba(255,152,0,0.3)',
              fill: true,
              tension: 0.3
            }
          ]
        };
        break;
      }

      case 'applied_jobs': {
        type = 'bar';
        const allApps = await fetchData();
        const byJob: Record<string, number> = {};
        for (const a of allApps) byJob[a.jobTitle || 'Unknown Job'] = (byJob[a.jobTitle] || 0) + 1;

        const sorted = Object.entries(byJob).sort((a, b) => b[1] - a[1]).slice(0, 10);

        chartOptions.indexAxis = 'y';
        chartData = {
          labels: sorted.map(([title]) => title),
          datasets: [{ label: 'Applicants', 
          data: sorted.map(([, count]) => count), 
          backgroundColor: '#9C27B0' }], 
        };
        break;
      }
    }

    // Create or update chart with animation apply
    if (chart && chart.config.type !== type) {
      chart.destroy();
      chart = null;
      await tick();
    }

    const animatedChartData = {
      labels: chartData.labels,
      datasets: chartData.datasets.map(ds => ({
        ...ds,
        data: ds.data.map(() => 0)
      }))
    };

    if (!chart) {
      chart = new Chart(canvas, {
        type,
        data: animatedChartData,
        options: {
          ...chartOptions,
          animation: {
            duration: 1000,
            easing: 'easeOutQuart',
            animateScale: true
          }
        }
      });

      setTimeout(() => {
        chart!.data.datasets = chartData.datasets;
        chart!.update({
          duration: 1000,
          easing: 'easeOutQuart'
        });
      }, 50);
    } else {
      chart.data = chartData;
      chart.options = chartOptions;
      chart.update({
        duration: 1000,
        easing: 'easeOutQuart'
      });
    }
  }

  selectedView.subscribe(async (view) => {
    if (company) await loadChartData(view);
  });

  onMount(async () => {
    const auth = await isAuthenticated();
    if (!auth) return goto('/login');
    company = getUserInfo();
    if (canvas) await loadChartData('applicants');
  });
</script>

<div class="flex flex-col w-full">
  <h1 class="text-2xl font-semibold text-gray-900 mb-1">
    Company Analytics
  </h1>
  <p class="mb-6 text-base text-gray-600">
    Welcome back, company HR Team. Here's analytics for your company job and recruitment.
  </p>
  <!-- View Switch Buttons -->
  <div class="flex flex-wrap gap-2 mb-6">
    {#each viewOptions as opt, i (i)}
      <button
        class="px-4 py-2 rounded-lg font-medium transition-colors duration-200 
               border border-gray-300
               hover:bg-green-300
               data-[active=true]:bg-green-600 data-[active=true]:text-white"
        data-active={$selectedView === opt.value}
        on:click={() => selectedView.set(opt.value)}
      >
        {opt.label}
      </button>
    {/each}
  </div>

  <!-- Chart Area -->
  <div class="relative h-[60vh] w-full bg-white rounded-2xl p-4 shadow-md">
    <canvas bind:this={canvas}></canvas>
  </div>
</div>
