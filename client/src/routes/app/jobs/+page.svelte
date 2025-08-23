<div class="flex h-screen">
  <main class="flex-1 p-0.5 space-y-6 w-full min-w-0">
    <div class="flex items-center gap-3">
      <!-- Search bar -->
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search jobs..."
          class="flex-1 ml-2 outline-none border-none"
        />
      </div>
      <!-- Filters -->
      <div class="flex gap-2">
        <button class="px-3 py-1 bg-gray-200 rounded-full text-sm">Full-time</button>
        <button class="px-3 py-1 bg-gray-200 rounded-full text-sm">1 day ago</button>
        <button class="px-3 py-1 bg-gray-200 rounded-full text-sm">Remote</button>
        <button class="px-3 py-1 bg-gray-200 rounded-full text-sm">AI</button>
        <button class="px-3 py-1 bg-gray-200 rounded-full text-sm flex items-center"><Filter class="h-4 w-4 mr-1" /> Filters</button>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 w-full">
      <!--Jobs List -->
      <section class="col-span-1 flex flex-col space-y-4 min-w-0 overflow-hidden">
        <div class="flex-1 overflow-y-auto w-full space-y-3">
          {#each jobs as job}
          <button
            on:click={() => (selectedJob = job)}
            class="w-full flex items-start gap-3 cursor-pointer bg-white p-2 rounded-lg shadow hover:ring-2 hover:ring-green-300 {selectedJob.id === job.id ? 'ring-2 ring-green-500' : ''}"
          >
            <img src={job.logo} alt={job.company} class="w-10 h-10 rounded-full object-cover flex-shrink-0" />
            <div class="flex flex-col flex-1 text-left">
              <div class="font-semibold">{job.title}</div>
              <div class="text-sm text-gray-600">{job.company}</div>
              <div class="text-sm text-black items-center flex gap-2 mt-1">
                <MapPin class="w-4 h-4 mt-0.5" />
                {job.location}
              </div>
            <div class="flex gap-2 mt-2 flex-wrap">
              {#each job.tags as tag}
                <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
              {/each}
            </div>
            <div class="text-xs text-gray-500 mt-1">{job.posted}</div>
            </div>
          </button>
          {/each}
        </div>
      </section>

      <!--Job Detail -->
      <section class="col-span-1 bg-white p-6 rounded-lg shadow space-y-2 min-w-0 overflow-hidden">
        {#if selectedJob}
          <div class="mb-4">
            <h2 class="text-xl font-bold">{selectedJob.title}</h2>
            <p class="text-gray-600">{selectedJob.company} – {selectedJob.location}</p>
            <p class="text-sm text-gray-500">{selectedJob.posted}</p>
          </div>
          <div class="flex gap-2">
            {#each selectedJob.tags as tag}
              <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
            {/each}
          </div>
          <div class="flex gap-2 mb-4">
            <button class="px-4 py-2 bg-green-500 text-white rounded-lg">Apply</button>
            <button class="px-4 py-2 bg-gray-200 rounded-lg">Bookmark</button>
          </div>
          <div class="space-y-4">
            <p>{selectedJob.description}</p>
          </div>
        {/if}
      </section>
    </div>
  </main>
</div>

<script>
  import { Search, Filter, MapPin } from 'lucide-svelte';

  let jobs = [
    {
      id: 1,
      title: "Software Engineer, AI Acceleration",
      company: "Google",
      location: "Singapore",
      type: "Full-time",
      tags: ["Remote", "AI"],
      posted: "1 day ago",
      description: "Google wants to have more software engineer for researching and improving AI abilities.",
      logo: "https://logo.clearbit.com/google.com",
    },
    {
      id: 2,
      title: "Software Developer, Back End",
      company: "Agoda",
      location: "Bangkok, Thailand",
      type: "Full-time",
      tags: ["Backend", "Node.js"],
      posted: "2 days ago",
      description: "Backend-focused developer role...",
      logo: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRbx0J4R5cqy54hvPd9lNP8ywO5kGWl5JlO9A&s",
    },
    {
      id: 3,
      title: "Fullstack Web Developer",
      company: "Kasikorn Bank",
      location: "Bangkok, Thailand",
      type: "Full-time",
      tags: ["React", "Node.js"],
      posted: "5 days ago",
      description: "Work across frontend and backend...",
      logo: "https://logo.clearbit.com/kasikornbank.com",
    },
  ];

  let selectedJob = jobs[0];
</script>